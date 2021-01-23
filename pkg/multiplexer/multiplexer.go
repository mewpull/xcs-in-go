// Package multiplexer provides a software implementation of a Boolean
// multiplexer -- a classic problem for learning classifier systems.
package multiplexer

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/matthewrkarlsen/xcs-in-go/pkg/mli"
)

type Multiplexer struct {
	MultiplexerSize int
	ControlBits     int
	CorrectAnswer   int
	LastAnswer      int
	EndState        bool
}

func New(multiplexerSize int) *Multiplexer {
	var controlBits = -1
	for k := 1; k < multiplexerSize; k++ {
		var maxNum = int(math.Pow(2, float64(k)))
		if k+maxNum == multiplexerSize {
			controlBits = k
		}
	}
	if controlBits == -1 {
		fmt.Println(strconv.Itoa(multiplexerSize) + " bits is not a valid multiplexer")
		os.Exit(-1)
	}
	return &Multiplexer{multiplexerSize, controlBits, -1, -1, false}
}

func (m *Multiplexer) IsAtEndState() bool {
	return m.EndState
}

func (m *Multiplexer) Reset() {
	m.EndState = false
}

func (m *Multiplexer) ObtainInput() mli.DataItem {
	var attributes = make([]int, m.MultiplexerSize)
	for j := 0; j < m.MultiplexerSize; j++ {
		attributes[j] = rand.Intn(2)
	}
	m.CorrectAnswer = m.GetMultiplexerAnswer(attributes)
	return &DataItemImpl{attributes, m.CorrectAnswer}
}

func (m *Multiplexer) Effect(action int) int {
	m.EndState = true
	if action == m.CorrectAnswer {
		return 1000
	}
	return 0
}

func (m *Multiplexer) GetMultiplexerAnswer(attributes []int) int {
	firstInt := m.ControlBits - 1
	exp := 0
	var total = 0
	for j := firstInt; j >= 0; j-- {
		var potentialValueAtByte = int(math.Pow(2, float64(exp)))
		var binaryValueAtByte = attributes[j]
		total += binaryValueAtByte * potentialValueAtByte
		exp += 1
	}
	return attributes[total]
}
