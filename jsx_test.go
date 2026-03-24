package javascript

import (
	"testing"
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
