package c3

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

var atnStateTypeMap = []string{
	"invalid",
	"basic",
	"rule start",
	"block start",
	"plus block start",
	"star block start",
	"token start",
	"rule stop",
	"block end",
	"star loop back",
	"star loop entry",
	"plus loop back",
	"loop end",
}

func generateBaseDescription(state antlr.ATNState, ruleNames []string) string {
	var stateValue string
	if state.GetStateNumber() == antlr.ATNStateInvalidStateNumber {
		stateValue = "Invalid"
	} else {
		stateValue = fmt.Sprintf("%d", state.GetStateNumber())
	}

	stateTypeIndex := int(state.GetStateType())
	stateType := "unknown"
	if stateTypeIndex >= 0 && stateTypeIndex < len(atnStateTypeMap) {
		stateType = atnStateTypeMap[stateTypeIndex]
	}

	ruleIndex := state.GetRuleIndex()
	ruleName := "unknown"
	if ruleIndex >= 0 && ruleIndex < len(ruleNames) {
		ruleName = ruleNames[ruleIndex]
	}

	return fmt.Sprintf("[%s %s] in %s", stateValue, stateType, ruleName)
}

func (c *CodeCompletionCore) printDescription(currentIndent string, state antlr.ATNState, baseDescription string, tokenIndex int) {
	if !c.printDebugOutput {
		return
	}

	var output string = currentIndent
	var transitionDescription string

	for _, transition := range state.GetTransitions() {
		var labels string
		label := GetLabel(transition)
		var symbols []int
		if label != nil {
			for _, interval := range label.GetIntervals() {
				for i := interval.Start; i <= interval.Stop; i++ {
					symbols = append(symbols, i)
				}
			}
		}
		if len(symbols) > 2 {
			labels += fmt.Sprintf("%s .. %s", c.getDisplayName(symbols[0]), c.getDisplayName(symbols[len(symbols)-1]))
		} else {
			for i, symbol := range symbols {
				if i > 0 {
					labels += ", "
				}
				labels += c.getDisplayName(symbol)
			}
		}
		if labels == "" {
			labels = "ε"
		}

		target := GetTransitionTarget(transition)
		stateTypeIndex := int(target.GetStateType())
		stateType := "unknown"
		if stateTypeIndex >= 0 && stateTypeIndex < len(atnStateTypeMap) {
			stateType = atnStateTypeMap[stateTypeIndex]
		}

		ruleIndex := target.GetRuleIndex()
		ruleName := "unknown"
		if ruleIndex >= 0 && ruleIndex < len(c.rules) {
			ruleName = c.rules[ruleIndex]
		}

		transitionDescription += fmt.Sprintf("\n%s\t(%s) [%d %s] in %s",
			currentIndent, labels, target.GetStateNumber(), stateType, ruleName)
	}

	if tokenIndex >= len(c.tokens)-1 {
		output += fmt.Sprintf("<<%d>> ", c.tokenStartIndex+tokenIndex)
	} else {
		output += fmt.Sprintf("<%d> ", c.tokenStartIndex+tokenIndex)
	}

	fmt.Printf("%s Current state: %s%s\n", output, baseDescription, transitionDescription)
}

func (c *CodeCompletionCore) printRuleState(stack []int) {
	if !c.printDebugOutput {
		return
	}

	if len(stack) == 0 {
		fmt.Println("<empty stack>")
		return
	}

	var sb string
	for _, rule := range stack {
		ruleName := "unknown"
		if rule >= 0 && rule < len(c.rules) {
			ruleName = c.rules[rule]
		}
		sb += fmt.Sprintf("  %s\n", ruleName)
	}
	fmt.Print(sb)
}

func (c *CodeCompletionCore) getDisplayName(token int) string {
	if token < 0 || token >= len(c.symbolicNames) {
		return fmt.Sprintf("Token %d", token)
	}
	if c.symbolicNames[token] != "" {
		return fmt.Sprintf("%s<%d>", c.symbolicNames[token], token)
	}
	if c.literalNames[token] != "" {
		return fmt.Sprintf("%s<%d>", c.literalNames[token], token)
	}
	return fmt.Sprintf("<%d>", token)
}
