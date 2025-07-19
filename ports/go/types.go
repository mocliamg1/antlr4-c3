package c3

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

type CandidatesCollection struct {
	Tokens        map[int][]int // Token type -> following tokens
	Rules         map[int][]int // Rule index -> call stack path
	RulePositions map[int][]int // Rule index -> [start, end] positions
}

func (c CandidatesCollection) String() string {
	return fmt.Sprintf("CandidatesCollection{tokens=%v, rules=%v, rulePositions=%v}",
		c.Tokens, c.Rules, c.RulePositions)
}

type FollowSetWithPath struct {
	Intervals []antlr.Interval
	Path      []int
	Following []int
}

func (f *FollowSetWithPath) ToTokenList() []int {
	tokens := make([]int, 0)
	for _, interval := range f.Intervals {
		for i := interval.Start; i < interval.Stop; i++ {
			tokens = append(tokens, i)
		}
	}
	return tokens
}

func (f *FollowSetWithPath) FollowingEqual(following []int) bool {
	if len(f.Following) != len(following) {
		return false
	}
	for i := range f.Following {
		if f.Following[i] != following[i] {
			return false
		}
	}
	return true
}

type FollowSetsHolder struct {
	Sets     []*FollowSetWithPath
	Combined []antlr.Interval
}

func (h *FollowSetsHolder) Contains(token int) bool {
	for _, interval := range h.Combined {
		if interval.Start <= token && token <= interval.Stop {
			return true
		}
	}
	return false
}

type PipelineEntry struct {
	State      antlr.ATNState
	TokenIndex int
}
