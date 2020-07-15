// Package mli provides machine learning interfaces to relate a machine
// learning problem to a machine learning algorithm.
package mli

type DataItem interface {
    GetInputs() []int
    GetAnswer() int
    GetAttribute(n int) int
    ToString() string
}
