# walk

[![CI](https://github.com/MJKWoolnough/javascript/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/javascript/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/javascript.svg)](https://pkg.go.dev/vimagination.zapto.org/javascript/walk)
[![Go Report Card](https://goreportcard.com/badge/vimagination.zapto.org/javascript/walk)](https://goreportcard.com/report/vimagination.zapto.org/javascript/walk)

--
    import "vimagination.zapto.org/javascript/walk"

Package walk provides a JavaScript type walker.

## Highlights

 - Simple interface to allow control over walking through parsed JavaScript.
 - Allows modification to the tree as it's being walked.

## Usage

```go
package main

import (
	"fmt"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/walk"
	"vimagination.zapto.org/parser"
)

func main() {
	src := "const a = 'b' - `c`"
	tk := parser.NewStringTokeniser(src)

	m, _ := javascript.ParseModule(&tk)

	var walkFn walk.Handler

	walkFn = walk.HandlerFunc(func(t javascript.Type) error {
		switch t := t.(type) {
		case *javascript.AdditiveExpression:
			if t.AdditiveOperator == javascript.AdditiveMinus {
				t.AdditiveOperator = javascript.AdditiveAdd
			}
		case *javascript.PrimaryExpression:
			if t.Literal != nil {
				t.Literal.Data = "'Hello'"
			}
		case *javascript.TemplateLiteral:
			t.NoSubstitutionTemplate.Data = "`, world`"
		}

		return walk.Walk(t, walkFn)
	})

	walk.Walk(m, walkFn)

	fmt.Printf("%s", m)

	// Output:
	// const a = 'Hello' + `, world`;
}
```

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/javascript/walk
