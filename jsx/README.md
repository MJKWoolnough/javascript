# jsx

[![CI](https://github.com/MJKWoolnough/javascript/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/javascript/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/javascript/jsx.svg)](https://pkg.go.dev/vimagination.zapto.org/javascript/jsx)
[![Go Report Card](https://goreportcard.com/badge/vimagination.zapto.org/javascript)](https://goreportcard.com/report/vimagination.zapto.org/javascript)

--
    import "vimagination.zapto.org/javascript/jsx"

Package jsx provides a simple, template-based transformer of JSX to plain JavaScript.

## Highlights

 - Flexible transpilation of JSX to JavaScript via a provided template.
 - Automatically adds required imports.

## Usage

```go
package main

import (
	"fmt"
	"text/template"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/jsx"
	"vimagination.zapto.org/parser"
)

func main() {
	js := `function MyElement() {
	return <div>Hello, World</div>
}`

	tk := parser.NewStringTokeniser(js)

	m, err := javascript.ParseModule(javascript.AsJSX(&tk))
	if err != nil {
		fmt.Println("unexepected error: ", err)

		return
	}

	tmpl := template.Must(template.New("").Parse("import {createElement} from '@elements';\ncreateElement(\"TAG_NAME\", PARAMS, CHILDREN);"))

	if err = jsx.Process(m, tmpl); err != nil {
		fmt.Println("unexepected error: ", err)

		return
	}

	fmt.Printf("%#s", m)

	// Output:
	// import{createElement}from"@elements"
	// function MyElement() {
	// 	return (createElement("div", {}, ["Hello, World"]))
	// }
}
```

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/javascript/jsx
