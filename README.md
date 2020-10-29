# smime-verify

![Build](https://github.com/danielb42/smime-verify/workflows/Build/badge.svg)
![Tag](https://img.shields.io/github/v/tag/danielb42/smime-verify)
![Go Version](https://img.shields.io/github/go-mod/go-version/danielb42/smime-verify)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/danielb42/smime-verify)](https://pkg.go.dev/github.com/danielb42/smime-verify)
[![Go Report Card](https://goreportcard.com/badge/github.com/danielb42/smime-verify)](https://goreportcard.com/report/github.com/danielb42/smime-verify)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

Verification of S/MIME messages signed by `TeleSec Business CA 1` intermediate certification.

## Install

### Either download a precompiled binary ...

Available for Linux, Windows and MacOS: [Latest Release](https://github.com/danielb42/smime-verify/releases/latest)

### ... or use go get

`go get github.com/danielb42/smime-verify`

## Usage

```bash
smime-verify <filename>
```
