
package expr // Expr
import "github.com/antlr4-go/antlr/v4"

// BaseExprListener is a complete listener for a parse tree produced by ExprParser.
type BaseExprListener struct{}

var _ ExprListener = &BaseExprListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseExprListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseExprListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseExprListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseExprListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterExpression is called when production expression is entered.
func (s *BaseExprListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BaseExprListener) ExitExpression(ctx *ExpressionContext) {}

// EnterAssignment is called when production assignment is entered.
func (s *BaseExprListener) EnterAssignment(ctx *AssignmentContext) {}

// ExitAssignment is called when production assignment is exited.
func (s *BaseExprListener) ExitAssignment(ctx *AssignmentContext) {}

// EnterSimpleExpression is called when production simpleExpression is entered.
func (s *BaseExprListener) EnterSimpleExpression(ctx *SimpleExpressionContext) {}

// ExitSimpleExpression is called when production simpleExpression is exited.
func (s *BaseExprListener) ExitSimpleExpression(ctx *SimpleExpressionContext) {}

// EnterVariableRef is called when production variableRef is entered.
func (s *BaseExprListener) EnterVariableRef(ctx *VariableRefContext) {}

// ExitVariableRef is called when production variableRef is exited.
func (s *BaseExprListener) ExitVariableRef(ctx *VariableRefContext) {}

// EnterFunctionRef is called when production functionRef is entered.
func (s *BaseExprListener) EnterFunctionRef(ctx *FunctionRefContext) {}

// ExitFunctionRef is called when production functionRef is exited.
func (s *BaseExprListener) ExitFunctionRef(ctx *FunctionRefContext) {}
