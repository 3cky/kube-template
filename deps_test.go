package main

import (
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"
	"k8s.io/kubernetes/pkg/controller/testutil"

	"testing"
)

func TestDependencyManager(t *testing.T) {
	fakeClient := &fake.Clientset{}

	fakeWatch := watch.NewFake()
	fakeClient.AddWatchReactor("pods", ktesting.DefaultWatchReactor(fakeWatch, nil))

	stopCh := make(chan struct{})
	defer close(stopCh)

	tc, err := newClient(fakeClient, stopCh)
	require.NoError(t, err)

	pod := testutil.NewPod("pod1", "host1")
	fakeWatch.Add(pod)

	dm := newDependencyManager(tc)
	require.Empty(t, dm.cachedDeps)

	pods, err := dm.Pods("", "")
	require.NoError(t, err)
	require.Len(t, pods, 1)
	require.NotEmpty(t, dm.cachedDeps)
	pod1 := pods[0]
	require.Equal(t, pod.Name, pod1.Name)

	pods, err = dm.Pods("", "")
	require.NoError(t, err)
	require.Len(t, pods, 1)
	require.Equal(t, pod1, pods[0])
}
