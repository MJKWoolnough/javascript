package jsx

import (
	"errors"
	"strconv"
	"text/template"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/scope"
	"vimagination.zapto.org/javascript/walk"
	"vimagination.zapto.org/parser"
)

type jsxWalker struct {
	identifiers map[string]map[string][]scope.Binding
	tmpl        *template.Template
	namespace   string
}

func (j *jsxWalker) Handle(t javascript.Type) error {
	ns := j.namespace

	if err := walk.Walk(t, j); err != nil {
		return err
	}

	j.namespace = ns

	switch t := t.(type) {
	case *javascript.JSXAttribute:
		if t.JSXElement != nil {
			pe, err := j.transform(t.JSXElement)
			if err != nil {
				return err
			}

			t.AssignmentExpression = &javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(pe),
			}
			t.JSXElement = nil
		} else if t.JSXFragment != nil {
			t.AssignmentExpression = &javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(childrenToArray(t.JSXFragment.Children)),
			}
			t.JSXFragment = nil
		} else if t.JSXString != nil {
			str, err := javascript.UnescapeJSXString(t.JSXString.Data)
			if err != nil {
				return err
			}

			t.AssignmentExpression = &javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
					Literal: &javascript.Token{
						Token: parser.Token{
							Data: strconv.Quote(str),
							Type: javascript.TokenStringLiteral,
						},
						Pos:     t.JSXString.Pos,
						Line:    t.JSXString.Line,
						LinePos: t.JSXString.LinePos,
					},
				}),
			}
			t.JSXString = nil
		} else if t.AssignmentExpression == nil {
			return javascript.ErrInvalidAssignment
		}
	case *javascript.JSXChild:
		if t.JSXElement != nil {
			pe, err := j.transform(t.JSXElement)
			if err != nil {
				return err
			}

			t.JSXChildExpression = &javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(pe),
			}
			t.JSXElement = nil
		} else if t.JSXFragment != nil {
			t.JSXChildExpression = &javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(childrenToArray(t.JSXFragment.Children)),
			}
			t.JSXFragment = nil
		}
	case *javascript.PrimaryExpression:
		if t.JSXElement != nil {
			pe, err := j.transform(t.JSXElement)
			if err != nil {
				return err
			}

			*t = *pe
		} else if t.JSXFragment != nil {
			t.ArrayLiteral = childrenToArray(t.JSXFragment.Children)
			t.JSXFragment = nil
		}
	}

	return nil
}

func childrenToArray(children []javascript.JSXChild) *javascript.ArrayLiteral {
	return nil
}

func (j *jsxWalker) transform(e *javascript.JSXElement) (*javascript.PrimaryExpression, error) {
	return nil, nil
}

func paramsToObject(attrs []javascript.JSXAttribute) *javascript.ObjectLiteral {
	return nil
}

func Process(m *javascript.Module, tmpl *template.Template) error {
	j := &jsxWalker{
		identifiers: make(map[string]map[string][]scope.Binding),
		tmpl:        tmpl,
	}

	if err := walk.Walk(m, j); err != nil {
		return err
	}

	return nil
}

var (
	ErrInvalidTransformation = errors.New("invalid transformation")
	ErrTooManyStatements     = errors.New("too many statments")
)
