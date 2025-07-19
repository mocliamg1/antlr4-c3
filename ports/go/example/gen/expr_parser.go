package expr // Expr
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type ExprParser struct {
	*antlr.BaseParser
}

var ExprParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func exprParserInit() {
	staticData := &ExprParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "'+'", "'-'", "'*'", "'/'", "'='", "'('", "')'",
	}
	staticData.SymbolicNames = []string{
		"", "VAR", "LET", "PLUS", "MINUS", "MULTIPLY", "DIVIDE", "EQUAL", "OPEN_PAR",
		"CLOSE_PAR", "ID", "WS",
	}
	staticData.RuleNames = []string{
		"expression", "assignment", "simpleExpression", "variableRef", "functionRef",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 11, 42, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 1, 0, 1, 0, 3, 0, 13, 8, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2,
		1, 2, 3, 2, 23, 8, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 5, 2, 31, 8,
		2, 10, 2, 12, 2, 34, 9, 2, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 0,
		1, 4, 5, 0, 2, 4, 6, 8, 0, 3, 1, 0, 1, 2, 1, 0, 3, 4, 1, 0, 5, 6, 40, 0,
		12, 1, 0, 0, 0, 2, 14, 1, 0, 0, 0, 4, 22, 1, 0, 0, 0, 6, 35, 1, 0, 0, 0,
		8, 37, 1, 0, 0, 0, 10, 13, 3, 2, 1, 0, 11, 13, 3, 4, 2, 0, 12, 10, 1, 0,
		0, 0, 12, 11, 1, 0, 0, 0, 13, 1, 1, 0, 0, 0, 14, 15, 7, 0, 0, 0, 15, 16,
		5, 10, 0, 0, 16, 17, 5, 7, 0, 0, 17, 18, 3, 4, 2, 0, 18, 3, 1, 0, 0, 0,
		19, 20, 6, 2, -1, 0, 20, 23, 3, 6, 3, 0, 21, 23, 3, 8, 4, 0, 22, 19, 1,
		0, 0, 0, 22, 21, 1, 0, 0, 0, 23, 32, 1, 0, 0, 0, 24, 25, 10, 4, 0, 0, 25,
		26, 7, 1, 0, 0, 26, 31, 3, 4, 2, 5, 27, 28, 10, 3, 0, 0, 28, 29, 7, 2,
		0, 0, 29, 31, 3, 4, 2, 4, 30, 24, 1, 0, 0, 0, 30, 27, 1, 0, 0, 0, 31, 34,
		1, 0, 0, 0, 32, 30, 1, 0, 0, 0, 32, 33, 1, 0, 0, 0, 33, 5, 1, 0, 0, 0,
		34, 32, 1, 0, 0, 0, 35, 36, 5, 10, 0, 0, 36, 7, 1, 0, 0, 0, 37, 38, 5,
		10, 0, 0, 38, 39, 5, 8, 0, 0, 39, 40, 5, 9, 0, 0, 40, 9, 1, 0, 0, 0, 4,
		12, 22, 30, 32,
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

// ExprParserInit initializes any static state used to implement ExprParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewExprParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func ExprParserInit() {
	staticData := &ExprParserStaticData
	staticData.once.Do(exprParserInit)
}

// NewExprParser produces a new parser instance for the optional input antlr.TokenStream.
func NewExprParser(input antlr.TokenStream) *ExprParser {
	ExprParserInit()
	this := new(ExprParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &ExprParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "Expr.g4"

	return this
}

// ExprParser tokens.
const (
	ExprParserEOF       = antlr.TokenEOF
	ExprParserVAR       = 1
	ExprParserLET       = 2
	ExprParserPLUS      = 3
	ExprParserMINUS     = 4
	ExprParserMULTIPLY  = 5
	ExprParserDIVIDE    = 6
	ExprParserEQUAL     = 7
	ExprParserOPEN_PAR  = 8
	ExprParserCLOSE_PAR = 9
	ExprParserID        = 10
	ExprParserWS        = 11
)

// ExprParser rules.
const (
	ExprParserRULE_expression       = 0
	ExprParserRULE_assignment       = 1
	ExprParserRULE_simpleExpression = 2
	ExprParserRULE_variableRef      = 3
	ExprParserRULE_functionRef      = 4
)

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Assignment() IAssignmentContext
	SimpleExpression() ISimpleExpressionContext

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ExprParserRULE_expression
	return p
}

func InitEmptyExpressionContext(p *ExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ExprParserRULE_expression
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) Assignment() IAssignmentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAssignmentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAssignmentContext)
}

func (s *ExpressionContext) SimpleExpression() ISimpleExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISimpleExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISimpleExpressionContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (p *ExprParser) Expression() (localctx IExpressionContext) {
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, ExprParserRULE_expression)
	p.SetState(12)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case ExprParserVAR, ExprParserLET:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(10)
			p.Assignment()
		}

	case ExprParserID:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(11)
			p.simpleExpression(0)
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAssignmentContext is an interface to support dynamic dispatch.
type IAssignmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	EQUAL() antlr.TerminalNode
	SimpleExpression() ISimpleExpressionContext
	VAR() antlr.TerminalNode
	LET() antlr.TerminalNode

	// IsAssignmentContext differentiates from other interfaces.
	IsAssignmentContext()
}

type AssignmentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAssignmentContext() *AssignmentContext {
	var p = new(AssignmentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ExprParserRULE_assignment
	return p
}

func InitEmptyAssignmentContext(p *AssignmentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ExprParserRULE_assignment
}

func (*AssignmentContext) IsAssignmentContext() {}

func NewAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignmentContext {
	var p = new(AssignmentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_assignment

	return p
}

func (s *AssignmentContext) GetParser() antlr.Parser { return s.parser }

func (s *AssignmentContext) ID() antlr.TerminalNode {
	return s.GetToken(ExprParserID, 0)
}

func (s *AssignmentContext) EQUAL() antlr.TerminalNode {
	return s.GetToken(ExprParserEQUAL, 0)
}

func (s *AssignmentContext) SimpleExpression() ISimpleExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISimpleExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISimpleExpressionContext)
}

func (s *AssignmentContext) VAR() antlr.TerminalNode {
	return s.GetToken(ExprParserVAR, 0)
}

func (s *AssignmentContext) LET() antlr.TerminalNode {
	return s.GetToken(ExprParserLET, 0)
}

func (s *AssignmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AssignmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterAssignment(s)
	}
}

func (s *AssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitAssignment(s)
	}
}

func (p *ExprParser) Assignment() (localctx IAssignmentContext) {
	localctx = NewAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, ExprParserRULE_assignment)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(14)
		_la = p.GetTokenStream().LA(1)

		if !(_la == ExprParserVAR || _la == ExprParserLET) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(15)
		p.Match(ExprParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(16)
		p.Match(ExprParserEQUAL)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(17)
		p.simpleExpression(0)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISimpleExpressionContext is an interface to support dynamic dispatch.
type ISimpleExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VariableRef() IVariableRefContext
	FunctionRef() IFunctionRefContext
	AllSimpleExpression() []ISimpleExpressionContext
	SimpleExpression(i int) ISimpleExpressionContext
	PLUS() antlr.TerminalNode
	MINUS() antlr.TerminalNode
	MULTIPLY() antlr.TerminalNode
	DIVIDE() antlr.TerminalNode

	// IsSimpleExpressionContext differentiates from other interfaces.
	IsSimpleExpressionContext()
}

type SimpleExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySimpleExpressionContext() *SimpleExpressionContext {
	var p = new(SimpleExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ExprParserRULE_simpleExpression
	return p
}

func InitEmptySimpleExpressionContext(p *SimpleExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ExprParserRULE_simpleExpression
}

func (*SimpleExpressionContext) IsSimpleExpressionContext() {}

func NewSimpleExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SimpleExpressionContext {
	var p = new(SimpleExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_simpleExpression

	return p
}

func (s *SimpleExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *SimpleExpressionContext) VariableRef() IVariableRefContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableRefContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableRefContext)
}

func (s *SimpleExpressionContext) FunctionRef() IFunctionRefContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionRefContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionRefContext)
}

func (s *SimpleExpressionContext) AllSimpleExpression() []ISimpleExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISimpleExpressionContext); ok {
			len++
		}
	}

	tst := make([]ISimpleExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISimpleExpressionContext); ok {
			tst[i] = t.(ISimpleExpressionContext)
			i++
		}
	}

	return tst
}

func (s *SimpleExpressionContext) SimpleExpression(i int) ISimpleExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISimpleExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISimpleExpressionContext)
}

func (s *SimpleExpressionContext) PLUS() antlr.TerminalNode {
	return s.GetToken(ExprParserPLUS, 0)
}

func (s *SimpleExpressionContext) MINUS() antlr.TerminalNode {
	return s.GetToken(ExprParserMINUS, 0)
}

func (s *SimpleExpressionContext) MULTIPLY() antlr.TerminalNode {
	return s.GetToken(ExprParserMULTIPLY, 0)
}

func (s *SimpleExpressionContext) DIVIDE() antlr.TerminalNode {
	return s.GetToken(ExprParserDIVIDE, 0)
}

func (s *SimpleExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimpleExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SimpleExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterSimpleExpression(s)
	}
}

func (s *SimpleExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitSimpleExpression(s)
	}
}

func (p *ExprParser) SimpleExpression() (localctx ISimpleExpressionContext) {
	return p.simpleExpression(0)
}

func (p *ExprParser) simpleExpression(_p int) (localctx ISimpleExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewSimpleExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx ISimpleExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 4
	p.EnterRecursionRule(localctx, 4, ExprParserRULE_simpleExpression, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(22)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(20)
			p.VariableRef()
		}

	case 2:
		{
			p.SetState(21)
			p.FunctionRef()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(32)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(30)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext()) {
			case 1:
				localctx = NewSimpleExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_simpleExpression)
				p.SetState(24)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
					goto errorExit
				}
				{
					p.SetState(25)
					_la = p.GetTokenStream().LA(1)

					if !(_la == ExprParserPLUS || _la == ExprParserMINUS) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(26)
					p.simpleExpression(5)
				}

			case 2:
				localctx = NewSimpleExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, ExprParserRULE_simpleExpression)
				p.SetState(27)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(28)
					_la = p.GetTokenStream().LA(1)

					if !(_la == ExprParserMULTIPLY || _la == ExprParserDIVIDE) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(29)
					p.simpleExpression(4)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(34)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IVariableRefContext is an interface to support dynamic dispatch.
type IVariableRefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode

	// IsVariableRefContext differentiates from other interfaces.
	IsVariableRefContext()
}

type VariableRefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVariableRefContext() *VariableRefContext {
	var p = new(VariableRefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ExprParserRULE_variableRef
	return p
}

func InitEmptyVariableRefContext(p *VariableRefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ExprParserRULE_variableRef
}

func (*VariableRefContext) IsVariableRefContext() {}

func NewVariableRefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariableRefContext {
	var p = new(VariableRefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_variableRef

	return p
}

func (s *VariableRefContext) GetParser() antlr.Parser { return s.parser }

func (s *VariableRefContext) ID() antlr.TerminalNode {
	return s.GetToken(ExprParserID, 0)
}

func (s *VariableRefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VariableRefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VariableRefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterVariableRef(s)
	}
}

func (s *VariableRefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitVariableRef(s)
	}
}

func (p *ExprParser) VariableRef() (localctx IVariableRefContext) {
	localctx = NewVariableRefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, ExprParserRULE_variableRef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(35)
		p.Match(ExprParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFunctionRefContext is an interface to support dynamic dispatch.
type IFunctionRefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	OPEN_PAR() antlr.TerminalNode
	CLOSE_PAR() antlr.TerminalNode

	// IsFunctionRefContext differentiates from other interfaces.
	IsFunctionRefContext()
}

type FunctionRefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionRefContext() *FunctionRefContext {
	var p = new(FunctionRefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ExprParserRULE_functionRef
	return p
}

func InitEmptyFunctionRefContext(p *FunctionRefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ExprParserRULE_functionRef
}

func (*FunctionRefContext) IsFunctionRefContext() {}

func NewFunctionRefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionRefContext {
	var p = new(FunctionRefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ExprParserRULE_functionRef

	return p
}

func (s *FunctionRefContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionRefContext) ID() antlr.TerminalNode {
	return s.GetToken(ExprParserID, 0)
}

func (s *FunctionRefContext) OPEN_PAR() antlr.TerminalNode {
	return s.GetToken(ExprParserOPEN_PAR, 0)
}

func (s *FunctionRefContext) CLOSE_PAR() antlr.TerminalNode {
	return s.GetToken(ExprParserCLOSE_PAR, 0)
}

func (s *FunctionRefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionRefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionRefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.EnterFunctionRef(s)
	}
}

func (s *FunctionRefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ExprListener); ok {
		listenerT.ExitFunctionRef(s)
	}
}

func (p *ExprParser) FunctionRef() (localctx IFunctionRefContext) {
	localctx = NewFunctionRefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, ExprParserRULE_functionRef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(37)
		p.Match(ExprParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(38)
		p.Match(ExprParserOPEN_PAR)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(39)
		p.Match(ExprParserCLOSE_PAR)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

func (p *ExprParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 2:
		var t *SimpleExpressionContext = nil
		if localctx != nil {
			t = localctx.(*SimpleExpressionContext)
		}
		return p.SimpleExpression_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *ExprParser) SimpleExpression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 4)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 3)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
