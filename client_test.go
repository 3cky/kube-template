package main

import (
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubernetes/pkg/controller/testutil"

	"github.com/stretchr/testify/require"
	"testing"
)

func TestClientGetPodsDirectly(t *testing.T) {
	testClientGetPods(t, false)
}

func TestClientGetPodsUsingInformer(t *testing.T) {
	testClientGetPods(t, true)
}

func testClientGetPods(t *testing.T, useInformer bool) {
	pod := testutil.NewPod("pod1", "host1")
	pod.ObjectMeta.Labels = make(map[string]string)
	pod.ObjectMeta.Labels["name"] = "pod1"

	fakeClient := fake.NewSimpleClientset(pod)

	stopCh := make(chan struct{})
	defer close(stopCh)

	tc, err := newClient(fakeClient, stopCh, useInformer)
	require.NoError(t, err)

	pods, err := tc.Pods("", "")
	require.NoError(t, err)
	require.Len(t, pods, 1)
	require.Equal(t, "pod1", pods[0].Name)
	require.Equal(t, "host1", pods[0].Spec.NodeName)

	pods, err = tc.Pods("", "name=unknown")
	require.NoError(t, err)
	require.Empty(t, pods)

	pods, err = tc.Pods("", "name=pod1")
	require.NoError(t, err)
	require.Len(t, pods, 1)
	require.Equal(t, "pod1", pods[0].Name)
	require.Equal(t, "host1", pods[0].Spec.NodeName)
}
