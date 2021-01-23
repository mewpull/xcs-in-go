// Package main runs the XCS algorithm on the 6-bit Boolean multiplexer
// problem.
package main

import (
	"github.com/matthewrkarlsen/xcs-in-go/pkg/multiplexer"
	"github.com/matthewrkarlsen/xcs-in-go/pkg/xcs"
)

func main() {
	prob := multiplexer.New(6)
	alg := &xcs.Xcs{}
	alg.OperateOn(prob)
}
