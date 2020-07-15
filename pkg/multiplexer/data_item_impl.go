// Package multiplexer provides a software implementation of a Boolean
// multiplexer -- a classic problem for learning classifier systems.
package multiplexer

import (
    "strings"
    "strconv"
)

type DataItemImpl struct {
    Inputs []int
    Answer int
}

func (d *DataItemImpl) ToString() string {


    var inputs = d.Inputs
    var builder = strings.Builder{}
    for i := 0; i < len(inputs); i++ {
        builder.WriteString(strconv.Itoa(inputs[i]))
    }
    var inpStr = builder.String()
    return inpStr + " --> " + strconv.Itoa(d.Answer)
}

func (d *DataItemImpl) GetInputs() []int {
    return d.Inputs
}

func (d *DataItemImpl) GetAnswer() int {
    return d.Answer
}

func (d *DataItemImpl) GetAttribute(n int) int {
    return d.Inputs[n]
}
