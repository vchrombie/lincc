package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-enry/go-license-detector/v4/licensedb"
)

// structs
type FileLicenseMapping struct {
	Software          Category `json:"software"`
	Documentation     Category `json:"documentation"`
	Multimedia        Category `json:"multimedia"`
	DataSetsAndModels Category `json:"data_sets_and_models"`
}

type Category struct {
	Extensions []string `json:"extensions"`
	Licenses   []string `json:"licenses"`
}

type Config struct {
	IgnorePatterns []string `json:"ignorePatterns"`
}

func loadConfig() (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(".lincc")
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err
}

func shouldIgnore(filePath string, config Config) bool {
	for _, pattern := range config.IgnorePatterns {
		// Normalize leading './' for consistent matching
		cleanPath := strings.TrimPrefix(filePath, "./")
		// Match using filepath.Match which supports glob patterns
		matched, err := filepath.Match(pattern, cleanPath)
		if err != nil {
			// Handle syntax error in pattern
			log.Printf("Invalid pattern syntax %s: %v", pattern, err)
			continue
		}
		if matched {
			return true
		}
		// Attempt to match pattern within any subdirectory if it's not explicitly rooted
		if !strings.HasPrefix(pattern, "/") {
			matched, err = filepath.Match("*/"+pattern, cleanPath)
			if err == nil && matched {
				return true
			}
		}
	}
	return false
}

func splitPattern(pattern string) (dirPattern, filePattern string) {
	lastSlash := strings.LastIndex(pattern, "/")
	dirPattern = pattern[:lastSlash]
	filePattern = pattern[lastSlash+1:]
	return
}

func matchDirAndFilePattern(path, dirPattern, filePattern string) bool {
	dir, file := filepath.Split(path)
	// Check if directory matches the dirPattern
	matchedDir, _ := filepath.Match(dirPattern+"/*", dir)
	if matchedDir {
		// Check if file matches the filePattern
		matchedFile, _ := filepath.Match(filePattern, file)
		if matchedFile {
			return true
		}
	}
	return false
}
func cloneRepo(repoURL, dir string) error {
	cmd := exec.Command("git", "clone", "--depth", "1", repoURL, dir)
	return cmd.Run()
}

func loadMapping() (FileLicenseMapping, error) {
	var mapping FileLicenseMapping
	data, err := ioutil.ReadFile("mapping.json")
	if err != nil {
		return mapping, err
	}
	err = json.Unmarshal(data, &mapping)
	return mapping, err
}

func getRootLicenses(dir string) ([]string, error) {
	results := licensedb.Analyse(dir)
	var rootLicenses []string

	for _, result := range results {
		fileLicenses := make(map[string]licensedb.Match)
		for _, match := range result.Matches {
			if existingMatch, exists := fileLicenses[match.File]; !exists || match.Confidence > existingMatch.Confidence {
				fileLicenses[match.File] = match
			}
		}
		for _, license := range fileLicenses {
			rootLicenses = append(rootLicenses, license.License)
		}
		if len(rootLicenses) == 0 {
			return nil, fmt.Errorf("no root licenses found")
		}
	}
	return rootLicenses, nil
}

func contains(slice []string, item string) bool {
	for _, sliceItem := range slice {
		if sliceItem == item {
			return true
		}
	}
	return false
}

func isLicenseApplicable(extension string, rootLicenses []string, mapping FileLicenseMapping) bool {
	var categories []Category = []Category{
		mapping.Software, mapping.Documentation,
		mapping.Multimedia, mapping.DataSetsAndModels,
	}

	for _, license := range rootLicenses {
		for _, category := range categories {
			if contains(category.Extensions, extension) && contains(category.Licenses, license) {
				return true
			}
		}
	}
	return false
}

func checkFiles(dir string, rootLicenses []string, mapping FileLicenseMapping, config Config) map[string]bool {
	fileChecks := make(map[string]bool)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// New logic to check if the path matches any ignore pattern
		if shouldIgnore(path, config) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			// check if its license is applicable
			relativePath, err := filepath.Rel(dir, path)
			if err != nil {
				log.Printf("Error getting relative path for %s: %v", path, err)
				return nil
			}
			extension := filepath.Ext(relativePath)
			applicable := isLicenseApplicable(extension, rootLicenses, mapping)
			fileChecks[relativePath] = applicable
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking directory: %v", err)
	}

	return fileChecks
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run script.go <repo-url>")
	}
	repoURL := os.Args[1]

	mapping, err := loadMapping()
	if err != nil {
		log.Fatalf("Failed to load mapping: %v", err)
	}

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load .lincc config: %v", err)
	}

	dir := filepath.Join(".", strings.Split(filepath.Base(repoURL), ".")[0])
	if err := cloneRepo(repoURL, dir); err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
	}
	defer func() { exec.Command("rm", "-rf", dir).Run() }()

	rootLicenses, err := getRootLicenses(dir)
	if err != nil {
		log.Fatalf("Failed to determine root licenses: %v", err)
	}

	fmt.Printf("Project: %s\n", strings.Split(filepath.Base(repoURL), ".")[0])
	for _, license := range rootLicenses {
		fmt.Printf("License: %s\n", license)
	}

	fileChecks := checkFiles(dir, rootLicenses, mapping, config)

	var files []string
	for file := range fileChecks {
		files = append(files, file)
	}

	sort.Strings(files)

	// some stats

	totalFiles := 0
	compliantFiles := 0
	notCompliantFiles := 0

	fmt.Println("\nFiles:")
	for _, file := range files {
		applicable := fileChecks[file]
		totalFiles++
		if applicable {
			fmt.Printf("%s: ✅\n", file)
			compliantFiles++
		} else {
			fmt.Printf("%s: ❌\n", file)
			notCompliantFiles++
		}
	}

	fmt.Printf("\nTotal files: %d\n", totalFiles)
	fmt.Printf("Compliant files: %d\n", compliantFiles)
	fmt.Printf("Non-compliant files: %d\n", notCompliantFiles)

	score := float64(compliantFiles) / float64(totalFiles) * 100
	fmt.Printf("\nScore: %.2f%%\n", score)
}
