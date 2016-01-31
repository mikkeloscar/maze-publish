package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadBuiltPkg(t *testing.T) {
	file := "packages.built"
	f, err := os.Create(file)
	assert.NoError(t, err, "should not fail")
	_, err = f.WriteString(`[{"package": "path/to/file", "signature":"sig"}]`)
	assert.NoError(t, err, "should not fail")
	err = f.Close()
	assert.NoError(t, err, "should not fail")

	var pkgs []*BuiltPkg
	err = loadBuiltPkgs(file, &pkgs)
	assert.NoError(t, err, "should not fail")
	assert.Len(t, pkgs, 1, "should have len 1")
	assert.Equal(t, pkgs[0].Package, "path/to/file", "should be equal")
	assert.Equal(t, pkgs[0].Signature, "sig", "should be equal")

	err = os.Remove(file)
	assert.NoError(t, err, "should not fail")
}
