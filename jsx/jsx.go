package jsx

import (
	"errors"
	"slices"
	"strconv"
	"strings"
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
	al := &javascript.ArrayLiteral{
		ElementList: make([]javascript.ArrayElement, 0, len(children)),
	}

	for _, child := range children {
		al.ElementList = append(al.ElementList, javascript.ArrayElement{
			Spread:               child.Spread,
			AssignmentExpression: *child.JSXChildExpression,
		})
	}

	return al
}

var (
	htmlElements = [...]string{"a", "abbr", "address", "area", "article", "aside", "audio", "b", "base", "bdi", "bdo", "blockquote", "body", "br", "button", "canvas", "caption", "cite", "code", "col", "colgroup", "data", "datalist", "dd", "del", "details", "dfn", "dialog", "div", "dl", "dt", "em", "embed", "fieldset", "figcaption", "figure", "footer", "form", "h1", "h2", "h3", "h4", "h5", "h6", "head", "header", "hgroup", "hr", "html", "i", "iframe", "img", "input", "ins", "kbd", "label", "legend", "li", "link", "main", "map", "mark", "menu", "meta", "meter", "nav", "noscript", "object", "ol", "optgroup", "option", "output", "p", "picture", "pre", "progress", "q", "rp", "rt", "ruby", "s", "samp", "script", "search", "section", "select", "slot", "small", "source", "span", "strong", "style", "sub", "summary", "sup", "table", "tbody", "td", "template", "textarea", "tfoot", "th", "thead", "time", "title", "tr", "track", "u", "ul", "var", "video", "wbr"}
	svgElements  = [...]string{"a", "animate", "animateMotion", "animateTransform", "circle", "clipPath", "defs", "desc", "ellipse", "feBlend", "feColorMatrix", "feComponentTransfer", "feComposite", "feConvolveMatrix", "feDiffuseLighting", "feDisplacementMap", "feDistantLight", "feDropShadow", "feFlood", "feFuncA", "feFuncB", "feFuncG", "feFuncR", "feGaussianBlur", "feImage", "feMerge", "feMergeNode", "feMorphology", "feOffset", "fePointLight", "feSpecularLighting", "feSpotLight", "feTile", "feTurbulence", "filter", "foreignObject", "g", "image", "line", "linearGradient", "marker", "mask", "metadata", "mpath", "path", "pattern", "polygon", "polyline", "radialGradient", "rect", "script", "set", "stop", "style", "svg", "switch", "symbol", "text", "textPath", "title", "tspan", "use", "view"}
)

type templateVars struct {
	Namespace     string
	InHTML, InSVG bool
}

func (j *jsxWalker) transform(e *javascript.JSXElement) (*javascript.PrimaryExpression, error) {
	name := e.ElementName.Identifier
	if name == nil {
		return nil, javascript.ErrInvalidAssignment
	}

	inHTML, inSVG := slices.Contains(htmlElements[:], name.Data), slices.Contains(svgElements[:], name.Data)
	if inHTML && !inSVG {
		j.namespace = "html"
	} else if inSVG && !inHTML {
		j.namespace = "svg"
	}

	var sb strings.Builder

	if err := j.tmpl.Execute(&sb, templateVars{
		Namespace: j.namespace,
		InHTML:    inHTML,
		InSVG:     inSVG,
	}); err != nil {
		return nil, err
	}

	tk := parser.NewStringTokeniser(sb.String())

	m, err := javascript.ParseModule(&tk)
	if err != nil {
		return nil, err
	}

	return j.process(e, m)
}

func (j *jsxWalker) process(e *javascript.JSXElement, m *javascript.Module) (*javascript.PrimaryExpression, error) {
	replaceTagName(m, e.ElementName.Identifier.Data)

	s, err := scope.ModuleScope(m, nil)
	if err != nil {
		return nil, err
	}

	delete(s.Bindings, "PARAMS")
	delete(s.Bindings, "CHILDREN")

	replaceParamsAndChildren(m, e)
	j.gatherIdentifiers(m, s)

	var expression *javascript.Expression

	for _, mi := range m.ModuleListItems {
		if mi.ImportDeclaration != nil {
			continue
		} else if expression != nil {
			return nil, ErrTooManyStatements
		} else if mi.StatementListItem != nil && mi.StatementListItem.Statement != nil && mi.StatementListItem.Statement.ExpressionStatement != nil {
			expression = mi.StatementListItem.Statement.ExpressionStatement
		} else {
			return nil, ErrInvalidTransformation
		}
	}

	return &javascript.PrimaryExpression{
		ParenthesizedExpression: &javascript.ParenthesizedExpression{
			Expressions: expression.Expressions,
		},
	}, nil
}

func replaceTagName(m *javascript.Module, name string) {
}

func replaceParamsAndChildren(m *javascript.Module, e *javascript.JSXElement) {
}

type importIdent struct {
	tk       *javascript.Token
	contains bool
}

func (i *importIdent) Handle(t javascript.Type) error {
	switch t := t.(type) {
	case *javascript.ImportClause:
		if t.ImportedDefaultBinding == i.tk || t.NameSpaceImport == i.tk {
			i.contains = true
		}
	case *javascript.ImportSpecifier:
		if t.IdentifierName == i.tk {
			i.contains = true
		}
	default:
		walk.Walk(t, i)
	}

	return nil
}

func (j *jsxWalker) gatherIdentifiers(m *javascript.Module, s *scope.Scope) {
	imports := make(map[*javascript.Token]string)

	for _, b := range s.Bindings {
		if b[0].BindingType == scope.BindingImport {
			for _, mi := range m.ModuleListItems {
				if mi.ImportDeclaration != nil {
					ii := importIdent{tk: b[0].Token}

					walk.Walk(mi, &ii)

					if ii.contains {
						imports[b[0].Token] = mi.ImportDeclaration.FromClause.ModuleSpecifier.Data

						break
					}
				}
			}
		}
	}

	for _, b := range s.Bindings {
		is := imports[b[0].Token]
		ident := b[0].Token.Data

		imp, ok := j.identifiers[ident]
		if !ok {
			imp = make(map[string][]scope.Binding)
			j.identifiers[ident] = imp
		}

		imp[is] = append(imp[is], b[1:]...)
	}
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
