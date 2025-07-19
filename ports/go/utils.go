package c3

import (
	"github.com/antlr4-go/antlr/v4"
)

func GetTransitionTarget(transition antlr.Transition) antlr.ATNState {
	return transition.GetTarget()
}

func GetTransitionFollowState(transition *antlr.RuleTransition) antlr.ATNState {
	return transition.FollowState
}

func GetEpsilon(transition antlr.Transition) bool {
	return transition.GetIsEpsilon()
}

func GetLabel(transition antlr.Transition) *antlr.IntervalSet {
	return transition.GetLabel()
}

func ContainsSymbol(intervals []antlr.Interval, symbol int) bool {
	for _, interval := range intervals {
		if interval.Contains(symbol) {
			return true
		}
	}
	return false
}

func ComplementIntervals(intervals []antlr.Interval, start, stop int) []antlr.Interval {
	result := make([]antlr.Interval, 0)
	result = append(result, antlr.NewInterval(start, stop+1))
	for _, iv := range intervals {
		result = removeIntervalRange(result, iv)
	}
	return result
}

func removeIntervalRange(intervals []antlr.Interval, remove antlr.Interval) []antlr.Interval {
	newIntervals := make([]antlr.Interval, 0)
	for _, iv := range intervals {
		// No overlap
		if remove.Stop <= iv.Start || remove.Start >= iv.Stop {
			newIntervals = append(newIntervals, iv)
			continue
		}
		// Overlap on the left
		if remove.Start > iv.Start {
			newIntervals = append(newIntervals, antlr.NewInterval(iv.Start, remove.Start))
		}
		// Overlap on the right
		if remove.Stop < iv.Stop {
			newIntervals = append(newIntervals, antlr.NewInterval(remove.Stop, iv.Stop))
		}
	}
	return newIntervals
}
