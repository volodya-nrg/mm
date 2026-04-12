package mylib

import (
	"fmt"
	"os"
	"time"
)

type NodeChangeContent[T string] struct {
	next Noder[T]
	name string
}

func (n *NodeChangeContent[T]) GetName() string {
	return n.name
}

func (n *NodeChangeContent[T]) GetNextNode() Noder[T] {
	return n.next
}
func (n *NodeChangeContent[T]) SetNextNode(noder Noder[T]) {
	n.next = noder
}
func (n *NodeChangeContent[T]) Execute(ch chan T) (chan T, error) {
	filePath := <-ch
	newContent := time.Now().Format(time.RFC3339Nano)

	if err := os.WriteFile(string(filePath), []byte(newContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write data in file (%s): %w", filePath, err)
	}

	chOut := make(chan T, 1)
	chOut <- T(newContent)
	close(chOut)

	return chOut, nil
}

func NewNodeChangeContent[T string]() *NodeChangeContent[T] {
	return &NodeChangeContent[T]{
		name: "change_content",
	}
}
