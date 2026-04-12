package scope_test

import (
	"fmt"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/scope"
	"vimagination.zapto.org/parser"
)

func Example() {
	src := `let a = 1; console.log(a)`

	tk := parser.NewStringTokeniser(src)

	ast, err := javascript.ParseModule(&tk)
	if err != nil {
		fmt.Println(err)

		return
	}

	s, err := scope.Build(ast, nil)
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
