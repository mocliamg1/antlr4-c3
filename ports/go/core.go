package c3

import (
	"fmt"
	"reflect"

	"github.com/antlr4-go/antlr/v4"
)

var followSetsByATN = make(map[string]map[int]*FollowSetsHolder)

type CodeCompletionCore struct {
	parser           antlr.Parser
	ruleToStartState []*antlr.RuleStartState
	ruleToStopState  []*antlr.RuleStopState
	rules            []string
	tokens           []antlr.Token
	atn              *antlr.ATN
	tokenStartIndex  int
	statesProcessed  int
	candidates       CandidatesCollection
	ignoredTokens    map[int]interface{}
	preferredRules   map[int]interface{}
	printDebugOutput bool
	printRuleStack   bool
	shortcutMap      map[int]map[int]map[int]struct{}
	symbolicNames    []string
	literalNames     []string
}

func NewCompletionCore(p antlr.Parser, printDebugOutput bool, printRuleStack bool) *CodeCompletionCore {
	return &CodeCompletionCore{
		parser:           p,
		ruleToStartState: make([]*antlr.RuleStartState, 0),
		ruleToStopState:  make([]*antlr.RuleStopState, 0),
		rules:            make([]string, 0),
		tokens:           make([]antlr.Token, 0),
		atn:              p.GetInterpreter().ATN(),
		tokenStartIndex:  0,
		statesProcessed:  0,
		candidates: CandidatesCollection{
			Tokens:        make(map[int][]int),
			Rules:         make(map[int][]int),
			RulePositions: make(map[int][]int),
		},
		printDebugOutput: printDebugOutput,
		printRuleStack:   printRuleStack,
		shortcutMap:      make(map[int]map[int]map[int]struct{}),
		symbolicNames:    p.GetSymbolicNames(),
		literalNames:     p.GetLiteralNames(),
	}
}

func (c *CodeCompletionCore) CollectCandidates(caretIndex int, ctx antlr.ParserRuleContext) CandidatesCollection {
	c.shortcutMap = make(map[int]map[int]map[int]struct{})
	c.candidates.Tokens = make(map[int][]int)
	c.candidates.Rules = make(map[int][]int)
	c.candidates.RulePositions = make(map[int][]int)

	c.rules = c.parser.GetRuleNames()
	c.ruleToStartState = c.atn.GetRuleToStartState()
	c.ruleToStopState = c.atn.GetRuleToStopState()

	c.statesProcessed = 0
	var startRule int
	if ctx != nil {
		c.tokenStartIndex = ctx.GetStart().GetTokenIndex()
		startRule = ctx.GetRuleIndex()
	} else {
		c.tokenStartIndex = 0
		startRule = 0
	}

	tokenStream := c.parser.GetTokenStream()
	currentIndex := tokenStream.Index()

	tokenStream.Seek(c.tokenStartIndex)
	c.tokens = make([]antlr.Token, 0)
	offset := 1
	for {
		token := tokenStream.LT(offset)
		offset++
		c.tokens = append(c.tokens, token)
		if token.GetTokenIndex() >= caretIndex || token.GetTokenType() == antlr.TokenEOF {
			break
		}
	}
	tokenStream.Seek(currentIndex)

	callStack := make([]int, 0)
	startRuleState := c.ruleToStartState[startRule]
	c.processRule(startRuleState, 0, callStack, "\n")

	c.calculateRulePositions()

	c.printDebugResults()

	tokenStream.Seek(currentIndex)
	return c.candidates
}

func (c *CodeCompletionCore) SetPreferredRules(preferredRules map[int]interface{}) {
	c.preferredRules = preferredRules
}

func (c *CodeCompletionCore) SetIgnoredTokens(ignoredTokens map[int]interface{}) {
	c.ignoredTokens = ignoredTokens
}

func (c *CodeCompletionCore) EnableDebugOutput() {
	c.printDebugOutput = true
}

func (c *CodeCompletionCore) calculateRulePositions() {
	for ruleID := range c.preferredRules {
		shortcut := c.shortcutMap[ruleID]
		if len(shortcut) == 0 {
			continue
		}

		var startToken int
		for k := range shortcut {
			if k > startToken {
				startToken = k
			}
		}

		endSet := shortcut[startToken]
		var endToken int
		if len(endSet) == 0 {
			endToken = len(c.tokens) - 1
		} else {
			for k := range endSet {
				if k > endToken {
					endToken = k
				}
			}
		}

		startOffset := c.tokens[startToken].GetStart()
		var endOffset int
		if c.tokens[endToken].GetTokenType() == antlr.TokenEOF {
			endOffset = c.tokens[endToken].GetStart()
		} else {
			endOffset = c.tokens[endToken-1].GetStop() + 1
		}

		c.candidates.RulePositions[ruleID] = []int{startOffset, endOffset}
	}
}

func (c *CodeCompletionCore) printDebugResults() {
	if !c.printDebugOutput {
		return
	}

	var logMessage string
	logMessage += fmt.Sprintf("States processed: %d\n", c.statesProcessed)

	logMessage += "Collected rules:\n"
	for key, value := range c.candidates.Rules {
		logMessage += fmt.Sprintf("  %d, path: ", key)
		for _, token := range value {
			ruleName := "unknown"
			if token >= 0 && token < len(c.rules) {
				ruleName = c.rules[token]
			}
			logMessage += fmt.Sprintf("%s ", ruleName)
		}
		logMessage += "\n"
	}

	logMessage += "Collected Tokens:\n"
	for key, value := range c.candidates.Tokens {
		logMessage += fmt.Sprintf("  %s", c.getDisplayName(key))
		for _, following := range value {
			logMessage += fmt.Sprintf(" %s", c.getDisplayName(following))
		}
		logMessage += "\n"
	}

	fmt.Print(logMessage)
}

func (c *CodeCompletionCore) getMaxTokenTypeReflect() int {
	val := reflect.ValueOf(c.atn).Elem().FieldByName("maxTokenType")
	return int(val.Int())
}

func (c *CodeCompletionCore) checkPredicate(transition *antlr.PredicateTransition) bool {
	return transition.GetPredicate().Evaluate(c.parser, antlr.ParserRuleContextEmpty)
}
