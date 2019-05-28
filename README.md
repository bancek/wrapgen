# wrapgen - A code generator for Go interfaces

This project can inject the details of any Go interface into a custom template.
The goal is to enable teams who either make frequent use of interface wrappers
to layer on behavior, such as using a decorator pattern, or teams who often need
to generate standard implementations of some interfaces, such as those genrating
stubs for a test.

This project is heavily inspired by <https://github.com/golang/mock> which
contains a tool called `mockgen` that can generate an advanced mock
implementation of any Go interface.

## Usage

### Quick Example

```bash
go install github.com/kevinconway/wrapgen
go install golang.org/x/tools/cmd/goimports

${GOPATH}/bin/wrapgen \
  --source=io \
  --interface=Reader \
  --interface=Writer \
  --package=wrappers \
  --template="https://github.com/kevinconway/wrapgen/templates/logtime.txt" \
  | ${GOPATH}/bin/goimports
```

The output will look like:

```golang
package wrappers

import (
        "io"
        "log"
        "time"
)

type WrapsReader struct {
        wrapped io.Reader
}

func (w *WrapsReader) Read(p []byte) (int, error) {
        start := time.Now()
        defer func() {
                log.Println("Read latency:", time.Now().Since(start))
        }()
        var n, err = w.Read(p)
        return n, err
}

type WrapsWriter struct {
        wrapped io.Writer
}

func (w *WrapsWriter) Write(p []byte) (int, error) {
        start := time.Now()
        defer func() {
                log.Println("Write latency:", time.Now().Since(start))
        }()
        var n, err = w.Write(p)
        return n, err
}
```

All output is written to `stdout` and `stderr`. Template paths may either be
URLs or file system paths. It is highly suggested to use `goimports` or another
formatter as a post-processor.

### CLI Options

```bash
wrapgen --help

Usage of wrapgen:
      --interface strings   The name of the interface to render.
      --leftdelim string    Left-hand side delimiter for the template. (default "#!")
      --package string      The package name that the resulting file will be in. Defaults to the source package.
      --rightdelim string   Right-hand side delimiter for the template. (default "!#")
      --source string       The import path of the package to render.
      --template string     The template to render.
      --timeout duration    Maximum runtime allowed for rendering. (default 1m0s)
```

Any number of interfaces may be given by providing more `--interface` flags.

### Writing Templates

The `templates/basic.txt` template from this project is the best way to get
started writing your own. It demonstrates the how to manage import statements,
iterate over the collected interfaces, and already covers the complexity of
rendering method arguments and outputs correctly.

Whether using `template/basics.txt` as a starter or generating a new template
from scratch, the template content must be valid `text/template` markup. By
default, the character sets `#!` and `!#` are used as the left and right
delimiters, respectively. Those characters are the default because they rarely,
if ever, conflict with common character sets in Go code. You can adjust these
with CLI flags.

The root context injected into the template is `Package` from the following:

```golang
// Package is a container for all exported interfaces of a Go package.
type Package struct {
	Name       string
	Source     *Import
	Interfaces []*Interface
	Imports    []*Import
}

// Import is a package name and path that is imported by another package.
type Import struct {
	Package string
	Path    string
}

// Interface is an exported interface defined in a package.
type Interface struct {
	Name    string
	Methods []*Method
}

// Method is a named function attached to an interface.
type Method struct {
	Name string
	In   []*Parameter
	Out  []*Parameter
}

// Parameter is a named parameter used by a Method.
type Parameter struct {
	Name string
	Type Type
}

// Type is a Go type definition that can be rendered into a valid
// Go code snippet.
type Type interface {
	String() string
}
```

Those are the structures available within any template. Each `Type` is
specialized and will render correctly when calling `String()`. For example, a
`Parameter` with `Type` of read-only channel of integers will render as `<-chan
int` when calling `String()` on the type. The `templates/basic.txt` template
contains examples of inspecting the type string in order to support variadics.
This practice can be extended to, for example, determine if the first parameter
is a context and optionall fetch a value from it.

## License

This project is available under the Apache2.0 license. See the `LICENSE` file
in this repository for the complete license text.

