//go:build builddeps
// +build builddeps

package main

import (
	// This package is specified under an alternate build tag because it is a
	// build-time dependency that needs to be tracked in go.mod in order to
	// replace it with an alternate version.
	_ "github.com/xlab/c-for-go"
)
