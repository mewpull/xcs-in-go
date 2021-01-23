package mli

type Problem interface {
	IsAtEndState() bool
	Reset()
	ObtainInput() DataItem
	Effect(action int) int
}
