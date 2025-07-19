package expr

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type ExprLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var ExprLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func exprlexerLexerInit() {
	staticData := &ExprLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.LiteralNames = []string{
		"", "", "", "'+'", "'-'", "'*'", "'/'", "'='", "'('", "')'",
	}
	staticData.SymbolicNames = []string{
		"", "VAR", "LET", "PLUS", "MINUS", "MULTIPLY", "DIVIDE", "EQUAL", "OPEN_PAR",
		"CLOSE_PAR", "ID", "WS",
	}
	staticData.RuleNames = []string{
		"VAR", "LET", "PLUS", "MINUS", "MULTIPLY", "DIVIDE", "EQUAL", "OPEN_PAR",
		"CLOSE_PAR", "ID", "WS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 11, 56, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2,
		1, 3, 1, 3, 1, 4, 1, 4, 1, 5, 1, 5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 8, 1, 8,
		1, 9, 1, 9, 5, 9, 48, 8, 9, 10, 9, 12, 9, 51, 9, 9, 1, 10, 1, 10, 1, 10,
		1, 10, 0, 0, 11, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 17,
		9, 19, 10, 21, 11, 1, 0, 9, 2, 0, 86, 86, 118, 118, 2, 0, 65, 65, 97, 97,
		2, 0, 82, 82, 114, 114, 2, 0, 76, 76, 108, 108, 2, 0, 69, 69, 101, 101,
		2, 0, 84, 84, 116, 116, 2, 0, 65, 90, 97, 122, 4, 0, 48, 57, 65, 90, 95,
		95, 97, 122, 3, 0, 9, 10, 13, 13, 32, 32, 56, 0, 1, 1, 0, 0, 0, 0, 3, 1,
		0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1,
		0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19,
		1, 0, 0, 0, 0, 21, 1, 0, 0, 0, 1, 23, 1, 0, 0, 0, 3, 27, 1, 0, 0, 0, 5,
		31, 1, 0, 0, 0, 7, 33, 1, 0, 0, 0, 9, 35, 1, 0, 0, 0, 11, 37, 1, 0, 0,
		0, 13, 39, 1, 0, 0, 0, 15, 41, 1, 0, 0, 0, 17, 43, 1, 0, 0, 0, 19, 45,
		1, 0, 0, 0, 21, 52, 1, 0, 0, 0, 23, 24, 7, 0, 0, 0, 24, 25, 7, 1, 0, 0,
		25, 26, 7, 2, 0, 0, 26, 2, 1, 0, 0, 0, 27, 28, 7, 3, 0, 0, 28, 29, 7, 4,
		0, 0, 29, 30, 7, 5, 0, 0, 30, 4, 1, 0, 0, 0, 31, 32, 5, 43, 0, 0, 32, 6,
		1, 0, 0, 0, 33, 34, 5, 45, 0, 0, 34, 8, 1, 0, 0, 0, 35, 36, 5, 42, 0, 0,
		36, 10, 1, 0, 0, 0, 37, 38, 5, 47, 0, 0, 38, 12, 1, 0, 0, 0, 39, 40, 5,
		61, 0, 0, 40, 14, 1, 0, 0, 0, 41, 42, 5, 40, 0, 0, 42, 16, 1, 0, 0, 0,
		43, 44, 5, 41, 0, 0, 44, 18, 1, 0, 0, 0, 45, 49, 7, 6, 0, 0, 46, 48, 7,
		7, 0, 0, 47, 46, 1, 0, 0, 0, 48, 51, 1, 0, 0, 0, 49, 47, 1, 0, 0, 0, 49,
		50, 1, 0, 0, 0, 50, 20, 1, 0, 0, 0, 51, 49, 1, 0, 0, 0, 52, 53, 7, 8, 0,
		0, 53, 54, 1, 0, 0, 0, 54, 55, 6, 10, 0, 0, 55, 22, 1, 0, 0, 0, 2, 0, 49,
		1, 0, 1, 0,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// ExprLexerInit initializes any static state used to implement ExprLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewExprLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func ExprLexerInit() {
	staticData := &ExprLexerLexerStaticData
	staticData.once.Do(exprlexerLexerInit)
}

// NewExprLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewExprLexer(input antlr.CharStream) *ExprLexer {
	ExprLexerInit()
	l := new(ExprLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &ExprLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "Expr.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// ExprLexer tokens.
const (
	ExprLexerVAR       = 1
	ExprLexerLET       = 2
	ExprLexerPLUS      = 3
	ExprLexerMINUS     = 4
	ExprLexerMULTIPLY  = 5
	ExprLexerDIVIDE    = 6
	ExprLexerEQUAL     = 7
	ExprLexerOPEN_PAR  = 8
	ExprLexerCLOSE_PAR = 9
	ExprLexerID        = 10
	ExprLexerWS        = 11
)
