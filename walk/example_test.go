package walk_test

import (
	"fmt"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/walk"
	"vimagination.zapto.org/parser"
)

func Example() {
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
