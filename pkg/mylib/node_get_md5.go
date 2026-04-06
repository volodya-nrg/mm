package mylib

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type NodeGetMD5 struct {
	next Noder
	name string
}

func (n *NodeGetMD5) GetName() string {
	return n.name
}

func (n *NodeGetMD5) GetNextNode() Noder {
	return n.next
}

func (n *NodeGetMD5) SetNextNode(noder Noder) {
	n.next = noder
}

func (n *NodeGetMD5) Execute(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file (%s): %w", filePath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	hash := md5.New()

	if _, err = io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calc MD5-hash: %w", err)
	}

	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, nil
}

func NewNodeGetMD5() *NodeGetMD5 {
	return &NodeGetMD5{
		name: "md5",
	}
}
