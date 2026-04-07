package mylib

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMyLib(t *testing.T) {
	node1 := NewNodeChangeContent()
	node2 := NewNodeGetMD5()

	node1.SetNextNode(node2) // соберем цепочку нод

	myLib, err := NewMyLib(
		"./test_data",
		3,
		node1,
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	for resp := range myLib.Run(ctx) {
		t.Logf("%s", resp)
	}
}
