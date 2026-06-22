package jsx_test

import (
	"fmt"
	"text/template"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/jsx"
	"vimagination.zapto.org/parser"
)

func Example() {
	js := `function MyElement() {
	return <div id="example">Hello, World</div>
}`

	tk := parser.NewStringTokeniser(js)

	m, err := javascript.ParseModule(javascript.AsJSX(&tk))
	if err != nil {
		fmt.Println("unexpected error: ", err)

		return
	}

	tmpl := template.Must(template.New("").Parse("import {createElement} from '@elements';\ncreateElement(\"TAG_NAME\", PARAMS, CHILDREN);"))

	if err = jsx.Process(m, tmpl); err != nil {
		fmt.Println("unexpected error: ", err)

		return
	}

	fmt.Printf("%#s", m)

	// Output:
	// import{createElement}from"@elements"
	// function MyElement() {
	// 	return (createElement("div", {"id":"example"}, ["Hello, World"]))
	// }
}
