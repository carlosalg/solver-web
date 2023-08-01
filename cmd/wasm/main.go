package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"syscall/js"
)

func main() {
	fmt.Println("hello from solver")
	js.Global().Set("solverExpr", solverWrapper())
	<-make(chan bool)
}

func solverWrapper() js.Func {
	solverFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid number of arguments passed"
		}
		inputLine := args[0].String()
		fmt.Printf("input %s\n", inputLine)
		expAst, err := parser.ParseExpr(inputLine)
		if err != nil {
			fmt.Println("Error parsing expression:", err)
		}
		result, err := expSolver(expAst)
		if err != nil {
			fmt.Println("Error evaluating expression:", err)
			return err.Error()
		}
		return result
	})
	return solverFunc
}

func expSolver(expr ast.Expr) (float64, error) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		left, err := expSolver(e.X)
		if err != nil {
			return 0, err
		}
		right, err := expSolver(e.Y)
		if err != nil {
			return 0, err
		}
		switch e.Op {
		case token.ADD:
			return left + right, nil
		case token.SUB:
			return left - right, nil
		case token.MUL:
			return left * right, nil
		case token.QUO:
			return left / right, nil
		}
	case *ast.BasicLit:
		value, err := strconv.ParseFloat(e.Value, 64)
		if err != nil {
			return 0, err
		}
		return value, nil
	}
	return 0, fmt.Errorf("unsupported expression: %T", expr)
}
