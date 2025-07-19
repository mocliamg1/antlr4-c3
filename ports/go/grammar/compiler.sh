#!/bin/bash
antlr-ng -v 4.13.2 -Dlanguage=Go Expr.g4 -o ../example/gen --package expr