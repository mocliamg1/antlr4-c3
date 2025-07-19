package c3

import (
	"fmt"
	"reflect"

	"github.com/antlr4-go/antlr/v4"
)

func (c *CodeCompletionCore) processRule(startRuleState *antlr.RuleStartState, tokenIndex int, callStack []int, indentation string) map[int]struct{} {
	result := make(map[int]struct{})

	positionMap := c.shortcutMap[startRuleState.GetRuleIndex()]
	if positionMap == nil {
		positionMap = make(map[int]map[int]struct{})
		c.shortcutMap[startRuleState.GetRuleIndex()] = positionMap
	} else {
		if positionMap[tokenIndex] != nil {
			if c.printDebugOutput {
				fmt.Printf("=====> shortcut")
			}
			return positionMap[tokenIndex]
		}
	}

	followSets := c.getOrCreateFollowSets(startRuleState)

	callStack = append(callStack, startRuleState.GetRuleIndex())

	currentSymbol := c.tokens[tokenIndex].GetTokenType()

	if tokenIndex >= len(c.tokens)-1 {
		c.handleEndOfInput(startRuleState, followSets, callStack)
		callStack = callStack[:len(callStack)-1]
		return result
	}

	if !followSets.Contains(antlr.TokenEpsilon) && !followSets.Contains(currentSymbol) {
		callStack = callStack[:len(callStack)-1]
		return result
	}

	result = c.processATNStates(startRuleState, tokenIndex, callStack, indentation)

	positionMap[tokenIndex] = result
	callStack = callStack[:len(callStack)-1]

	return result
}

func (c *CodeCompletionCore) getOrCreateFollowSets(startRuleState *antlr.RuleStartState) *FollowSetsHolder {
	parserType := reflect.TypeOf(c.parser).String()
	setsPerState, ok := followSetsByATN[parserType]
	if !ok {
		setsPerState = make(map[int]*FollowSetsHolder)
		followSetsByATN[parserType] = setsPerState
	}

	followSets := setsPerState[startRuleState.GetStateNumber()]
	if followSets == nil {
		followSets = &FollowSetsHolder{
			Sets:     make([]*FollowSetWithPath, 0),
			Combined: make([]antlr.Interval, 0),
		}
		setsPerState[startRuleState.GetStateNumber()] = followSets

		stop := c.ruleToStopState[startRuleState.GetRuleIndex()]
		followSets.Sets = c.determineFollowSets(startRuleState, stop)

		intervals := make([]antlr.Interval, 0)
		for _, set := range followSets.Sets {
			for _, interval := range set.Intervals {
				intervals = append(intervals, interval)
			}
		}
		followSets.Combined = intervals
	}
	return followSets
}

func (c *CodeCompletionCore) handleEndOfInput(startRuleState *antlr.RuleStartState, followSets *FollowSetsHolder, callStack []int) {
	if c.preferredRules[startRuleState.GetRuleIndex()] != nil {
		c.translateToRuleIndex(callStack)
	} else {
		for _, followSet := range followSets.Sets {
			fullPath := append(callStack, followSet.Path...)
			if !c.translateToRuleIndex(fullPath) {
				for _, symbol := range followSet.ToTokenList() {
					if c.ignoredTokens[symbol] == nil {
						if c.printDebugOutput {
							fmt.Printf("Debug: =====> collected: %s for rule %d\n", c.getDisplayName(symbol), startRuleState.GetRuleIndex())
						}
						if c.candidates.Tokens[symbol] == nil {
							c.candidates.Tokens[symbol] = followSet.Following
						} else {
							following := c.candidates.Tokens[symbol]
							if !followSet.FollowingEqual(following) {
								c.candidates.Tokens[symbol] = make([]int, 0)
							}
						}
					} else {
						if c.printDebugOutput {
							fmt.Printf("Debug: =====> ignoring: %s for rule %d\n", c.getDisplayName(symbol), startRuleState.GetRuleIndex())
						}
					}
				}
			}
		}
	}
}

func (c *CodeCompletionCore) processATNStates(startRuleState *antlr.RuleStartState, tokenIndex int, callStack []int, indentation string) map[int]struct{} {
	result := make(map[int]struct{})

	statePipeline := []PipelineEntry{
		{
			State:      startRuleState,
			TokenIndex: tokenIndex,
		},
	}

	for len(statePipeline) > 0 {
		entry := statePipeline[len(statePipeline)-1]
		if c.printDebugOutput {
			c.printDescription(indentation, entry.State, generateBaseDescription(entry.State, c.rules), entry.TokenIndex)
			if c.printRuleStack {
				c.printRuleState(callStack)
			}
		}

		statePipeline = statePipeline[:len(statePipeline)-1]
		c.statesProcessed++

		currentSymbol := c.tokens[entry.TokenIndex].GetTokenType()
		atCaret := entry.TokenIndex >= len(c.tokens)-1

		switch entry.State.GetStateType() {
		case antlr.ATNStateRuleStart:
			indentation += " "
		case antlr.ATNStateRuleStop:
			result[entry.TokenIndex] = struct{}{}
			continue
		}

		statePipeline = c.processTransitions(entry, callStack, indentation, atCaret, currentSymbol, statePipeline)
	}

	return result
}

func (c *CodeCompletionCore) processTransitions(entry PipelineEntry, callStack []int, indentation string, atCaret bool, currentSymbol int, statePipeline []PipelineEntry) []PipelineEntry {
	for _, transition := range entry.State.GetTransitions() {
		switch t := transition.(type) {
		case *antlr.RuleTransition:
			statePipeline = c.processRuleTransition(t, entry, callStack, indentation, statePipeline)
		case *antlr.PredicateTransition:
			if c.checkPredicate(t) {
				statePipeline = append(statePipeline, PipelineEntry{
					State:      GetTransitionTarget(t),
					TokenIndex: entry.TokenIndex,
				})
			}
		case *antlr.WildcardTransition:
			statePipeline = c.processWildcardTransition(entry, callStack, atCaret, statePipeline)
		default:
			statePipeline = c.processGeneralTransition(transition, entry, callStack, atCaret, currentSymbol, statePipeline)
		}
	}
	return statePipeline
}

func (c *CodeCompletionCore) processRuleTransition(ruleTransition *antlr.RuleTransition, entry PipelineEntry, callStack []int, indentation string, statePipeline []PipelineEntry) []PipelineEntry {
	atnState := GetTransitionTarget(ruleTransition)
	endStatus := c.processRule(c.ruleToStartState[atnState.GetRuleIndex()], entry.TokenIndex, callStack, indentation)
	for position := range endStatus {
		statePipeline = append(statePipeline, PipelineEntry{
			State:      GetTransitionFollowState(ruleTransition),
			TokenIndex: position,
		})
	}
	return statePipeline
}

func (c *CodeCompletionCore) processWildcardTransition(entry PipelineEntry, callStack []int, atCaret bool, statePipeline []PipelineEntry) []PipelineEntry {
	if atCaret {
		if c.translateToRuleIndex(callStack) {
			for token := antlr.TokenMinUserTokenType; token <= c.getMaxTokenTypeReflect(); token++ {
				if c.ignoredTokens[token] == nil {
					c.candidates.Tokens[token] = make([]int, 0)
				}
			}
		}
	} else {
		statePipeline = append(statePipeline, PipelineEntry{
			State:      GetTransitionTarget(entry.State.GetTransitions()[0]),
			TokenIndex: entry.TokenIndex + 1,
		})
	}
	return statePipeline
}

func (c *CodeCompletionCore) processGeneralTransition(transition antlr.Transition, entry PipelineEntry, callStack []int, atCaret bool, currentSymbol int, statePipeline []PipelineEntry) []PipelineEntry {
	if GetEpsilon(transition) {
		statePipeline = append(statePipeline, PipelineEntry{
			State:      GetTransitionTarget(transition),
			TokenIndex: entry.TokenIndex,
		})
		return statePipeline
	}

	intervalsLabel := GetLabel(transition).GetIntervals()
	if len(intervalsLabel) > 0 {
		if _, ok := transition.(*antlr.NotSetTransition); ok {
			intervalsLabel = ComplementIntervals(intervalsLabel, antlr.TokenMinUserTokenType, c.getMaxTokenTypeReflect())
		}
	}

	if atCaret {
		if !c.translateToRuleIndex(callStack) {
			c.collectAtCaretTokens(intervalsLabel, transition)
		}
	} else {
		if ContainsSymbol(intervalsLabel, currentSymbol) {
			statePipeline = append(statePipeline, PipelineEntry{
				State:      GetTransitionTarget(transition),
				TokenIndex: entry.TokenIndex + 1,
			})
		}
	}
	return statePipeline
}

func (c *CodeCompletionCore) collectAtCaretTokens(intervalsLabel []antlr.Interval, transition antlr.Transition) {
	list := make([]int, 0)
	for _, interval := range intervalsLabel {
		for i := interval.Start; i < interval.Stop; i++ {
			list = append(list, i)
		}
	}

	addFollowing := len(list) == 1
	for _, symbol := range list {
		if c.ignoredTokens == nil || c.ignoredTokens[symbol] == nil {
			if c.printDebugOutput {
				fmt.Printf("=====> collected: %s\n", c.getDisplayName(symbol))
			}
			if addFollowing {
				c.candidates.Tokens[symbol] = c.getFollowingTokens(transition)
			} else {
				c.candidates.Tokens[symbol] = []int{}
			}
		} else {
			if c.printDebugOutput {
				fmt.Printf("=====> collected: Ignoring token: %d\n", symbol)
			}
		}
	}
}
