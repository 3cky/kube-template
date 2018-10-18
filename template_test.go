package main

import (
	"fmt"
	"io/ioutil"

	"github.com/stretchr/testify/require"
	"testing"
)

func testTemplateName(t *testing.T) string {
	return testDataFilePrefix(t) + ".template"
}

func TestTemplateRender(t *testing.T) {
	client := newTestClient()
	dm := newDependencyManager(client)
	testtemplate := testTemplateName(t)
	testgolden := testGoldenFileName(t)
	td, err := parseTemplateDescriptor(fmt.Sprintf("%s:%s", testtemplate, testgolden))
	require.NoError(t, err)
	cfg := new(Config)
	template, err := newTemplate(dm, td, cfg)
	require.NoError(t, err)
	if *update {
		_, err := template.Process(false)
		require.NoError(t, err)
	}
	actual, err := template.Render()
	require.NoError(t, err)
	expected, err := ioutil.ReadFile(testgolden)
	require.NoError(t, err)
	require.Equal(t, string(expected), actual)
}
