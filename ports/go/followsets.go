package c3

import (
	"github.com/antlr4-go/antlr/v4"
)

func (c *CodeCompletionCore) determineFollowSets(startRuleState *antlr.RuleStartState, stopState *antlr.RuleStopState) []*FollowSetWithPath {
	result := make([]*FollowSetWithPath, 0)
	seen := make(map[int]struct{})
	ruleStack := make([]int, 0)
	c.collectFollowSets(startRuleState, stopState, &result, seen, ruleStack)
	return result
}

func (c *CodeCompletionCore) translateToRuleIndex(ruleStack []int) bool {
	if len(c.preferredRules) == 0 {
		return false
	}

	for i := 0; i < len(ruleStack); i++ {
		if _, ok := c.preferredRules[ruleStack[i]]; ok {
			path := make([]int, i)
			copy(path, ruleStack[:i])
			addNew := true
			for key, value := range c.candidates.Rules {
				if key == ruleStack[i] && len(value) == len(path) {
					same := true
					for j := range value {
						if value[j] != path[j] {
							same = false
							break
						}
					}
					if same {
						addNew = false
						break
					}
				}
			}
			if addNew {
				c.candidates.Rules[ruleStack[i]] = path
			}
			return true
		}
	}
	return false
}

func (c *CodeCompletionCore) getFollowingTokens(initialTransition antlr.Transition) []int {
	result := []int{}
	seen := map[int]bool{}
	pipeline := []antlr.ATNState{GetTransitionTarget(initialTransition)}

	for len(pipeline) > 0 {
		state := pipeline[len(pipeline)-1]
		pipeline = pipeline[:len(pipeline)-1]

		if seen[state.GetStateNumber()] {
			continue
		}
		seen[state.GetStateNumber()] = true

		for i := 0; i < len(state.GetTransitions()); i++ {
			transition := state.GetTransitions()[i]
			if _, ok := transition.(*antlr.AtomTransition); ok {
				if !GetEpsilon(transition) {
					label := GetLabel(transition)
					if label != nil && label.GetIntervals() != nil && len(label.GetIntervals()) == 1 {
						interval := label.GetIntervals()[0]
						tokenType := interval.Start
						if c.ignoredTokens == nil || c.ignoredTokens[tokenType] == nil {
							result = append(result, tokenType)
							pipeline = append(pipeline, GetTransitionTarget(transition))
						}
					}
				} else {
					pipeline = append(pipeline, GetTransitionTarget(transition))
				}
			}
		}
	}

	return result
}

func (c *CodeCompletionCore) collectFollowSets(
	s antlr.ATNState,
	stopState antlr.ATNState,
	followSets *[]*FollowSetWithPath,
	seen map[int]struct{},
	ruleStack []int,
) {
	if _, ok := seen[s.GetStateNumber()]; ok {
		return
	}
	seen[s.GetStateNumber()] = struct{}{}

	if s == stopState || s.GetStateType() == antlr.ATNStateRuleStop {
		intervals := antlr.NewInterval(antlr.TokenEpsilon, antlr.TokenEpsilon)
		set := &FollowSetWithPath{
			Intervals: []antlr.Interval{intervals},
			Path:      append([]int{}, ruleStack...),
			Following: []int{},
		}
		*followSets = append(*followSets, set)
		return
	}

	for _, transition := range s.GetTransitions() {
		if _, ok := transition.(*antlr.RuleTransition); ok {
			ruleIndex := GetTransitionTarget(transition).GetRuleIndex()
			alreadyInStack := false
			for _, idx := range ruleStack {
				if idx == ruleIndex {
					alreadyInStack = true
					break
				}
			}
			if alreadyInStack {
				continue
			}
			ruleStack = append(ruleStack, ruleIndex)
			c.collectFollowSets(
				GetTransitionTarget(transition),
				stopState,
				followSets,
				seen,
				ruleStack,
			)
			ruleStack = ruleStack[:len(ruleStack)-1]
		} else if predicateTransition, ok := transition.(*antlr.PredicateTransition); ok {
			if c.checkPredicate(predicateTransition) {
				c.collectFollowSets(
					GetTransitionTarget(transition),
					stopState,
					followSets,
					seen,
					ruleStack,
				)
			}
		} else if _, ok := transition.(*antlr.WildcardTransition); ok {
			interval := antlr.NewInterval(antlr.TokenMinUserTokenType, c.getMaxTokenTypeReflect())
			set := &FollowSetWithPath{
				Intervals: []antlr.Interval{interval},
				Path:      append([]int{}, ruleStack...),
				Following: []int{},
			}
			*followSets = append(*followSets, set)
		} else if GetEpsilon(transition) {
			c.collectFollowSets(
				GetTransitionTarget(transition),
				stopState,
				followSets,
				seen,
				ruleStack,
			)
		} else {
			label := GetLabel(transition).GetIntervals()
			if len(label) > 0 {
				if _, ok := transition.(*antlr.NotSetTransition); ok {
					label = ComplementIntervals(label, antlr.TokenMinUserTokenType, c.getMaxTokenTypeReflect())
				}
				set := &FollowSetWithPath{
					Intervals: label,
					Path:      append([]int{}, ruleStack...),
					Following: c.getFollowingTokens(transition),
				}
				*followSets = append(*followSets, set)
			}
		}
	}
}
