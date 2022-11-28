//go:build tools
// +build tools

// Package tools tracks dependencies for tools.
// See https://github.com/golang/go/issues/25922.
package tools

import (
	_ "github.com/onsi/ginkgo/v2/ginkgo"
	_ "github.com/square/certstrap"
)
