package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDependencyManager(t *testing.T) {
	client := newTestClient()
	dm := newDependencyManager(client)
	require.Empty(t, dm.cachedDeps)
	pods, err := dm.Pods("", "")
	require.NoError(t, err)
	require.Len(t, pods, 1)
	require.NotEmpty(t, dm.cachedDeps)
	pod := pods[0]
	pods, err = dm.Pods("", "")
	require.NoError(t, err)
	require.Len(t, pods, 1)
	require.Equal(t, pod, pods[0])
}
