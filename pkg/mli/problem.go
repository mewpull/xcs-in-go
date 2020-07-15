// Package mli provides machine learning interfaces to relate a machine
// learning problem to a machine learning algorithm.
package mli

type Problem interface {
    IsAtEndState() bool
    Reset()
    ObtainInput() DataItem
    Effect(action int) int
}
