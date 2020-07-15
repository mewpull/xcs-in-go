// Package mli provides machine learning interfaces to relate a machine
// learning problem to a machine learning algorithm.
package mli

type Algorithm interface {
    OperateOn(problem Problem)
}
