// Package main runs the XCS algorithm on the 6-bit Boolean multiplexer
// problem.
package main

import (
    "../pkg/xcs"
    "../pkg/multiplexer"
    "../pkg/mli"
)

func main() {
    var prob mli.Problem = multiplexer.New(6)
    var alg mli.Algorithm = &xcs.Xcs{}
    alg.OperateOn(prob)
}
