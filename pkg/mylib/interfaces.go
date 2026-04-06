package mylib

type Noder interface {
	GetName() string
	GetNextNode() Noder
	SetNextNode(n Noder)
	Execute(filePath string) (string, error)
}
