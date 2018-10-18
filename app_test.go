package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/stretchr/testify/require"
	"testing"
)

func TestAppRunOnce(t *testing.T) {
	cmd := newCmd()
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	testgolden := testGoldenFileName(t)
	var testout string
	if *update {
		testout = testgolden
	} else {
		testoutfile, err := ioutil.TempFile("", "testout")
		require.NoError(t, err)
		testout = testoutfile.Name()
		defer testFileUnlink(testout)
	}
	err := cmd.ParseFlags([]string{
		"--once",
		fmt.Sprintf("--template=%s:%s", testTemplateName(t), testout),
	})
	require.NoError(t, err)
	cfg, err := newConfig(cmd)
	require.NoError(t, err)
	app, err := newApp(cfg)
	require.NoError(t, err)
	tc := newTestClient()
	app.client = tc
	app.dm.client = tc
	app.RunOnce()
	actual, err := ioutil.ReadFile(testout)
	require.NoError(t, err)
	expected, err := ioutil.ReadFile(testgolden)
	require.NoError(t, err)
	require.Equal(t, string(expected), string(actual))
}
