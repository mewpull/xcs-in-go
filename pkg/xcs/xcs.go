// Package xcs provides an implementation of the eXtented Classifier
// System algorithm as described within [Butz, M. V., & Wilson, S. W.
// (2000, September). An algorithmic description of XCS. In
// International Workshop on Learning Classifier Systems (pp. 253-272).
// Springer, Berlin, Heidelberg].
package xcs

import (
	"container/list"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"

	"github.com/matthewrkarlsen/xcs-in-go/pkg/mli"
)

const (
	beta                 = 0.2
	pHash                = 0.33
	nu                   = 5.0
	pExplore             = 0.5
	maxPop               = 400
	doGaSubsumption      = true
	thetaGa              = 50
	chi                  = 0.8
	mu                   = 0.04
	maxAction            = 1
	errorZero            = 10
	thetaSub             = 20
	thetaMna             = maxAction + 1
	epsilon0             = 0.001
	alpha                = 0.1
	thetaDel             = 20
	delta                = 0.1
	gamma                = 0.71
	actionSetSubsumption = false
	fitnessI             = 0.0
	initialError         = 0.0
)

type Xcs struct {
}

func (x *Xcs) RuleMatchesState(rule *Classifier, state mli.DataItem) bool {
	condition := rule.GetCondition()
	inputs := state.GetInputs()
	numAttributes := len(inputs)
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
	matchSet := x.ObtainMatchingClassifiers(ruleSet, dataItem)
	for matchSet.Len() < thetaMna {
		cl := x.GenerateClassifier(matchSet, dataItem, step)
		ruleSet.PushBack(cl)
		x.DeleteFromPop(ruleSet)
		matchSet.PushBack(cl)
	}
	return matchSet
}

func (x *Xcs) ObtainMatchingClassifiers(ruleSet *list.List, dataItem mli.DataItem) *list.List {
	matchSet := list.New()
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
	actionsPresent := make(map[int]bool, maxAction-1)
	allActions := list.New()
	for i := 0; i <= maxAction; i++ {
		allActions.PushBack(i)
	}
	for i := 0; i <= maxAction; i++ {
		actionsPresent[i] = false
	}

	for e := matchSet.Front(); e != nil; e = e.Next() {
		cl := e.Value.(*Classifier)
		act := cl.GetAction()
		actionsPresent[act] = true
	}

	toChooseFrom := list.New()
	for i := 0; i <= maxAction; i++ {
		if actionsPresent[i] == false {
			toChooseFrom.PushBack(i)
		}
	}

	if toChooseFrom.Len() == 0 {
		toChooseFrom = allActions
	}

	randomIdx := rand.Intn(toChooseFrom.Len())
	currentIdx := 0
	answer := -1
	for e := toChooseFrom.Front(); e != nil; e = e.Next() {
		if currentIdx == randomIdx {
			answer = e.Value.(int)
			break
		}
		currentIdx++
	}
	if answer == -1 {
		log.Fatal("Error. answer == -1.")
	}
	return &Classifier{condition, answer, 0.0, initialError, initialError, fitnessI, fitnessI, 1, 0, 0, thetaSub, step, nu, make([]int, 80000), thetaDel, delta, errorZero}
}

func (x *Xcs) CountMicroClassifiers(ruleSet *list.List) int32 {
	microPop := int32(0)
	for e := ruleSet.Front(); e != nil; e = e.Next() {
		cls := e.Value.(*Classifier)
		microPop += cls.GetNumerosity()
	}
	return microPop
}

func (x *Xcs) GetAverageFitnessOfPop(ruleSet *list.List, microPopCount int32) float64 {
	fitnessSum := 0.0
	for e := ruleSet.Front(); e != nil; e = e.Next() {
		cls := e.Value.(*Classifier)
		fitnessSum += cls.GetFitness()
	}
	return fitnessSum / float64(microPopCount)
}

func (x *Xcs) DeleteFromPop(ruleSet *list.List) {
	microPop := x.CountMicroClassifiers(ruleSet)
	if microPop < maxPop {
		return
	}
	averageFitnessOfPop := x.GetAverageFitnessOfPop(ruleSet, microPop)
	voteSum := 0.0
	for e := ruleSet.Front(); e != nil; e = e.Next() {
		cls := e.Value.(*Classifier)
		voteSum += cls.GetDeletionVote(averageFitnessOfPop)
	}
	choicePoint := voteSum * rand.Float64()
	voteSum = 0.0
	for e := ruleSet.Front(); e != nil; e = e.Next() {
		cl := e.Value.(*Classifier)
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
		c := e.Value.(*Classifier)
		if c.CouldSubsume() {
			if cl == nil || c.GetHashCount() > cl.GetHashCount() || (c.GetHashCount() == cl.GetHashCount() && rand.Float64() > 0.5) {
				cl = c
			}
		}
	}

	toDelete := list.New()
	if cl != nil {
		for e := actionSet.Front(); e != nil; e = e.Next() {
			classifier := e.Value.(*Classifier)
			if cl.IsMoreGeneralThan(classifier) {
				toDelete.PushBack(e)
			}
		}
	}
	for e := toDelete.Front(); e != nil; e = e.Next() {
		ruleSet.Remove(e)
		actionSet.Remove(e)
		td := e.Value.(*Classifier)
		cl.IncrementNumerosityBy(td.GetNumerosity())
	}
}

func (x *Xcs) ApplyMutation(classifier *Classifier, dataItem mli.DataItem) {
	condition := classifier.GetCondition()
	for k := 0; k < len(condition); k++ {
		if rand.Float64() < mu {
			if condition[k] == "#" {
				attrib := dataItem.GetAttribute(k)
				classifier.SetConditionComponent(k, strconv.Itoa(attrib))
			} else {
				classifier.SetConditionComponent(k, "#")
			}
		}
	}
	if rand.Float64() < mu {
		clAction := classifier.GetAction()
		allActions := x.GetSetOfActionsLessSpecified(clAction)
		randomIdx := rand.Intn(allActions.Len())
		currentIdx := 0
		action := -1
		for e := allActions.Front(); e != nil; e = e.Next() {
			if currentIdx == randomIdx {
				action = e.Value.(int)
				break
			}
			currentIdx++
		}
		if action == -1 {
			log.Fatal("Error. action == -1.")
		}
		classifier.SetAction(action)
	}
}

func (x *Xcs) GetSetOfActionsLessSpecified(action int) *list.List {
	allActions := list.New()
	for j := 0; j < thetaMna; j++ {
		if j != action {
			allActions.PushBack(j)
		}
	}
	return allActions
}

func (x *Xcs) InsertInPopulation(classifier *Classifier, ruleSet *list.List) {
	for e := ruleSet.Front(); e != nil; e = e.Next() {
		cl := e.Value.(*Classifier)
		if cl.DoesMatch(classifier) {
			cl.IncrementNumerosity()
			return
		}
	}
	ruleSet.PushBack(classifier)
}

func (xs *Xcs) ApplyCrossover(classifier1 *Classifier, classifier2 *Classifier) {
	x := rand.Intn(len(classifier1.GetCondition()))
	y := rand.Intn(len(classifier2.GetCondition()))
	if x > y {
		z := x
		x = y
		y = z
	}
	condition1 := classifier1.GetCondition()
	condition2 := classifier2.GetCondition()

	for m := 0; m < len(classifier1.GetCondition()); m++ {
		if x <= m && m < y {
			cc1 := condition1[m]
			cc2 := condition2[m]
			condition1[m] = cc2
			condition2[m] = cc1
		}
	}

	newFitness := (classifier1.GetFitness() + classifier2.GetFitness()) / 2
	classifier1.SetFitness(newFitness)
	classifier2.SetFitness(newFitness)

	newError := (classifier1.GetError() + classifier2.GetError()) / 2
	classifier1.SetError(newError)
	classifier2.SetError(newError)

	newPayoff := (classifier1.GetPayoff() + classifier2.GetPayoff()) / 2
	classifier1.SetPayoff(newPayoff)
	classifier2.SetPayoff(newPayoff)
}

func (x *Xcs) SelectOffspring(actionSet *list.List) *Classifier {
	if actionSet.Len() == 1 {
		return actionSet.Front().Value.(*Classifier)
	}
	fitnessSum := 0.0
	for e := actionSet.Front(); e != nil; e = e.Next() {
		cl := e.Value.(*Classifier)
		fitnessSum = fitnessSum + cl.GetFitness()
	}
	choicePoint := rand.Float64() * fitnessSum
	fitnessSum = 0.0
	for e := actionSet.Front(); e != nil; e = e.Next() {
		cl := e.Value.(*Classifier)
		fitnessSum = fitnessSum + cl.GetFitness()
		if fitnessSum > choicePoint {
			return cl
		}
	}

	randomIdx := rand.Intn(actionSet.Len())
	currentIdx := 0
	var cls *Classifier
	for e := actionSet.Front(); e != nil; e = e.Next() {
		if currentIdx == randomIdx {
			cls = e.Value.(*Classifier)
			break
		}
		currentIdx++
	}
	if cls == nil {
		log.Fatal("Error. cls == nil.")
	}
	return cls
}

func (x *Xcs) RunGeneticAlgorithm(actionSet *list.List, dataItem mli.DataItem, ruleSet *list.List, step int64) {
	numerositySum := int32(0)
	timeStampSum := int64(0)
	for e := actionSet.Front(); e != nil; e = e.Next() {
		cl := e.Value.(*Classifier)
		numerositySum = numerositySum + cl.GetNumerosity()
		timeStampSum = timeStampSum + (cl.GetTimeStamp() * int64(cl.GetNumerosity()))
	}
	if float64(step)-float64(timeStampSum)/float64(numerositySum) > float64(thetaGa) {

		for e := actionSet.Front(); e != nil; e = e.Next() {
			cl := e.Value.(*Classifier)
			cl.SetTimeStamp(step)
		}

		parent1 := x.SelectOffspring(actionSet)
		parent2 := x.SelectOffspring(actionSet)
		child1 := parent1.GetOffspring()
		child2 := parent2.GetOffspring()

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
	actionSet := make(map[int]bool, maxAction)
	for e := matchSet.Front(); e != nil; e = e.Next() {
		cl := e.Value.(*Classifier)
		actionSet[cl.GetAction()] = true
	}
	pa := make(map[int]float64, len(actionSet))
	fsa := make(map[int]float64, len(actionSet))
	for k := range actionSet {
		fsa[k] = 0.0
	}
	for e := matchSet.Front(); e != nil; e = e.Next() {
		cl := e.Value.(*Classifier)
		a := cl.GetAction()
		_, exists := pa[a]
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
	actionSet := list.New()
	for e := matchSet.Front(); e != nil; e = e.Next() {
		cl := e.Value.(*Classifier)
		if cl.GetAction() == action {
			actionSet.PushBack(cl)
		}
	}
	return actionSet
}

func (x *Xcs) UpdateFitnessInSet(actionSet *list.List) {
	accuracySum := 0.0
	k := make([]float64, actionSet.Len())
	i := 0
	for e := actionSet.Front(); e != nil; e = e.Next() {
		cl := e.Value.(*Classifier)
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
		cl := e.Value.(*Classifier)
		fitness := cl.GetFitness() + beta*(k[i]*float64(cl.GetNumerosity())/accuracySum-cl.GetFitness())
		cl.SetFitness(fitness)
		i++
	}
}

func (x *Xcs) UpdateActionSet(capitalP float64, actionSet *list.List, ruleSet *list.List) {
	totAsNum := x.CountMicroClassifiers(actionSet)
	for e := actionSet.Front(); e != nil; e = e.Next() {
		cl := e.Value.(*Classifier)
		cl.IncrementExperience()
		classifierExp := cl.GetExperience()
		inexperienced := float64(classifierExp) < 1/beta
		payoff := cl.GetPayoff()
		if inexperienced {
			payoff += (capitalP - payoff) / float64(classifierExp)
		} else {
			payoff += beta * (capitalP - payoff)
		}
		cl.SetPayoff(payoff)
		predictionError := cl.GetPredictionError()
		if inexperienced {
			predictionError += (math.Abs(capitalP-payoff) - predictionError) / float64(classifierExp)
		} else {
			predictionError += beta * (math.Abs(capitalP-payoff) - predictionError)
		}
		cl.SetPredictionError(predictionError)
		actionSetSize := cl.GetActionSetSize()
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
	numCorrect := 0
	numIncorrect := 0
	for j := 0; j < 100; j++ {
		problem.Reset()
		for problem.IsAtEndState() == false {
			dataItem := problem.ObtainInput()
			matchSet := x.ObtainMatchingClassifiers(ruleSet, dataItem)
			if matchSet.Len() == 0 {
				numIncorrect += 1
				break
			}
			predictionArray := x.CreatePredictionArray(matchSet)
			bestAction := 0
			expP := 0.0
			for key, expPTmp := range predictionArray {
				if expPTmp > expP {
					expP = expPTmp
					bestAction = key
				}
			}
			reward := problem.Effect(bestAction)
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
		propCorrect := (float64(numCorrect) / (float64(numCorrect) + float64(numIncorrect)))
		fmt.Printf("Post-cycle eval #%v. Proportion correct: %v\n", macroStep, propCorrect)
	}
}

func (x *Xcs) OperateOn(problem mli.Problem) {

	ruleSet := list.New()

	cumulativeMicroSteps := int64(0)
	macroStep := 0
	for i := 0; i < 80001; i++ {
		microStep := 0
		problem.Reset()
		var lastActionSet *list.List
		var lastReward int
		for problem.IsAtEndState() == false {
			dataItem := problem.ObtainInput()
			matchSet := x.CreateMatchSet(ruleSet, dataItem, cumulativeMicroSteps)
			predictionArray := x.CreatePredictionArray(matchSet)
			paRandomIdx := rand.Intn(len(predictionArray))
			idx := 0
			bestAction := 0
			for k := range predictionArray {
				if idx == paRandomIdx {
					bestAction = k
					break
				}
				idx++
			}
			expP := 0.0
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
			actionSet := x.CreateActionSet(matchSet, bestAction)
			reward := problem.Effect(bestAction)
			if lastActionSet != nil {
				capitalP := float64(lastReward) + gamma*expP
				x.UpdateActionSet(capitalP, lastActionSet, ruleSet)
				if cumulativeMicroSteps%thetaGa == 0 {
					x.RunGeneticAlgorithm(lastActionSet, dataItem, ruleSet, cumulativeMicroSteps)
				}
			}
			if problem.IsAtEndState() {
				capitalP := float64(reward)
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
		cl := e.Value.(*Classifier)
		fmt.Println(cl.ToString())
	}
}
