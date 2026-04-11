# javascript

[![CI](https://github.com/MJKWoolnough/javascript/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/javascript/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/javascript.svg)](https://pkg.go.dev/vimagination.zapto.org/javascript/scope)
[![Go Report Card](https://goreportcard.com/badge/vimagination.zapto.org/javascript)](https://goreportcard.com/report/vimagination.zapto.org/javascript)

--
    import "vimagination.zapto.org/javascript/scope"

Package scope parses out a scope tree for a JavaScript module or script.

## Highlights

 - Process a JavaScript AST into a scope tree, resolving identifiers to their matching declaration.
 - Easily rename identifiers.

## Usage

```go
package main

import (
	"fmt"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/scope"
	"vimagination.zapto.org/parser"
)

func main() {
	src := `let a = 1; console.log(a)`

	tk := parser.NewStringTokeniser(src)

	ast, err := javascript.ParseModule(&tk)
	if err != nil {
		fmt.Println(err)

		return
	}

	s, err := scope.ModuleScope(ast, nil)
	if err != nil {
		fmt.Println(err)

		return
	}

	s.Rename("a", "b")

	fmt.Printf("%s", ast)

	// Output:
	// let b = 1;
	//
	// console.log(b);
}
```

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/javascript/scope
