// Package xcs provides an implementation of the eXtented Classifier
// System algorithm as described within [Butz, M. V., & Wilson, S. W.
// (2000, September). An algorithmic description of XCS. In
// International Workshop on Learning Classifier Systems (pp. 253-272).
// Springer, Berlin, Heidelberg].
package xcs

import (
	"container/list"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/matthewrkarlsen/xcs-in-go/pkg/mli"
)

var beta float64 = 0.2
var pHash float64 = 0.33
var nu float64 = 5.0
var pExplore float64 = 0.5
var maxPop int32 = 400
var doGaSubsumption = true
var thetaGa int64 = 50
var chi float64 = 0.8
var mu float64 = 0.04
var maxAction int = 1
var errorZero float64 = 10
var thetaSub int64 = 20
var thetaMna int = maxAction + 1
var epsilon0 float64 = 0.001
var alpha float64 = 0.1
var thetaDel int64 = 20
var delta float64 = 0.1
var gamma float64 = 0.71
var actionSetSubsumption = false
var fitnessI float64 = 0.0
var initialError float64 = 0.0

type Xcs struct {
}

func (x *Xcs) RuleMatchesState(rule *Classifier, state mli.DataItem) bool {
	var condition = rule.GetCondition()
	var inputs = state.GetInputs()
	var numAttributes = len(inputs)
	for idx := 0; idx < numAttributes; idx++ {
		if condition[idx] == "0" {
			if inputs[idx] != 0 {
				return false
			}
		} else if condition[idx] == "1" {
			if inputs[idx] != 1 {
				return false
			}
		}
	}
	return true
}

func (x *Xcs) CreateMatchSet(ruleSet *list.List, dataItem mli.DataItem, step int64) *list.List {
	var matchSet = x.ObtainMatchingClassifiers(ruleSet, dataItem)
	for matchSet.Len() < thetaMna {
		var cl *Classifier = x.GenerateClassifier(matchSet, dataItem, step)
		ruleSet.PushBack(cl)
		x.DeleteFromPop(ruleSet)
		matchSet.PushBack(cl)
	}
	return matchSet
}

func (x *Xcs) ObtainMatchingClassifiers(ruleSet *list.List, dataItem mli.DataItem) *list.List {
	var matchSet = list.New()
	for r := ruleSet.Front(); r != nil; r = r.Next() {
		if x.RuleMatchesState(r.Value.(*Classifier), dataItem) {
			matchSet.PushBack(r.Value.(*Classifier))
		}
	}
	return matchSet
}

func (x *Xcs) GenerateClassifier(matchSet *list.List, dataItem mli.DataItem, step int64) *Classifier {
	condition := make([]string, len(dataItem.GetInputs()))
	for i, attrib := range dataItem.GetInputs() {
		if rand.Float64() < pHash {
			condition[i] = "#"
		} else {
			condition[i] = strconv.Itoa(attrib)
		}
	}
	var actionsPresent = make(map[int]bool, maxAction-1)
	var allActions = list.New()
	for i := 0; i <= maxAction; i++ {
		allActions.PushBack(i)
	}
	for i := 0; i <= maxAction; i++ {
		actionsPresent[i] = false
	}

	for e := matchSet.Front(); e != nil; e = e.Next() {
		var cl = e.Value.(*Classifier)
		var act = cl.GetAction()
		actionsPresent[act] = true
	}

	var toChooseFrom = list.New()
	for i := 0; i <= maxAction; i++ {
		if actionsPresent[i] == false {
			toChooseFrom.PushBack(i)
		}
	}

	if toChooseFrom.Len() == 0 {
		toChooseFrom = allActions
	}

	var randomIdx = rand.Intn(toChooseFrom.Len())
	var currentIdx = 0
	var answer = -1
	for e := toChooseFrom.Front(); e != nil; e = e.Next() {
		if currentIdx == randomIdx {
			answer = e.Value.(int)
			break
		}
		currentIdx++
	}
	if answer == -1 {
		fmt.Println("Error. answer == -1.")
		os.Exit(-1)
	}
	return &Classifier{condition, answer, 0.0, initialError, initialError, fitnessI, fitnessI, 1, 0, 0, thetaSub, step, nu, make([]int, 80000), thetaDel, delta, errorZero}
}

func (x *Xcs) CountMicroClassifiers(ruleSet *list.List) int32 {
	var microPop int32 = 0
	for e := ruleSet.Front(); e != nil; e = e.Next() {
		var cls = e.Value.(*Classifier)
		microPop += cls.GetNumerosity()
	}
	return microPop
}

func (x *Xcs) GetAverageFitnessOfPop(ruleSet *list.List, microPopCount int32) float64 {
	var fitnessSum = 0.0
	for e := ruleSet.Front(); e != nil; e = e.Next() {
		var cls = e.Value.(*Classifier)
		fitnessSum += cls.GetFitness()
	}
	return fitnessSum / float64(microPopCount)
}

func (x *Xcs) DeleteFromPop(ruleSet *list.List) {
	var microPop = x.CountMicroClassifiers(ruleSet)
	if microPop < maxPop {
		return
	}
	var averageFitnessOfPop float64 = x.GetAverageFitnessOfPop(ruleSet, microPop)
	var voteSum float64 = 0.0
	for e := ruleSet.Front(); e != nil; e = e.Next() {
		var cls = e.Value.(*Classifier)
		voteSum += cls.GetDeletionVote(averageFitnessOfPop)
	}
	var choicePoint = voteSum * rand.Float64()
	voteSum = 0.0
	for e := ruleSet.Front(); e != nil; e = e.Next() {
		var cl = e.Value.(*Classifier)
		voteSum += cl.GetDeletionVote(averageFitnessOfPop)
		if voteSum > choicePoint {
			if cl.GetNumerosity() > 1 {
				cl.DecrementNumerosity()
			} else {
				ruleSet.Remove(e)
			}
			break
		}
	}
}

func (x *Xcs) DoActionSetSubsumption(actionSet *list.List, ruleSet *list.List) {
	var cl *Classifier
	for e := actionSet.Front(); e != nil; e = e.Next() {
		var c = e.Value.(*Classifier)
		if c.CouldSubsume() {
			if cl == nil || c.GetHashCount() > cl.GetHashCount() || (c.GetHashCount() == cl.GetHashCount() && rand.Float64() > 0.5) {
				cl = c
			}
		}
	}

	var toDelete = list.New()
	if cl != nil {
		for e := actionSet.Front(); e != nil; e = e.Next() {
			var classifier = e.Value.(*Classifier)
			if cl.IsMoreGeneralThan(classifier) {
				toDelete.PushBack(e)
			}
		}
	}
	for e := toDelete.Front(); e != nil; e = e.Next() {
		ruleSet.Remove(e)
		actionSet.Remove(e)
		var td = e.Value.(*Classifier)
		cl.IncrementNumerosityBy(td.GetNumerosity())
	}
}

func (x *Xcs) ApplyMutation(classifier *Classifier, dataItem mli.DataItem) {
	var condition = classifier.GetCondition()
	for k := 0; k < len(condition); k++ {
		if rand.Float64() < mu {
			if condition[k] == "#" {
				var attrib int = dataItem.GetAttribute(k)
				classifier.SetConditionComponent(k, strconv.Itoa(attrib))
			} else {
				classifier.SetConditionComponent(k, "#")
			}
		}
	}
	if rand.Float64() < mu {
		var clAction = classifier.GetAction()
		var allActions = x.GetSetOfActionsLessSpecified(clAction)
		var randomIdx = rand.Intn(allActions.Len())
		var currentIdx = 0
		var action = -1
		for e := allActions.Front(); e != nil; e = e.Next() {
			if currentIdx == randomIdx {
				action = e.Value.(int)
				break
			}
			currentIdx++
		}
		if action == -1 {
			fmt.Println("Error. action == -1.")
			os.Exit(-1)
		}
		classifier.SetAction(action)
	}
}

func (x *Xcs) GetSetOfActionsLessSpecified(action int) *list.List {
	var allActions = list.New()
	for j := 0; j < thetaMna; j++ {
		if j != action {
			allActions.PushBack(j)
		}
	}
	return allActions
}

func (x *Xcs) InsertInPopulation(classifier *Classifier, ruleSet *list.List) {
	for e := ruleSet.Front(); e != nil; e = e.Next() {
		var cl = e.Value.(*Classifier)
		if cl.DoesMatch(classifier) {
			cl.IncrementNumerosity()
			return
		}
	}
	ruleSet.PushBack(classifier)
}

func (xs *Xcs) ApplyCrossover(classifier1 *Classifier, classifier2 *Classifier) {
	var x = rand.Intn(len(classifier1.GetCondition()))
	var y = rand.Intn(len(classifier2.GetCondition()))
	if x > y {
		var z = x
		x = y
		y = z
	}
	var condition1 = classifier1.GetCondition()
	var condition2 = classifier2.GetCondition()

	for m := 0; m < len(classifier1.GetCondition()); m++ {
		if x <= m && m < y {
			var cc1 = condition1[m]
			var cc2 = condition2[m]
			condition1[m] = cc2
			condition2[m] = cc1
		}
	}

	var newFitness = (classifier1.GetFitness() + classifier2.GetFitness()) / 2
	classifier1.SetFitness(newFitness)
	classifier2.SetFitness(newFitness)

	var newError = (classifier1.GetError() + classifier2.GetError()) / 2
	classifier1.SetError(newError)
	classifier2.SetError(newError)

	var newPayoff = (classifier1.GetPayoff() + classifier2.GetPayoff()) / 2
	classifier1.SetPayoff(newPayoff)
	classifier2.SetPayoff(newPayoff)
}

func (x *Xcs) SelectOffspring(actionSet *list.List) *Classifier {
	if actionSet.Len() == 1 {
		return actionSet.Front().Value.(*Classifier)
	}
	var fitnessSum = 0.0
	for e := actionSet.Front(); e != nil; e = e.Next() {
		var cl = e.Value.(*Classifier)
		fitnessSum = fitnessSum + cl.GetFitness()
	}
	var choicePoint = rand.Float64() * fitnessSum
	fitnessSum = 0.0
	for e := actionSet.Front(); e != nil; e = e.Next() {
		var cl = e.Value.(*Classifier)
		fitnessSum = fitnessSum + cl.GetFitness()
		if fitnessSum > choicePoint {
			return cl
		}
	}

	var randomIdx = rand.Intn(actionSet.Len())
	var currentIdx = 0
	var cls *Classifier
	for e := actionSet.Front(); e != nil; e = e.Next() {
		if currentIdx == randomIdx {
			cls = e.Value.(*Classifier)
			break
		}
		currentIdx++
	}
	if cls == nil {
		fmt.Println("Error. cls == nil.")
		os.Exit(-1)
	}
	return cls
}

func (x *Xcs) RunGeneticAlgorithm(actionSet *list.List, dataItem mli.DataItem, ruleSet *list.List, step int64) {
	var numerositySum int32 = 0
	var timeStampSum int64 = 0
	for e := actionSet.Front(); e != nil; e = e.Next() {
		var cl = e.Value.(*Classifier)
		numerositySum = numerositySum + cl.GetNumerosity()
		timeStampSum = timeStampSum + (cl.GetTimeStamp() * int64(cl.GetNumerosity()))
	}
	if float64(step)-float64(timeStampSum)/float64(numerositySum) > float64(thetaGa) {

		for e := actionSet.Front(); e != nil; e = e.Next() {
			var cl = e.Value.(*Classifier)
			cl.SetTimeStamp(step)
		}

		var parent1 = x.SelectOffspring(actionSet)
		var parent2 = x.SelectOffspring(actionSet)
		var child1 = parent1.GetOffspring()
		var child2 = parent2.GetOffspring()

		if rand.Float64() < chi {
			x.ApplyCrossover(child1, child2)
		}

		child1.SetFitness(child1.GetFitness() * 0.1)
		child2.SetFitness(child2.GetFitness() * 0.1)

		for _, child := range []*Classifier{child1, child2} {
			x.ApplyMutation(child, dataItem)
			if doGaSubsumption {
				if parent1.DoesSubsume(child) {
					parent1.IncrementNumerosity()
				} else if parent2.DoesSubsume(child) {
					parent2.IncrementNumerosity()
				} else {
					x.InsertInPopulation(child, ruleSet)
				}
			} else {
				x.InsertInPopulation(child, ruleSet)
			}
			x.DeleteFromPop(ruleSet)
		}
	}
}

func (x *Xcs) CreatePredictionArray(matchSet *list.List) map[int]float64 {
	var actionSet = make(map[int]bool, maxAction)
	for e := matchSet.Front(); e != nil; e = e.Next() {
		var cl = e.Value.(*Classifier)
		actionSet[cl.GetAction()] = true
	}
	var pa = make(map[int]float64, len(actionSet))
	var fsa = make(map[int]float64, len(actionSet))
	for k := range actionSet {
		fsa[k] = 0.0
	}
	for e := matchSet.Front(); e != nil; e = e.Next() {
		var cl = e.Value.(*Classifier)
		var a = cl.GetAction()
		var _, exists = pa[a]
		if !exists {
			pa[a] = cl.GetPayoff() * cl.GetFitness()
		} else {
			pa[a] = pa[a] + cl.GetPayoff()*cl.GetFitness()
		}
		fsa[a] = fsa[a] + cl.GetFitness()
	}
	for a := range actionSet {
		if fsa[a] > 0.0 {
			pa[a] = pa[a] / fsa[a]
		}
	}
	return pa
}

func (x *Xcs) CreateActionSet(matchSet *list.List, action int) *list.List {
	var actionSet = list.New()
	for e := matchSet.Front(); e != nil; e = e.Next() {
		var cl *Classifier = e.Value.(*Classifier)
		if cl.GetAction() == action {
			actionSet.PushBack(cl)
		}
	}
	return actionSet
}

func (x *Xcs) UpdateFitnessInSet(actionSet *list.List) {
	var accuracySum float64 = 0.0
	var k = make([]float64, actionSet.Len())
	var i = 0
	for e := actionSet.Front(); e != nil; e = e.Next() {
		var cl = e.Value.(*Classifier)
		var val float64
		if cl.GetError() < errorZero {
			val = 1.0
		} else {
			val = alpha * (math.Pow((cl.GetError() / errorZero), -nu))
		}
		k[i] = val
		accuracySum += val * float64(cl.GetNumerosity())
		i++
	}
	i = 0
	for e := actionSet.Front(); e != nil; e = e.Next() {
		var cl = e.Value.(*Classifier)
		var fitness = cl.GetFitness() + beta*(k[i]*float64(cl.GetNumerosity())/accuracySum-cl.GetFitness())
		cl.SetFitness(fitness)
		i++
	}
}

func (x *Xcs) UpdateActionSet(capitalP float64, actionSet *list.List, ruleSet *list.List) {
	var totAsNum = x.CountMicroClassifiers(actionSet)
	for e := actionSet.Front(); e != nil; e = e.Next() {
		var cl = e.Value.(*Classifier)
		cl.IncrementExperience()
		var classifierExp = cl.GetExperience()
		var inexperienced = float64(classifierExp) < 1/beta
		var payoff = cl.GetPayoff()
		if inexperienced {
			payoff += (capitalP - payoff) / float64(classifierExp)
		} else {
			payoff += beta * (capitalP - payoff)
		}
		cl.SetPayoff(payoff)
		var predictionError = cl.GetPredictionError()
		if inexperienced {
			predictionError += (math.Abs(capitalP-payoff) - predictionError) / float64(classifierExp)
		} else {
			predictionError += beta * (math.Abs(capitalP-payoff) - predictionError)
		}
		cl.SetPredictionError(predictionError)
		var actionSetSize = cl.GetActionSetSize()
		if inexperienced {
			actionSetSize += (float64(totAsNum) - actionSetSize) / float64(classifierExp)
		} else {
			actionSetSize += beta * (float64(totAsNum) - actionSetSize)
		}
		cl.SetActionSetSize(actionSetSize)
	}
	x.UpdateFitnessInSet(actionSet)
	if actionSetSubsumption {
		x.DoActionSetSubsumption(actionSet, ruleSet)
	}
}

func (x *Xcs) Evaluate(problem mli.Problem, ruleSet *list.List, macroStep int) {
	var numCorrect int = 0
	var numIncorrect int = 0
	for j := 0; j < 100; j++ {
		problem.Reset()
		for problem.IsAtEndState() == false {
			var dataItem = problem.ObtainInput()
			var matchSet = x.ObtainMatchingClassifiers(ruleSet, dataItem)
			if matchSet.Len() == 0 {
				numIncorrect += 1
				break
			}
			var predictionArray = x.CreatePredictionArray(matchSet)
			var bestAction = 0
			var expP float64 = 0.0
			for key, expPTmp := range predictionArray {
				if expPTmp > expP {
					expP = expPTmp
					bestAction = key
				}
			}
			var reward = problem.Effect(bestAction)
			if problem.IsAtEndState() {
				if reward == 1000 {
					numCorrect += 1
				} else {
					numIncorrect += 1
				}
			}
		}
	}
	if (numCorrect + numIncorrect) > 0 {
		var propCorrect = (float64(numCorrect) / (float64(numCorrect) + float64(numIncorrect)))
		fmt.Println("Post-cycle eval #", macroStep, ". Proportion correct: ", propCorrect)
	}
}

func (x *Xcs) OperateOn(problem mli.Problem) {

	ruleSet := list.New()

	var cumulativeMicroSteps int64 = 0
	macroStep := 0
	for i := 0; i < 80001; i++ {
		microStep := 0
		problem.Reset()
		var lastActionSet *list.List
		var lastReward int
		for problem.IsAtEndState() == false {
			var dataItem = problem.ObtainInput()
			var matchSet = x.CreateMatchSet(ruleSet, dataItem, cumulativeMicroSteps)
			var predictionArray = x.CreatePredictionArray(matchSet)
			var paRandomIdx = rand.Intn(len(predictionArray))
			var idx = 0
			var bestAction = 0
			for k := range predictionArray {
				if idx == paRandomIdx {
					bestAction = k
					break
				}
				idx++
			}
			var expP float64 = 0.0
			if rand.Float64() > pExplore {
				for k, val := range predictionArray {
					if val > expP {
						expP = val
						bestAction = k
					}
				}
			} else {
				expP = predictionArray[bestAction]
			}
			var actionSet = x.CreateActionSet(matchSet, bestAction)
			var reward = problem.Effect(bestAction)
			if lastActionSet != nil {
				var capitalP = float64(lastReward) + gamma*expP
				x.UpdateActionSet(capitalP, lastActionSet, ruleSet)
				if cumulativeMicroSteps%thetaGa == 0 {
					x.RunGeneticAlgorithm(lastActionSet, dataItem, ruleSet, cumulativeMicroSteps)
				}
			}
			if problem.IsAtEndState() {
				var capitalP = float64(reward)
				x.UpdateActionSet(capitalP, actionSet, ruleSet)
				if cumulativeMicroSteps%thetaGa == 0 {
					x.RunGeneticAlgorithm(actionSet, dataItem, ruleSet, cumulativeMicroSteps)
				}
				lastActionSet = nil
				lastReward = -1.0
			} else {
				lastActionSet = actionSet
				lastReward = reward
			}
			microStep += 1
		}
		if i != 0 && i%50 == 0 {
			x.Evaluate(problem, ruleSet, macroStep)
		}
		macroStep += 1
		cumulativeMicroSteps += 1
	}

	for e := ruleSet.Front(); e != nil; e = e.Next() {
		var cl = e.Value.(*Classifier)
		fmt.Println(cl.ToString())
	}
}
