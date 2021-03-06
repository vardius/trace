📡 trace
================
[![Build Status](https://travis-ci.org/vardius/trace.svg?branch=master)](https://travis-ci.org/vardius/trace)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/trace)](https://goreportcard.com/report/github.com/vardius/trace)
[![codecov](https://codecov.io/gh/vardius/trace/branch/master/graph/badge.svg)](https://codecov.io/gh/vardius/trace)
[![](https://godoc.org/github.com/vardius/trace?status.svg)](https://pkg.go.dev/github.com/vardius/trace)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/vardius/trace/blob/master/LICENSE.md)

<img align="right" height="180px" src="https://github.com/vardius/gorouter/blob/master/website/src/static/img/logo.png?raw=true" alt="logo" />

trace - simple helper to trace the function calls, errors or logs reference

📖 ABOUT
==================================================
Contributors:

* [Rafał Lorenz](http://rafallorenz.com)

Want to contribute ? Feel free to send pull requests!

Have problems, bugs, feature ideas?
We are using the github [issue tracker](https://github.com/vardius/trace/issues) to manage them.

## 📚 Documentation

For __examples__ **visit [godoc#pkg-examples](http://godoc.org/github.com/vardius/trace#pkg-examples)**

For **GoDoc** reference, **visit [pkg.go.dev](https://pkg.go.dev/github.com/vardius/trace)**

🚏 HOW TO USE
==================================================

## 🏫 Basic example
```go
package main

import (
    "log"

    "github.com/vardius/trace"
)

func main() {
    log.Printf("%s\n\t%s", "Hello from:", trace.Here(trace.Lfile | trace.Lline | trace.Lfunction))
}
```

📜 [License](LICENSE.md)
-------

This package is released under the MIT license. See the complete license in the package
