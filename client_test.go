package main

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubernetes/pkg/controller/testutil"

	"github.com/stretchr/testify/require"
	"testing"
)

func newTestClient() *Client {
	client := Client{}
	objects := []runtime.Object{
		&v1.PodList{Items: []v1.Pod{*testutil.NewPod("pod1", "host1")}},
	}
	client.kubeClient = fake.NewSimpleClientset(objects...)
	return &client
}

func TestClientGetPods(t *testing.T) {
	tc := newTestClient()
	pods, err := tc.Pods("", "")
	require.NoError(t, err)
	require.Len(t, pods, 1)
	require.Equal(t, "pod1", pods[0].Name)
	require.Equal(t, "host1", pods[0].Spec.NodeName)
}
