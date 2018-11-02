package main

import (
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubernetes/pkg/controller/testutil"

	"github.com/stretchr/testify/require"
	"testing"
)

func testTemplateName(t *testing.T) string {
	return testDataFilePrefix(t) + ".template"
}

func TestTemplateRender(t *testing.T) {
	pod := testutil.NewPod("pod1", "host1")
	fakeClient := fake.NewSimpleClientset(pod)

	stopCh := make(chan struct{})
	defer close(stopCh)

	tc, err := newClient(fakeClient, stopCh, true)
	require.NoError(t, err)

	dm := newDependencyManager(tc)

	testtemplate := testTemplateName(t)
	testgolden := testGoldenFileName(t)
	td, err := parseTemplateDescriptor(fmt.Sprintf("%s:%s", testtemplate, testgolden))
	require.NoError(t, err)

	cfg := new(Config)
	template, err := newTemplate(cfg, dm, td)
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
