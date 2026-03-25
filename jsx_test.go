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

func TestJSXAttribute(t *testing.T) {
	doTests(t, []sourceFn{
		{"</>", func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrMissingIdentifier,
				Parsing: "JSXAttribute",
				Token:   tk[1],
			}
		}},
		{"<a/>", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrMissingEquals,
				Parsing: "JSXAttribute",
				Token:   tk[2],
			}
		}},
		{"<a=/>", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err:     ErrMissingAttribute,
				Parsing: "JSXAttribute",
				Token:   tk[3],
			}
		}},
		{"<a=''/>", func(t *test, tk Tokens) { // 4
			t.Output = JSXAttribute{
				Identifier: &tk[1],
				JSXString:  &tk[3],
				Tokens:     tk[1:4],
			}
		}},
		{"<a:=''/>", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     ErrMissingIdentifier,
				Parsing: "JSXAttribute",
				Token:   tk[3],
			}
		}},
		{"<a:b=''/>", func(t *test, tk Tokens) { // 6
			t.Output = JSXAttribute{
				Namespace:  &tk[1],
				Identifier: &tk[3],
				JSXString:  &tk[5],
				Tokens:     tk[1:6],
			}
		}},
		{"<a={b}/>", func(t *test, tk Tokens) { // 7
			t.Output = JSXAttribute{
				Identifier: &tk[1],
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: WrapConditional(&PrimaryExpression{
						IdentifierReference: &tk[4],
						Tokens:              tk[4:5],
					}),
					Tokens: tk[4:5],
				},
				Tokens: tk[1:6],
			}
		}},
		{"<a={b c}/>", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err:     ErrMissingClosingBrace,
				Parsing: "JSXAttribute",
				Token:   tk[6],
			}
		}},
		{"<a={,}/>", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err:     assignmentCustomError(tk[4], ErrMissingIdentifier),
				Parsing: "JSXAttribute",
				Token:   tk[4],
			}
		}},
		{"<a=<></>/>", func(t *test, tk Tokens) { // 10
			t.Output = JSXAttribute{
				Identifier: &tk[1],
				JSXFragment: &JSXFragment{
					Tokens: tk[3:8],
				},
				Tokens: tk[1:8],
			}
		}},
		{"<a=<></b>/>", func(t *test, tk Tokens) { // 11
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingTagClose,
					Parsing: "JSXFragment",
					Token:   tk[7],
				},
				Parsing: "JSXAttribute",
				Token:   tk[3],
			}
		}},
		{"<a=<b/>/>", func(t *test, tk Tokens) { // 12
			t.Output = JSXAttribute{
				Identifier: &tk[1],
				JSXElement: &JSXElement{
					ElementName: JSXElementName{
						Identifier: &tk[4],
						Tokens:     tk[4:5],
					},
					SelfClosing: true,
					Tokens:      tk[3:7],
				},
				Tokens: tk[1:7],
			}
		}},
		{"<a=<b></c>/>", func(t *test, tk Tokens) { // 13
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidClosingTag,
					Parsing: "JSXElement",
					Token:   tk[9],
				},
				Parsing: "JSXAttribute",
				Token:   tk[3],
			}
		}},
		{"<{}/>", func(t *test, tk Tokens) { // 14
			t.Err = Error{
				Err:     ErrMissingSpread,
				Parsing: "JSXAttribute",
				Token:   tk[2],
			}
		}},
		{"<{...,}/>", func(t *test, tk Tokens) { // 15
			t.Err = Error{
				Err:     assignmentCustomError(tk[3], ErrMissingIdentifier),
				Parsing: "JSXAttribute",
				Token:   tk[3],
			}
		}},
		{"<{...a}/>", func(t *test, tk Tokens) { // 16
			t.Output = JSXAttribute{
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: WrapConditional(&PrimaryExpression{
						IdentifierReference: &tk[3],
						Tokens:              tk[3:4],
					}),
					Tokens: tk[3:4],
				},
				Tokens: tk[1:5],
			}
		}},
		{"<{...a b}/>", func(t *test, tk Tokens) { // 17
			t.Err = Error{
				Err:     ErrMissingClosingBrace,
				Parsing: "JSXAttribute",
				Token:   tk[5],
			}
		}},
	}, func(t *test) (Type, error) {
		var ja JSXAttribute

		t.Tokens = t.Tokens[1:1]

		err := ja.parse(&t.Tokens)

		return ja, err
	}, true)
}
