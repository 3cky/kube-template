package main

import (
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"
	"k8s.io/kubernetes/pkg/controller/testutil"

	"github.com/stretchr/testify/require"
	"testing"
)

func TestClientGetPods(t *testing.T) {
	fakeClient := &fake.Clientset{}

	fakeWatch := watch.NewFake()
	fakeClient.AddWatchReactor("pods", ktesting.DefaultWatchReactor(fakeWatch, nil))

	stopCh := make(chan struct{})
	defer close(stopCh)

	tc, err := newClient(fakeClient, stopCh)
	require.NoError(t, err)

	pods, err := tc.Pods("", "")
	require.NoError(t, err)
	require.Empty(t, pods)

	pod := testutil.NewPod("pod1", "host1")
	fakeWatch.Add(pod)

	pods, err = tc.Pods("", "")
	require.NoError(t, err)
	require.Len(t, pods, 1)
	require.Equal(t, "pod1", pods[0].Name)
	require.Equal(t, "host1", pods[0].Spec.NodeName)
}
