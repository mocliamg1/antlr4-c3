# ANTLR4-C3 Go Port

This is the Go version of the antlr4-c3 library. Translated from: [antlr4-c3](https://github.com/mike-lischke/antlr4-c3)

## Usage

Write or obtain an antlr4 grammar file, then use antlr4 to generate Go code.

```bash
# Expr.g4 is an antlr4 grammar file, which you can write yourself or obtain from the internet
# 4.13.2 is determined by the antlr4 version in go.mod
# ../example/gen is the directory for the generated code, please modify according to your actual situation
antlr4 -v 4.13.2 -Dlanguage=Go Expr.g4 -o ../example/gen --package expr
```

Then use the generated code in your project:

```go
package main

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	expr "path/to/your/example/gen"
	antlr_c3 "path/to/antlr4-c3/ports/go/lib"
)

func main() {
	expression := "var c = a + b()"
	input := antlr.NewInputStream(expression)
	lexer := expr.NewExprLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	parser := expr.NewExprParser(stream)

	lexer.RemoveErrorListeners()
	parser.RemoveErrorListeners()

	parser.Expression()
	core := antlr_c3.NewCompletionCore(parser, false, false)

	candidates := core.CollectCandidates(stream.Size(), nil)
	fmt.Println(candidates)
}
```