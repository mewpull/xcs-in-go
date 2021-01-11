// Package main runs the XCS algorithm on the 6-bit Boolean multiplexer
// problem.
package main

import (
	"github.com/matthewrkarlsen/xcs-in-go/pkg/mli"
	"github.com/matthewrkarlsen/xcs-in-go/pkg/multiplexer"
	"github.com/matthewrkarlsen/xcs-in-go/pkg/xcs"
)

func main() {
	var prob mli.Problem = multiplexer.New(6)
	var alg mli.Algorithm = &xcs.Xcs{}
	alg.OperateOn(prob)
}
