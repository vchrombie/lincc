# licc

[![Go Report Card](https://goreportcard.com/badge/github.com/vchrombie/licc)](https://goreportcard.com/report/github.com/vchrombie/licc)
[![GoDoc](https://godoc.org/github.com/vchrombie/licc?status.svg)](https://godoc.org/github.com/vchrombie/licc)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

licc is a tool designed to automate the verification and correction of software
license configurations across projects.

## Installation

```bash
go get -u github.com/vchrombie/licc
```

## Usage

```bash
go run main.go https://github.com/vchrombie/licc
```

## Roadmap

- [x] Evaluation of content types for appropriate license application
- [ ] Integration with CI/CD pipelines and generates a compliance badge
- [ ] Use cli framework for proper cli interface
- [ ] Add option to ignore/exclude parts of the repository
- [ ] Identification and correction of license inconsistencies

## License

Apache-2.0
