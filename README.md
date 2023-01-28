# 🐫 Go Bumpy

[![Go package](https://github.com/survivorbat/go-bumpy/actions/workflows/test.yaml/badge.svg)](https://github.com/survivorbat/go-bumpy/actions/workflows/test.yaml)
![GitHub](https://img.shields.io/github/license/survivorbat/go-bumpy)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/survivorbat/go-bumpy)

Go-bumpy is a simple tool for bumping the version of your go project based on [Semantic Versioning](https://semver.org/).
Not only does it look at existing tags in your repository, it also
reads the version from your `go.mod` file to determine what the major version of your project is.

It is also capable of pushing the new tag to your remote repository.

## ⬇️ Installation

`go install github.com/survivorbat/go-bumpy/cmd/bumpy`

Or check out the [releases](https://github.com/survivorbat/go-bumpy/releases).

## 📋 Usage

`bumpy [-minor] [-push="origin"] <directory>`

It will output the new tag name to stdout and logging to stderr.

### Options

- `-minor` Bump the minor version instead of the patch version
- `-push` Push the new tag to the specified remote. If not specified, the tag will not be pushed.

### Examples

| Module Suffix | Latest Tag | Output |
|---------------|------------|--------|
| None          | None       | v0.0.0 |
| None          | v2.5.0     | v2.5.1 |
| v3            | None       | v3.0.0 |
| v3            | v3.2.0     | v3.2.1 |
| v3            | v5.4.3     | v3.0.0 |

## 🔭 Plans

None yet
