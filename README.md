# 🐫 Go Bumpy

[![Go package](https://github.com/survivorbat/go-bumpy/actions/workflows/test.yaml/badge.svg)](https://github.com/survivorbat/go-bumpy/actions/workflows/test.yaml)
![GitHub](https://img.shields.io/github/license/survivorbat/go-bumpy)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/survivorbat/go-bumpy)

Go-bumpy is a simple tool for bumping the version of your go project based on [Semantic Versioning](https://semver.org/).
Not only does it look at existing tags in your repository, it also
reads the version from your `go.mod` file to determine what the major version of your project is.

## ⬇️ Installation

Check out the [releases](https://github.com/survivorbat/go-bumpy/releases) or use `go install github.com/survivorbat/go-bumpy/cmd/bumpy`.

## 📋 Usage

`bumpy [-minor] .`

## 🔭 Plans

Add option for auto pushing the new tag to the remote.
