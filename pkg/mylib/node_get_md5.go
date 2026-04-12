package mylib

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type NodeGetMD5[T string] struct {
	next Noder[T]
	name string
}

func (n *NodeGetMD5[T]) GetName() string {
	return n.name
}

func (n *NodeGetMD5[T]) GetNextNode() Noder[T] {
	return n.next
}

func (n *NodeGetMD5[T]) SetNextNode(noder Noder[T]) {
	n.next = noder
}

func (n *NodeGetMD5[T]) Execute(ch chan T) (chan T, error) {
	filePath := <-ch

	file, err := os.Open(string(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to open file (%s): %w", filePath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	hash := md5.New()

	if _, err = io.Copy(hash, file); err != nil {
		return nil, fmt.Errorf("failed to calc MD5-hash: %w", err)
	}

	checksum := hex.EncodeToString(hash.Sum(nil))

	chOut := make(chan T, 1)
	chOut <- T(checksum)
	close(chOut)

	return chOut, nil
}

func NewNodeGetMD5[T string]() *NodeGetMD5[T] {
	return &NodeGetMD5[T]{
		name: "md5",
	}
}
