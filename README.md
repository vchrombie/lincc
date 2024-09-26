# lincc

[![Go Report Card](https://goreportcard.com/badge/github.com/vchrombie/lincc)](https://goreportcard.com/report/github.com/vchrombie/lincc)
[![GoDoc](https://godoc.org/github.com/vchrombie/lincc?status.svg)](https://godoc.org/github.com/vchrombie/lincc)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

lincc is a tool designed to automate the verification and correction of software
license configurations across projects.

It checks if the licenses applied to project components, such as documentation and
code, are properly aligned with the chosen project license.

## Installation

```bash
go get -u github.com/vchrombie/lincc
```

## Usage

```bash
go run main.go https://github.com/vchrombie/lincc
```

## Roadmap

- [x] Evaluation of content types for appropriate license application
- [ ] Integration with CI/CD pipelines and generates a compliance badge
- [ ] Use cli framework for proper cli interface
- [ ] Add option to ignore/exclude parts of the repository
- [ ] Identification and correction of license inconsistencies

## License

Apache-2.0
