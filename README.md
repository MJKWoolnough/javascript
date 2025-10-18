# javascript

[![CI](https://github.com/MJKWoolnough/javascript/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/javascript/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/javascript.svg)](https://pkg.go.dev/vimagination.zapto.org/javascript)
[![Go Report Card](https://goreportcard.com/badge/vimagination.zapto.org/javascript)](https://goreportcard.com/report/vimagination.zapto.org/javascript)

--
    import "vimagination.zapto.org/javascript"

Package javascript implements a javascript tokeniser and AST.

## Highlights

 - Parse javascript code into AST.
 - Modify parsed code.
 - Consistant javascript formatting.

## Usage

```go
package main

import (
	"fmt"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/parser"
)

func main() {
	src := `function greet(name) {console.log("Hello, " + name)} for (const name of ["Alice", "Bob", "Charlie"]) greet(name)`

	tk := parser.NewStringTokeniser(src)

	ast, err := javascript.ParseScript(&tk)
	if err != nil {
		fmt.Println(err)

		return
	}

	javascript.UnwrapConditional(javascript.WrapConditional(javascript.UnwrapConditional(javascript.UnwrapConditional(ast.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).Arguments.ArgumentList[0].AssignmentExpression.ConditionalExpression).(*javascript.AdditiveExpression).AdditiveExpression)).(*javascript.PrimaryExpression).Literal.Data = `"Hi, "`

	fmt.Printf("%s", ast)

	// Output:
	// function greet(name) {
	//	console.log("Hi, " + name);
	// }
	//
	// for (const name of ["Alice", "Bob", "Charlie"]) greet(name);
}
```

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/javascript
