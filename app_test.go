package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"
	"k8s.io/kubernetes/pkg/controller/testutil"

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
		defer UnlinkQuietly(testout)
	}

	err := cmd.ParseFlags([]string{
		"--once",
		fmt.Sprintf("--template=%s:%s", testTemplateName(t), testout),
	})
	require.NoError(t, err)

	cfg, err := newConfig(cmd)
	require.NoError(t, err)

	fakeClient := &fake.Clientset{}

	fakeWatch := watch.NewFake()
	fakeClient.AddWatchReactor("pods", ktesting.DefaultWatchReactor(fakeWatch, nil))

	stopCh := make(chan struct{})
	doneCh := make(chan struct{})
	defer close(stopCh)
	defer close(doneCh)

	tc, err := newClient(fakeClient, stopCh)
	require.NoError(t, err)

	pod := testutil.NewPod("pod1", "host1")
	fakeWatch.Add(pod)

	dm := newDependencyManager(tc)

	templates, err := newTemplatesFromConfig(cfg, dm)
	require.NoError(t, err)

	app := &App{
		stopCh:       stopCh,
		doneCh:       doneCh,
		dm:           dm,
		templates:    templates,
		dryRun:       cfg.DryRun,
		updatePeriod: cfg.PollPeriod,
	}

	app.RunOnce()
	actual, err := ioutil.ReadFile(testout)
	require.NoError(t, err)

	expected, err := ioutil.ReadFile(testgolden)
	require.NoError(t, err)
	require.Equal(t, string(expected), string(actual))
}
