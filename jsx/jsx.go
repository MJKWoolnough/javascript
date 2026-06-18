// Package jsx allows for the transforming of JSX within a JavaScript AST.
package jsx

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/scope"
	"vimagination.zapto.org/javascript/walk"
	"vimagination.zapto.org/parser"
)

const (
	tagName  = "TAG_NAME"
	children = "CHILDREN"
	params   = "PARAMS"
)

type jsxTransformer struct {
	tmpl      *template.Template
	namespace string
	imports   map[string]struct{}
}

func (j *jsxTransformer) Handle(t javascript.Type) error {
	switch t := t.(type) {
	case *javascript.JSXElement:
		if t.ElementName.Identifier == nil {
			return javascript.ErrMissingIdentifier
		}

		defer j.setNamespace(t.ElementName)()
	}

	if err := walk.Walk(t, j); err != nil {
		return err
	}

	switch t := t.(type) {
	case *javascript.MemberExpression:
		if t.PrimaryExpression != nil {
			pe, err := j.handlePrimaryExpression(t.PrimaryExpression)
			if err != nil {
				return err
			}

			t.PrimaryExpression = pe
		}
	}

	return nil
}

func (j *jsxTransformer) setNamespace(name javascript.JSXElementName) func() {
	ns := j.namespace

	if name.Namespace != nil {
		j.namespace = name.Namespace.Data
	} else if inHTML, inSVG, inMathML := nsIn(name.Identifier); inHTML && !inSVG {
		j.namespace = "html"
	} else if inSVG && !inHTML {
		j.namespace = "svg"
	} else if inMathML {
		j.namespace = "mathml"
	}

	return func() { j.namespace = ns }
}

func (j *jsxTransformer) handlePrimaryExpression(t *javascript.PrimaryExpression) (*javascript.PrimaryExpression, error) {
	if t.JSXElement != nil {
		return j.transform(t.JSXElement)
	} else if t.JSXFragment != nil {
		al, err := j.childrenToArray(t.JSXFragment.Children)
		if err != nil {
			return nil, err
		}

		al.Tokens = t.JSXFragment.Tokens

		return &javascript.PrimaryExpression{
			ArrayLiteral: al,
			Tokens:       t.Tokens,
		}, nil
	}

	return t, nil
}

func (j *jsxTransformer) childrenToArray(children []javascript.JSXChild) (*javascript.ArrayLiteral, error) {
	al := &javascript.ArrayLiteral{
		ElementList: make([]javascript.ArrayElement, 0, len(children)),
	}

	for _, child := range children {
		ae, err := j.childToAE(child)
		if err != nil {
			return nil, err
		}

		ae.Tokens = child.Tokens

		al.ElementList = append(al.ElementList, javascript.ArrayElement{
			Spread:               child.Spread,
			AssignmentExpression: ae,
			Comments:             child.Comments,
			Tokens:               child.Tokens,
		})
	}

	return al, nil
}

func (j *jsxTransformer) childToAE(t javascript.JSXChild) (javascript.AssignmentExpression, error) {
	if t.JSXElement != nil {
		pe, err := j.transform(t.JSXElement)
		if err != nil {
			return javascript.AssignmentExpression{}, err
		}

		return javascript.AssignmentExpression{
			ConditionalExpression: javascript.WrapConditional(pe),
		}, nil
	} else if t.JSXFragment != nil {
		al, err := j.childrenToArray(t.JSXFragment.Children)
		if err != nil {
			return javascript.AssignmentExpression{}, err
		}

		return javascript.AssignmentExpression{
			ConditionalExpression: javascript.WrapConditional(al),
		}, nil
	} else if t.JSXText != nil {
		return javascript.AssignmentExpression{
			ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
				Literal: &javascript.Token{
					Token: parser.Token{
						Data: strconv.Quote(t.JSXText.Data),
					},
				},
			}),
		}, nil
	} else if t.JSXChildExpression == nil {
		return javascript.AssignmentExpression{}, ErrMissingChild
	}

	return *t.JSXChildExpression, nil
}

var (
	htmlElements  = [...]string{"a", "abbr", "address", "area", "article", "aside", "audio", "b", "base", "bdi", "bdo", "blockquote", "body", "br", "button", "canvas", "caption", "cite", "code", "col", "colgroup", "data", "datalist", "dd", "del", "details", "dfn", "dialog", "div", "dl", "dt", "em", "embed", "fieldset", "figcaption", "figure", "footer", "form", "h1", "h2", "h3", "h4", "h5", "h6", "head", "header", "hgroup", "hr", "html", "i", "iframe", "img", "input", "ins", "kbd", "label", "legend", "li", "link", "main", "map", "mark", "menu", "meta", "meter", "nav", "noscript", "object", "ol", "optgroup", "option", "output", "p", "picture", "pre", "progress", "q", "rp", "rt", "ruby", "s", "samp", "script", "search", "section", "select", "slot", "small", "source", "span", "strong", "style", "sub", "summary", "sup", "table", "tbody", "td", "template", "textarea", "tfoot", "th", "thead", "time", "title", "tr", "track", "u", "ul", "var", "video", "wbr"}
	svgElements   = [...]string{"a", "animate", "animateMotion", "animateTransform", "circle", "clipPath", "defs", "desc", "ellipse", "feBlend", "feColorMatrix", "feComponentTransfer", "feComposite", "feConvolveMatrix", "feDiffuseLighting", "feDisplacementMap", "feDistantLight", "feDropShadow", "feFlood", "feFuncA", "feFuncB", "feFuncG", "feFuncR", "feGaussianBlur", "feImage", "feMerge", "feMergeNode", "feMorphology", "feOffset", "fePointLight", "feSpecularLighting", "feSpotLight", "feTile", "feTurbulence", "filter", "foreignObject", "g", "image", "line", "linearGradient", "marker", "mask", "metadata", "mpath", "path", "pattern", "polygon", "polyline", "radialGradient", "rect", "script", "set", "stop", "style", "svg", "switch", "symbol", "text", "textPath", "title", "tspan", "use", "view"}
	mathMLElement = [...]string{"annotation", "annotation-xml", "maction", "math", "merror", "mfrac", "mi", "mmultiscripts", "mn", "mo", "mover", "mpadded", "mphantom", "mprescripts", "mroot", "mrow", "ms", "mspace", "msqrt", "mstyle", "msub", "msubsup", "msup", "mtable", "mtd", "mtext", "mtr", "munder", "munderover", "semantics"}
)

type templateVars struct {
	Namespace               string
	InHTML, InSVG, InMathML bool
	HasParams, HasChildren  bool
	NumParams, NumChildren  int
}

func (j *jsxTransformer) transform(e *javascript.JSXElement) (*javascript.PrimaryExpression, error) {
	if e.ElementName.Identifier == nil {
		return nil, javascript.ErrMissingIdentifier
	}

	defer j.setNamespace(e.ElementName)()

	inHTML, inSVG, inMathML := nsIn(e.ElementName.Identifier)

	var sb strings.Builder

	if err := j.tmpl.Execute(&sb, templateVars{
		Namespace:   j.namespace,
		InHTML:      inHTML,
		InSVG:       inSVG,
		InMathML:    inMathML,
		HasParams:   len(e.Attributes) > 0,
		HasChildren: len(e.Children) > 0,
		NumParams:   len(e.Attributes),
		NumChildren: len(e.Children),
	}); err != nil {
		return nil, fmt.Errorf("error while executing JSX template: %w", err)
	}

	tk := parser.NewStringTokeniser(sb.String())

	m, err := javascript.ParseModule(&tk)
	if err != nil {
		return nil, fmt.Errorf("error while parsing transformed code: %w", err)
	}

	return j.process(e, m)
}

func nsIn(name *javascript.Token) (bool, bool, bool) {
	return slices.Contains(htmlElements[:], name.Data), slices.Contains(svgElements[:], name.Data), slices.Contains(mathMLElement[:], name.Data)
}

func (j *jsxTransformer) process(e *javascript.JSXElement, m *javascript.Module) (*javascript.PrimaryExpression, error) {
	if err := replaceTagName(m, e); err != nil {
		return nil, err
	}

	s, err := scope.Build(m, nil)
	if err != nil {
		return nil, fmt.Errorf("error building scope: %w", err)
	}

	delete(s.Bindings, params)
	delete(s.Bindings, children)

	j.handleImports(m, s)

	if err := j.replaceParamsAndChildren(m, e); err != nil {
		return nil, err
	}

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
			Comments: [2]javascript.Comments{
				nil,
				append(append(make(javascript.Comments, 0, len(e.Comments[2])+len(e.Comments[3])), e.Comments[2]...), e.Comments[3]...),
			},
			Tokens: e.Tokens,
		},
		Tokens: e.Tokens,
	}, nil
}

func replaceTagName(m *javascript.Module, e *javascript.JSXElement) error {
	var h walk.Handler

	h = walk.HandlerFunc(func(t javascript.Type) error {
		walk.Walk(t, h)

		switch t := t.(type) {
		case *javascript.MemberExpression:
			if pe := t.PrimaryExpression; pe != nil {
				if pe.Literal != nil && pe.Literal.Type == javascript.TokenStringLiteral {
					if str, _ := javascript.Unquote(pe.Literal.Data); str == tagName {
						if e.ElementName.MemberExpression != nil {
							return ErrInvalidTagTemplate
						}

						pe.Literal.Data = strconv.Quote(e.ElementName.Identifier.Data)
						pe.Tokens = e.ElementName.Tokens
						t.Comments[0] = e.Comments[0]
						t.Comments[4] = e.ElementName.Comments[2]
						t.Tokens = e.ElementName.Tokens
					}
				} else if pe.IdentifierReference != nil && pe.IdentifierReference.Data == tagName {
					pe.IdentifierReference = e.ElementName.Identifier
					pe.Tokens = e.ElementName.Tokens
					t.Comments[0] = e.Comments[0]

					for _, ct := range e.ElementName.MemberExpression {
						t.Comments[4] = ct.Comments[0]
						x := *t
						*t = javascript.MemberExpression{
							MemberExpression: &x,
							IdentifierName:   ct.Token,
							Comments:         [5]javascript.Comments{nil, ct.Comments[1]},
						}
					}

					t.Comments[4] = e.ElementName.Comments[2]
					t.Tokens = e.ElementName.Tokens
				}
			}
		case *javascript.ImportClause:
			if t.ImportedDefaultBinding != nil && t.ImportedDefaultBinding.Data == tagName {
				t.ImportedDefaultBinding = e.ElementName.Identifier
			} else if t.NameSpaceImport != nil && t.NameSpaceImport.Data == tagName {
				t.NameSpaceImport = e.ElementName.Identifier
			}
		case *javascript.ImportSpecifier:
			if t.IdentifierName != nil && t.IdentifierName.Data == tagName {
				t.IdentifierName = e.ElementName.Identifier
			}

			if t.ImportedBinding != nil && t.ImportedBinding.Data == tagName {
				t.ImportedBinding = e.ElementName.Identifier
			}
		}

		return nil
	})

	return walk.Walk(m, h)
}

func (j *jsxTransformer) replaceParamsAndChildren(m *javascript.Module, e *javascript.JSXElement) error {
	var h walk.Handler

	h = walk.HandlerFunc(func(t javascript.Type) error {
		walk.Walk(t, h)

		switch t := t.(type) {
		case *javascript.PrimaryExpression:
			if t.IdentifierReference != nil {
				switch t.IdentifierReference.Data {
				case "PARAMS":
					ol, err := j.paramsToObject(e.Attributes)
					if err != nil {
						return err
					}

					t.IdentifierReference = nil
					t.ObjectLiteral = ol
				case "CHILDREN":
					al, err := j.childrenToArray(e.Children)
					if err != nil {
						return err
					}

					t.IdentifierReference = nil
					t.ArrayLiteral = al
				}
			}
		}

		return nil
	})

	return walk.Walk(m, h)
}

func (j *jsxTransformer) handleImports(m *javascript.Module, s *scope.Scope) {
	old := s.Bindings
	s.Bindings = make(map[string][]scope.Binding, len(s.Bindings))

	for _, mli := range m.ModuleListItems {
		if mli.ImportDeclaration == nil {
			continue
		}

		from, _ := javascript.Unquote(mli.ImportDeclaration.FromClause.ModuleSpecifier.Data)

		j.imports[from] = struct{}{}
	}

	for binding, bs := range old {
		s.Bindings[binding] = bs

		if len(bs) < 2 || bs[0].BindingType != scope.BindingImport {
			continue
		}

		s.Rename(binding, getImportID(m, bs[0].Token))
	}
}

func getImportID(m *javascript.Module, tk *javascript.Token) string {
	for _, mli := range m.ModuleListItems {
		if mli.ImportDeclaration == nil {
			continue
		} else if from, _ := javascript.Unquote(mli.ImportDeclaration.FromClause.ModuleSpecifier.Data); mli.ImportDeclaration.ImportedDefaultBinding == tk {
			return "\x00\x00" + from
		} else if mli.ImportDeclaration.NameSpaceImport == tk {
			return "\x00*\x00" + from
		} else if bi := hasIdentifier(mli.ImportDeclaration.NamedImports, tk); bi != "" {
			return "\x00" + bi + "\x00" + from
		}
	}

	return ""
}

func hasIdentifier(ni *javascript.NamedImports, tk *javascript.Token) string {
	for _, i := range ni.ImportList {
		if i.IdentifierName.Data == tk.Data {
			return i.ImportedBinding.Data
		}
	}

	return ""
}

func (j *jsxTransformer) paramsToObject(attrs []javascript.JSXAttribute) (*javascript.ObjectLiteral, error) {
	ol := &javascript.ObjectLiteral{
		PropertyDefinitionList: make([]javascript.PropertyDefinition, 0, len(attrs)),
	}

	for _, attr := range attrs {
		ae, err := j.paramTo(attr)
		if err != nil {
			return nil, err
		}

		if attr.Identifier == nil {
			ol.PropertyDefinitionList = append(ol.PropertyDefinitionList, javascript.PropertyDefinition{
				AssignmentExpression: ae,
			})
		} else {
			ol.PropertyDefinitionList = append(ol.PropertyDefinitionList, javascript.PropertyDefinition{
				PropertyName: &javascript.PropertyName{
					LiteralPropertyName: &javascript.Token{
						Token: parser.Token{
							Type: javascript.TokenStringLiteral,
							Data: strconv.Quote(attr.Identifier.Data),
						},
					},
					Tokens: attr.Tokens,
				},
				AssignmentExpression: ae,
			})
		}

		ol.Comments = [2]javascript.Comments{attr.Comments[0]}
	}

	return ol, nil
}

func (j *jsxTransformer) paramTo(t javascript.JSXAttribute) (*javascript.AssignmentExpression, error) {
	if t.JSXElement != nil {
		pe, err := j.transform(t.JSXElement)
		if err != nil {
			return nil, err
		}

		return &javascript.AssignmentExpression{
			ConditionalExpression: javascript.WrapConditional(pe),
		}, nil
	} else if t.JSXFragment != nil {
		al, err := j.childrenToArray(t.JSXFragment.Children)
		if err != nil {
			return nil, err
		}

		return &javascript.AssignmentExpression{
			ConditionalExpression: javascript.WrapConditional(al),
		}, nil
	} else if t.JSXString != nil {
		str, err := javascript.UnescapeJSXString(t.JSXString.Data)
		if err != nil {
			return nil, err
		}

		return &javascript.AssignmentExpression{
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
		}, nil
	} else if t.AssignmentExpression != nil {
		return t.AssignmentExpression, nil
	}

	return &javascript.AssignmentExpression{
		ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
			Literal: &javascript.Token{
				Token: parser.Token{
					Data: "true",
					Type: javascript.TokenBooleanLiteral,
				},
			},
		}),
	}, nil
}

type importData struct {
	*javascript.ImportDeclaration
	bindings map[string]string
}

// Process transforms any JSX within the given parsed Module using the template
// to generate the required JavaScript.
//
// Within the template, you can use the TAG_NAME placeholder in place of the
// element name. It can be represented as an identifier or as a string literal.
//
// The PARAMS placeholder will be replaced with an object containing the
// parameters.
//
// The CHILDREN placeholder will be replaced with an array of the child
// elements.
//
// In addition, the following template variables are available:
//
//	.Namespace:  Specified namespace, or automatically determined to one of html, svg, mathml.
//	.InHTML:     Set to true if tag name is a known HTML tag.
//	.InSVG       Set to true if tag name is a known SVG tag.
//	.InMathML    Set to true if tag name is a known MathML tag.
//	.HasParams   Set to true if params have been set.
//	.HasChildren Set to true if children have been set.
//	.NumParams   Number of parameters.
//	.NumChildren Number of children.
//
// Any import statement will be added to the Module, with import bindings being
// potentially renamed on a clash.
func Process(m *javascript.Module, tmpl *template.Template) error {
	j := &jsxTransformer{
		tmpl:    tmpl,
		imports: make(map[string]struct{}),
	}

	if err := walk.Walk(m, j); err != nil {
		return err
	}

	imports, err := existingImports(m)
	if err != nil {
		return err
	}

	addImports(m, imports, j.imports)

	s, err := scope.Build(m, nil)
	if err != nil {
		return err
	}

	newIdentsToRename(s, imports)

	return nil
}

func addImports(m *javascript.Module, existing map[string]*importData, imports map[string]struct{}) {
	imps := slices.Collect(maps.Keys(imports))

	slices.Sort(imps)

	for _, from := range imps {
		_, ok := existing[from]
		if ok {
			continue
		}

		id := &javascript.ImportDeclaration{
			FromClause: javascript.FromClause{
				ModuleSpecifier: &javascript.Token{
					Token: parser.Token{
						Data: strconv.Quote(from),
					},
				},
			},
		}

		m.ModuleListItems = slices.Insert(m.ModuleListItems, 0, javascript.ModuleItem{
			ImportDeclaration: id,
		})

		imp := &importData{
			ImportDeclaration: id,
			bindings:          make(map[string]string),
		}
		existing[from] = imp
	}
}

func existingImports(m *javascript.Module) (map[string]*importData, error) {
	imports := make(map[string]*importData)

	for _, mi := range m.ModuleListItems {
		if mi.ImportDeclaration == nil {
			continue
		}

		from, err := javascript.Unquote(mi.ImportDeclaration.FromClause.ModuleSpecifier.Data)
		if err != nil {
			return nil, err
		}

		id := &importData{
			ImportDeclaration: mi.ImportDeclaration,
			bindings:          make(map[string]string),
		}

		if id.ImportedDefaultBinding != nil {
			id.bindings[""] = id.ImportedDefaultBinding.Data
		}

		if id.NameSpaceImport != nil {
			id.bindings["*"] = id.NameSpaceImport.Data
		}

		if id.NamedImports != nil {
			for _, ni := range id.NamedImports.ImportList {
				id.bindings[ni.IdentifierName.Data] = ni.ImportedBinding.Data
			}
		}

		imports[from] = id
	}

	return imports, nil
}

func newIdentsToRename(s *scope.Scope, imports map[string]*importData) {
	b := make(map[string][]scope.Binding, len(s.Bindings))

	var rename, bindings []string

	for binding, bs := range s.Bindings {
		if strings.HasPrefix(binding, "\x00") {
			bindings = append(bindings, binding)
		} else {
			b[binding] = bs
		}
	}

	slices.Sort(bindings)

	for _, binding := range bindings {
		ident, from, _ := strings.Cut(binding[1:], "\x00")
		imp := imports[from]

		if imp.ImportClause == nil {
			imp.ImportClause = new(javascript.ImportClause)
		}

		ni := imp.getExistingBindingOrMake(s, b, binding, ident)

		b[ni] = append(b[ni], s.Bindings[binding]...)
		rename = append(rename, ni)
	}

	s.Bindings = b

	renameNewBindings(s, rename)
}

func (imp *importData) getExistingBindingOrMake(s *scope.Scope, b map[string][]scope.Binding, binding, ident string) string {
	ni, ok := imp.bindings[ident]
	if !ok {
		tk := &javascript.Token{
			Token: parser.Token{
				Type: javascript.TokenIdentifier,
				Data: ident,
			},
		}
		imp.bindings[ident] = tk.Data

		switch ident {
		case "":
			imp.ImportedDefaultBinding = tk
		case "*":
			imp.NameSpaceImport = tk
		default:
			if imp.NamedImports == nil {
				imp.NamedImports = new(javascript.NamedImports)
			}

			imp.NamedImports.ImportList = append(imp.NamedImports.ImportList, javascript.ImportSpecifier{
				IdentifierName: &javascript.Token{
					Token: parser.Token{
						Type: javascript.TokenIdentifier,
						Data: ident,
					},
				},
				ImportedBinding: tk,
			})

			slices.SortFunc(imp.NamedImports.ImportList, func(a, b javascript.ImportSpecifier) int {
				return strings.Compare(a.ImportedBinding.Data, b.ImportedBinding.Data)
			})
		}

		ni = binding
		b[ni] = append(b[ni], scope.Binding{
			BindingType: scope.BindingImport,
			Scope:       s,
			Token:       tk,
		})
	}

	return ni
}

func renameNewBindings(s *scope.Scope, rename []string) {
	for _, name := range rename {
		s.Rename(name, "\x00")

		if strings.HasPrefix(name, "\x00") {
			name, _, _ = strings.Cut(name[1:], "\x00")
		}

		switch name {
		case "":
			name = "def"
		case "*":
			name = "ns"
		}

		num := 0
		newName := name

		for s.IdentifierInUse(newName) {
			num++
			newName = name + "_" + strconv.Itoa(num)
		}

		s.Rename("\x00", newName)
	}
}

var (
	ErrInvalidTransformation = errors.New("invalid transformation")
	ErrTooManyStatements     = errors.New("too many statements")
	ErrMissingChild          = errors.New("missing JSX child")
	ErrInvalidTagTemplate    = errors.New("cannot convert member expression JSX Element name to string")
)
