package mylib

import (
	"fmt"
	"os"
	"time"
)

type NodeChangeContent struct {
	next Noder
	name string
}

func (n *NodeChangeContent) GetName() string {
	return n.name
}

func (n *NodeChangeContent) GetNextNode() Noder {
	return n.next
}
func (n *NodeChangeContent) SetNextNode(noder Noder) {
	n.next = noder
}
func (n *NodeChangeContent) Execute(filepath string) (string, error) {
	newContent := time.Now().Format(time.RFC3339Nano)

	if err := os.WriteFile(filepath, []byte(newContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write data in file (%s): %w", filepath, err)
	}

	return newContent, nil
}

func NewNodeChangeContent() *NodeChangeContent {
	return &NodeChangeContent{
		name: "change_content",
	}
}
