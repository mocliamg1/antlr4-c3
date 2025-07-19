package main

import (
	"fmt"
	"testing"

	"github.com/antlr4-go/antlr/v4"

	expr "antlr4-c3/ports/go/example/gen"
	antlr_c3 "antlr4-c3/ports/go/lib"
)

// CountingErrorListener counts syntax errors
type CountingErrorListener struct {
	*antlr.DefaultErrorListener
	ErrorCount int
}

func (c *CountingErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	c.ErrorCount++
}

func TestSimpleExpression(t *testing.T) {
	fmt.Println("\nsimpleExpressionTest")

	expression := "var c = a + b()"
	input := antlr.NewInputStream(expression)
	lexer := expr.NewExprLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	parser := expr.NewExprParser(stream)

	// Remove default error listeners
	lexer.RemoveErrorListeners()
	parser.RemoveErrorListeners()

	// Add custom error listener
	errorListener := &CountingErrorListener{}
	parser.AddErrorListener(errorListener)

	parser.Expression()

	if errorListener.ErrorCount != 0 {
		t.Errorf("Expected 0 errors, got %d", errorListener.ErrorCount)
	}

	core := antlr_c3.NewCompletionCore(parser, true, true)

	candidates := core.CollectCandidates(0, nil)

	if len(candidates.Tokens) != 3 {
		t.Errorf("Expected 3 tokens, got %d", len(candidates.Tokens))
	}
	if !containsToken(candidates.Tokens, expr.ExprLexerVAR) {
		t.Error("Expected VAR token to be present")
	}
	if !containsToken(candidates.Tokens, expr.ExprLexerLET) {
		t.Error("Expected LET token to be present")
	}
	if !containsToken(candidates.Tokens, expr.ExprLexerID) {
		t.Error("Expected ID token to be present")
	}

	if len(candidates.Tokens[expr.ExprLexerID]) != 0 {
		t.Errorf("Expected empty token list for ID, got %v", candidates.Tokens[expr.ExprLexerID])
	}

	candidates = core.CollectCandidates(1, nil)
	if len(candidates.Tokens) != 1 {
		t.Errorf("Expected 1 token at position 1, got %d", len(candidates.Tokens))
	}
	if !containsToken(candidates.Tokens, expr.ExprLexerID) {
		t.Error("Expected ID token to be present at position 1")
	}

	candidates = core.CollectCandidates(2, nil)
	if len(candidates.Tokens) != 1 {
		t.Errorf("Expected 1 token at position 2, got %d", len(candidates.Tokens))
	}
	if !containsToken(candidates.Tokens, expr.ExprLexerID) {
		t.Error("Expected ID token to be present at position 2")
	}

	candidates = core.CollectCandidates(4, nil)
	if len(candidates.Tokens) != 1 {
		t.Errorf("Expected 1 token at position 4, got %d", len(candidates.Tokens))
	}
	if !containsToken(candidates.Tokens, expr.ExprLexerEQUAL) {
		t.Error("Expected EQUAL token to be present at position 4")
	}

	candidates = core.CollectCandidates(6, nil)
	if len(candidates.Tokens) != 1 {
		t.Errorf("Expected 1 token at position 6, got %d", len(candidates.Tokens))
	}
	if !containsToken(candidates.Tokens, expr.ExprLexerID) {
		t.Error("Expected ID token to be present at position 6")
	}

	candidates = core.CollectCandidates(8, nil)
	if len(candidates.Tokens) != 5 {
		t.Errorf("Expected 5 tokens at position 8, got %d", len(candidates.Tokens))
	}
	if !containsToken(candidates.Tokens, expr.ExprLexerPLUS) {
		t.Error("Expected PLUS token to be present at position 8")
	}
	if !containsToken(candidates.Tokens, expr.ExprLexerMINUS) {
		t.Error("Expected MINUS token to be present at position 8")
	}
	if !containsToken(candidates.Tokens, expr.ExprLexerMULTIPLY) {
		t.Error("Expected MULTIPLY token to be present at position 8")
	}
	if !containsToken(candidates.Tokens, expr.ExprLexerDIVIDE) {
		t.Error("Expected DIVIDE token to be present at position 8")
	}
	if !containsToken(candidates.Tokens, expr.ExprLexerOPEN_PAR) {
		t.Error("Expected OPEN_PAR token to be present at position 8")
	}
}

func TestTypicalExpression(t *testing.T) {
	fmt.Println("\ntypicalExpressionTest")

	expression := "var c = a + b"
	input := antlr.NewInputStream(expression)
	lexer := expr.NewExprLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	parser := expr.NewExprParser(stream)

	// Set prediction mode
	parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeLLExactAmbigDetection)

	// Remove default error listeners
	lexer.RemoveErrorListeners()
	parser.RemoveErrorListeners()

	// Add custom error listener
	errorListener := &CountingErrorListener{}
	parser.AddErrorListener(errorListener)

	parser.Expression()

	if errorListener.ErrorCount != 0 {
		t.Errorf("Expected 0 errors, got %d", errorListener.ErrorCount)
	}

	preferredRules := map[int]interface{}{
		expr.ExprParserRULE_functionRef: true,
		expr.ExprParserRULE_variableRef: true,
	}

	ignoredTokens := map[int]interface{}{
		expr.ExprLexerID:       true,
		expr.ExprLexerPLUS:     true,
		expr.ExprLexerMINUS:    true,
		expr.ExprLexerMULTIPLY: true,
		expr.ExprLexerDIVIDE:   true,
		expr.ExprLexerEQUAL:    true,
	}

	core := antlr_c3.NewCompletionCore(parser, true, true)
	core.SetPreferredRules(preferredRules)
	core.SetIgnoredTokens(ignoredTokens)

	candidates := core.CollectCandidates(0, nil)

	if len(candidates.Tokens) != 2 {
		t.Errorf("Expected 2 tokens at position 0, got %d", len(candidates.Tokens))
	}
	if !containsToken(candidates.Tokens, expr.ExprLexerVAR) {
		t.Error("Expected VAR token to be present at position 0")
	}
	if !containsToken(candidates.Tokens, expr.ExprLexerLET) {
		t.Error("Expected LET token to be present at position 0")
	}

	candidates = core.CollectCandidates(2, nil)
	if len(candidates.Tokens) != 0 {
		t.Errorf("Expected 0 tokens at position 2, got %d", len(candidates.Tokens))
	}

	candidates = core.CollectCandidates(4, nil)
	if len(candidates.Tokens) != 0 {
		t.Errorf("Expected 0 tokens at position 4, got %d", len(candidates.Tokens))
	}

	candidates = core.CollectCandidates(6, nil)
	if len(candidates.Tokens) != 0 {
		t.Errorf("Expected 0 tokens at position 6, got %d", len(candidates.Tokens))
	}
	if len(candidates.Rules) != 2 {
		t.Errorf("Expected 2 rules at position 6, got %d", len(candidates.Rules))
	}

	found := 0
	for key := range candidates.Rules {
		if key == expr.ExprParserRULE_functionRef || key == expr.ExprParserRULE_variableRef {
			found++
		}
	}
	if found != 2 {
		t.Errorf("Expected 2 matching rules at position 6, got %d", found)
	}

	candidates = core.CollectCandidates(7, nil)
	if len(candidates.Tokens) != 0 {
		t.Errorf("Expected 0 tokens at position 7, got %d", len(candidates.Tokens))
	}
	if len(candidates.Rules) != 1 {
		t.Errorf("Expected 1 rule at position 7, got %d", len(candidates.Rules))
	}

	found = 0
	for key := range candidates.Rules {
		if key == expr.ExprParserRULE_functionRef || key == expr.ExprParserRULE_variableRef {
			found++
		}
	}
	if found != 1 {
		t.Errorf("Expected 1 matching rule at position 7, got %d", found)
	}
}

// Helper function to check if a token exists in candidates
func containsToken(tokens map[int][]int, tokenType int) bool {
	_, exists := tokens[tokenType]
	return exists
}
