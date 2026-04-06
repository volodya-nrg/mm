package mylib

import (
	"testing"
)

func TestMyLib(t *testing.T) {
	node1 := NewNodeChangeContent()
	node2 := NewNodeGetMD5()

	node1.SetNextNode(node2) // соберем цепочку нод

	myLib := NewMyLib(
		"./test_data",
		3,
		node1,
	)

	for resp := range myLib.Run(t.Context()) {
		t.Logf("%s", resp)
	}
}
