package javascript

import (
	"testing"

	"vimagination.zapto.org/parser"
)

func TestJSXElementName(t *testing.T) {
	doTests(t, []sourceFn{
		{"</>", func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrMissingIdentifier,
				Parsing: "JSXElementName",
				Token:   tk[1],
			}
		}},
		{"<a />", func(t *test, tk Tokens) { // 2
			t.Output = JSXElementName{
				Identifier: &tk[1],
				Tokens:     tk[1:2],
			}
		}},
		{"<a: />", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err:     ErrMissingIdentifier,
				Parsing: "JSXElementName",
				Token:   tk[3],
			}
		}},
		{"<a:b />", func(t *test, tk Tokens) { // 4
			t.Output = JSXElementName{
				Namespace:  &tk[1],
				Identifier: &tk[3],
				Tokens:     tk[1:4],
			}
		}},
		{"<a. />", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     ErrMissingIdentifier,
				Parsing: "JSXElementName",
				Token:   tk[3],
			}
		}},
		{"<a.b />", func(t *test, tk Tokens) { // 6
			t.Output = JSXElementName{
				MemberExpression: []*Token{
					&tk[1],
				},
				Identifier: &tk[3],
				Tokens:     tk[1:4],
			}
		}},
		{"<a.b. />", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err:     ErrMissingIdentifier,
				Parsing: "JSXElementName",
				Token:   tk[5],
			}
		}},
		{"<a.b.c />", func(t *test, tk Tokens) { // 8
			t.Output = JSXElementName{
				MemberExpression: []*Token{
					&tk[1],
					&tk[3],
				},
				Identifier: &tk[5],
				Tokens:     tk[1:6],
			}
		}},
	}, func(t *test) (Type, error) {
		var jn JSXElementName

		t.Tokens = t.Tokens[1:1]

		err := jn.parse(&t.Tokens)

		return jn, err
	}, true)
}

func TestJSXElementNameEqual(t *testing.T) {
	for n, test := range [...]struct {
		A, B  JSXElementName
		Match bool
	}{
		{ // 1
			A:     JSXElementName{},
			B:     JSXElementName{},
			Match: true,
		},
		{ // 2
			A:     JSXElementName{Identifier: &Token{Token: parser.Token{Data: "A"}}},
			B:     JSXElementName{},
			Match: false,
		},
		{ // 3
			A:     JSXElementName{Identifier: &Token{Token: parser.Token{Data: "A"}}},
			B:     JSXElementName{Identifier: &Token{Token: parser.Token{Data: "B"}}},
			Match: false,
		},
		{ // 4
			A:     JSXElementName{Identifier: &Token{Token: parser.Token{Data: "A"}}},
			B:     JSXElementName{Identifier: &Token{Token: parser.Token{Data: "A"}}},
			Match: true,
		},
		{ // 5
			A:     JSXElementName{Namespace: &Token{Token: parser.Token{Data: "B"}}, Identifier: &Token{Token: parser.Token{Data: "A"}}},
			B:     JSXElementName{Identifier: &Token{Token: parser.Token{Data: "B"}}},
			Match: false,
		},
		{ // 6
			A:     JSXElementName{Namespace: &Token{Token: parser.Token{Data: "B"}}, Identifier: &Token{Token: parser.Token{Data: "A"}}},
			B:     JSXElementName{Namespace: &Token{Token: parser.Token{Data: "B"}}, Identifier: &Token{Token: parser.Token{Data: "A"}}},
			Match: true,
		},
		{ // 7
			A:     JSXElementName{MemberExpression: []*Token{{Token: parser.Token{Data: "B"}}}, Identifier: &Token{Token: parser.Token{Data: "A"}}},
			B:     JSXElementName{Identifier: &Token{Token: parser.Token{Data: "A"}}},
			Match: false,
		},
		{ // 8
			A:     JSXElementName{MemberExpression: []*Token{{Token: parser.Token{Data: "B"}}}, Identifier: &Token{Token: parser.Token{Data: "A"}}},
			B:     JSXElementName{MemberExpression: []*Token{{Token: parser.Token{Data: "B"}}}, Identifier: &Token{Token: parser.Token{Data: "A"}}},
			Match: true,
		},
		{ // 9
			A:     JSXElementName{MemberExpression: []*Token{{Token: parser.Token{Data: "C"}}, {Token: parser.Token{Data: "B"}}}, Identifier: &Token{Token: parser.Token{Data: "A"}}},
			B:     JSXElementName{MemberExpression: []*Token{{Token: parser.Token{Data: "C"}}, {Token: parser.Token{Data: "D"}}}, Identifier: &Token{Token: parser.Token{Data: "A"}}},
			Match: false,
		},
		{ // 10
			A:     JSXElementName{MemberExpression: []*Token{{Token: parser.Token{Data: "C"}}, {Token: parser.Token{Data: "B"}}}, Identifier: &Token{Token: parser.Token{Data: "A"}}},
			B:     JSXElementName{MemberExpression: []*Token{{Token: parser.Token{Data: "C"}}, {Token: parser.Token{Data: "B"}}}, Identifier: &Token{Token: parser.Token{Data: "A"}}},
			Match: true,
		},
	} {
		if test.A.equal(&test.B) != test.Match {
			t.Errorf("test %d.1: expected match = %v, got %v", n+1, test.Match, !test.Match)
		} else if test.B.equal(&test.A) != test.Match {
			t.Errorf("test %d.2: expected match = %v, got %v", n+1, test.Match, !test.Match)
		}
	}
}
