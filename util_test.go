package main

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

func testDataFilePrefix(t *testing.T) string {
	return filepath.Join("testdata", t.Name())
}

func testGoldenFileName(t *testing.T) string {
	return testDataFilePrefix(t) + ".out.golden"
}

func TestIsPresent(t *testing.T) {
	assert.True(t, IsPresent([]string{"a", "b", "c"}, "a"))
	assert.True(t, IsPresent([]string{""}, ""))
	assert.False(t, IsPresent([]string{}, ""))
	assert.False(t, IsPresent([]string{}, "a"))
	assert.False(t, IsPresent([]string{"d", "e", "f"}, "a"))
}
