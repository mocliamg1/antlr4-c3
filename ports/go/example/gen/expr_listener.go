package expr // Expr
import "github.com/antlr4-go/antlr/v4"

// ExprListener is a complete listener for a parse tree produced by ExprParser.
type ExprListener interface {
	antlr.ParseTreeListener

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterAssignment is called when entering the assignment production.
	EnterAssignment(c *AssignmentContext)

	// EnterSimpleExpression is called when entering the simpleExpression production.
	EnterSimpleExpression(c *SimpleExpressionContext)

	// EnterVariableRef is called when entering the variableRef production.
	EnterVariableRef(c *VariableRefContext)

	// EnterFunctionRef is called when entering the functionRef production.
	EnterFunctionRef(c *FunctionRefContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitAssignment is called when exiting the assignment production.
	ExitAssignment(c *AssignmentContext)

	// ExitSimpleExpression is called when exiting the simpleExpression production.
	ExitSimpleExpression(c *SimpleExpressionContext)

	// ExitVariableRef is called when exiting the variableRef production.
	ExitVariableRef(c *VariableRefContext)

	// ExitFunctionRef is called when exiting the functionRef production.
	ExitFunctionRef(c *FunctionRefContext)
}
