package mli

type DataItem interface {
	GetInputs() []int
	GetAnswer() int
	GetAttribute(n int) int
	ToString() string
}
