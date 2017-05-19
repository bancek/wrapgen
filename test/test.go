// Package wrapgentest is used by the wrapgen test suite to test the package parser.
// While it is completely possible to generate AST nodes and unit test the
// parser code it is much easier to define a package and set of elements
// to be parsed.
//
package wrapgentest

import (
	"io"
	nethttp "net/http"
	"os"
)

type unexportedStruct struct {
	A string
	B int
	C chan bool
	D []uint64
}

type unexportedInterface interface {
	A()
	B(int, bool) (string, error)
	C(one int, two bool) (string, error)
	D(one unexportedStruct, two *unexportedStruct) error
}

type unexportedInterfaceWithEmbedded interface {
	unexportedInterface
}

type ExportedStruct struct {
	A string
	B int
	C chan bool
	D []uint64
}

type ExportedInterface interface {
	A()
	B(int, bool) (string, error)
	C(one int, two bool) (string, error)
	D(one ExportedStruct, two *ExportedStruct) error
	E(one func(), two func(int) bool) error
	F(one chan bool, two <-chan bool, three chan<- bool) error
	G(one []string, two [100]string) error
	H(one os.File, two *os.File) error
	I(one os.FileInfo) error
	J(one map[string]string) error
	K(one ...string) error
	L(one interface{}, two struct{}) error
	M(one nethttp.Handler) error
}

type ExportedInterfaceWithEmbedded interface {
	ExportedInterface
}

type ExportedInterfaceWithRemoteEmbedded interface {
	io.Reader
}
