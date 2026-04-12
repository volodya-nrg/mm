package mylib

type Noder[T string] interface {
	GetName() string
	GetNextNode() Noder[T]
	SetNextNode(n Noder[T])
	Execute(chan T) (chan T, error)
}
