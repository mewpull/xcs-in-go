// Package xcs provides an implementation of the eXtented Classifier
// System algorithm as described within [Butz, M. V., & Wilson, S. W.
// (2000, September). An algorithmic description of XCS. In
// International Workshop on Learning Classifier Systems (pp. 253-272).
// Springer, Berlin, Heidelberg].
package xcs

import (
	"strconv"
	"strings"
)

type Classifier struct {
	Condition       []string
	Action          int
	Payoff          float64
	InitialError    float64
	PredictionError float64
	FitnessI        float64
	Fitness         float64
	Numerosity      int32
	ActionSetSize   float64
	Exp             int64
	ThetaSub        int64
	TimeStamp       int64
	V               float64
	CorrectSets     []int
	ThetaDel        int64
	Delta           float64
	ErrorZero       float64
}

func (c *Classifier) ToString() string {
	var condition = c.Condition
	var builder = strings.Builder{}
	for i := 0; i < len(condition); i++ {
		builder.WriteString(condition[i])
	}
	var condStr = builder.String()
	return condStr + " --> " + strconv.Itoa(c.Action) +
		" [PAY:" + strconv.FormatFloat(c.Payoff, 'f', -1, 32) +
		"; ERR: " + strconv.FormatFloat(c.PredictionError, 'f', -1, 32) + "]"
}

func (c *Classifier) UpdateCorrectSetSize(setSize int) float64 {
	c.CorrectSets[len(c.CorrectSets)] = setSize
	var tot = 0
	for i := 0; i < len(c.CorrectSets); i++ {
		var ss = c.CorrectSets[i]
		tot += ss
	}
	var correctSetSize float64 = float64(tot) / float64(len(c.CorrectSets))
	return correctSetSize
}

func (c *Classifier) SetConditionComponent(n int, val string) {
	c.Condition[n] = val
}

func (c *Classifier) GetExperience() int64 {
	return c.Exp
}

func (c *Classifier) SetExperience(exp int64) {
	c.Exp = exp
}

func (c *Classifier) GetCondition() []string {
	return c.Condition
}

func (c *Classifier) GetAction() int {
	return c.Action
}

func (c *Classifier) GetFitness() float64 {
	return c.Fitness
}

func (c *Classifier) GetPayoff() float64 {
	return c.Payoff
}

func (c *Classifier) GetNumerosity() int32 {
	return c.Numerosity
}

func (c *Classifier) GetDeletionVote(averageFitnessOfPop float64) float64 {
	var f64Numerosity = float64(c.Numerosity)
	var deletionVote float64 = float64(c.ActionSetSize) * f64Numerosity
	if c.Exp > c.ThetaDel && c.Fitness/f64Numerosity < c.Delta*averageFitnessOfPop {
		deletionVote = (deletionVote * averageFitnessOfPop) / (c.Fitness / f64Numerosity)
	}
	return deletionVote
}

func (c *Classifier) IncrementNumerosityBy(num int32) {
	c.Numerosity = c.Numerosity + int32(num)
}

func (c *Classifier) DecrementNumerosity() {
	c.Numerosity -= 1
}

func (c *Classifier) IncrementNumerosity() {
	c.Numerosity += 1
}

func (c *Classifier) IncrementExperience() {
	c.Exp += 1
}

func (c *Classifier) CouldSubsume() bool {
	if c.Exp > c.ThetaSub {
		if c.PredictionError < c.ErrorZero {
			return true
		}
	}
	return false
}

func (c *Classifier) DoesSubsume(classifier *Classifier) bool {
	if c.Action == classifier.GetAction() {
		if c.CouldSubsume() {
			if c.IsMoreGeneralThan(classifier) {
				return true
			}
		}
	}
	return false
}

func (c *Classifier) DoesMatch(classifier *Classifier) bool {
	if c.GetAction() != classifier.GetAction() {
		return false
	}
	var condition1 = c.GetCondition()
	var condition2 = classifier.GetCondition()
	for i := 0; i < len(condition1); i++ {
		if condition1[i] != condition2[i] {
			return false
		}
	}
	return true
}

func (c *Classifier) IsMoreGeneralThan(classifier *Classifier) bool {
	if c.GetHashCount() <= classifier.GetHashCount() {
		return false
	}
	var condition1 = c.GetCondition()
	var condition2 = classifier.GetCondition()
	for i := 0; i < len(condition1); i++ {
		if condition1[i] != "#" && condition1[i] != condition2[i] {
			return false
		}
	}
	return true
}

func (c *Classifier) GetError() float64 {
	return c.PredictionError
}

func (c *Classifier) GetHashCount() int {
	var hashCount = 0
	for _, x := range c.Condition {
		if x == "#" {
			hashCount += 1
		}
	}
	return hashCount
}

func (c *Classifier) GetTimeStamp() int64 {
	return c.TimeStamp
}

func (c *Classifier) SetPayoff(payoff float64) {
	c.Payoff = payoff
}

func (c *Classifier) SetActionSetSize(actionSetSize float64) {
	c.ActionSetSize = actionSetSize
}

func (c *Classifier) GetActionSetSize() float64 {
	return float64(c.ActionSetSize)
}

func (c *Classifier) GetPredictionError() float64 {
	return c.PredictionError
}

func (c *Classifier) SetPredictionError(predictionError float64) {
	c.PredictionError = predictionError
}

func (c *Classifier) SetTimeStamp(timeStamp int64) {
	c.TimeStamp = timeStamp
}

func (c *Classifier) GetOffspring() *Classifier {
	var condition = make([]string, len(c.Condition))
	for i, a := range c.Condition {
		condition[i] = a
	}
	var cl = Classifier{condition, c.Action, 0.0, c.InitialError, c.InitialError, c.FitnessI, c.FitnessI, 1, 0, 0, c.ThetaSub, c.TimeStamp, c.V, make([]int, 80000), c.ThetaDel, c.Delta, c.ErrorZero}
	cl.SetFitness(c.Fitness)
	cl.SetPayoff(c.Payoff)
	cl.SetPredictionError(c.PredictionError)
	cl.SetExperience(0)
	cl.SetNumerosity(1)
	cl.SetActionSetSize(c.ActionSetSize)
	cl.SetPayoff(c.Payoff)
	return &cl
}

func (c *Classifier) SetFitness(fitness float64) {
	c.Fitness = fitness
}

func (c *Classifier) SetAction(action int) {
	c.Action = action
}

func (c *Classifier) SetError(predictionError float64) {
	c.PredictionError = predictionError
}

func (c *Classifier) SetNumerosity(numerosity int32) {
	c.Numerosity = numerosity
}
