package jsx

import (
	"errors"
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

func (j *jsxTransformer) transform(e *javascript.JSXElement) (*javascript.PrimaryExpression, error) {
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

func (j *jsxTransformer) process(e *javascript.JSXElement, m *javascript.Module) (*javascript.PrimaryExpression, error) {
	replaceTagName(m, e.ElementName.Identifier.Data)

	s, err := scope.Build(m, nil)
	if err != nil {
		return nil, err
	}

	delete(s.Bindings, params)
	delete(s.Bindings, children)

	j.handleImports(m, s)
	replaceParamsAndChildren(m, e)

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
	var h walk.Handler

	h = walk.HandlerFunc(func(t javascript.Type) error {
		walk.Walk(t, h)

		switch t := t.(type) {
		case *javascript.PrimaryExpression:
			if t.Literal != nil && t.Literal.Type == javascript.TokenStringLiteral {
				if str, _ := javascript.Unquote(t.Literal.Data); str == tagName {
					t.Literal.Data = strconv.Quote(name)
				}
			} else if t.IdentifierReference != nil && t.IdentifierReference.Data == tagName {
				t.IdentifierReference.Data = name
			}
		case *javascript.ImportClause:
			if t.ImportedDefaultBinding != nil && t.ImportedDefaultBinding.Data == tagName {
				t.ImportedDefaultBinding.Data = name
			} else if t.NameSpaceImport != nil && t.NameSpaceImport.Data == tagName {
				t.NameSpaceImport.Data = name
			}
		case *javascript.ImportSpecifier:
			if t.IdentifierName != nil && t.IdentifierName.Data == tagName {
				t.IdentifierName.Data = name
			}

			if t.ImportedBinding != nil && t.ImportedBinding.Data == tagName {
				t.ImportedBinding.Data = name
			}
		}

		return nil
	})

	walk.Walk(m, h)
}

func replaceParamsAndChildren(m *javascript.Module, e *javascript.JSXElement) {
	var h walk.Handler

	h = walk.HandlerFunc(func(t javascript.Type) error {
		walk.Walk(t, h)

		switch t := t.(type) {
		case *javascript.PrimaryExpression:
			if t.IdentifierReference != nil {
				switch t.IdentifierReference.Data {
				case "PARAMS":
					t.IdentifierReference = nil
					t.ObjectLiteral = paramsToObject(e.Attributes)
				case "CHILDREN":
					t.IdentifierReference = nil
					t.ArrayLiteral = childrenToArray(e.Children)
				}
			}
		}

		return nil
	})

	walk.Walk(m, h)
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
		}

		from, err := javascript.Unquote(mli.ImportDeclaration.FromClause.ModuleSpecifier.Data)
		if err != nil {
			continue
		}

		if mli.ImportDeclaration.ImportedDefaultBinding == tk {
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

func paramsToObject(attrs []javascript.JSXAttribute) *javascript.ObjectLiteral {
	ol := &javascript.ObjectLiteral{
		PropertyDefinitionList: make([]javascript.PropertyDefinition, 0, len(attrs)),
	}

	for _, attr := range attrs {
		if attr.Identifier == nil {
			if attr.AssignmentExpression == nil {
				return nil
			}

			ol.PropertyDefinitionList = append(ol.PropertyDefinitionList, javascript.PropertyDefinition{
				AssignmentExpression: attr.AssignmentExpression,
			})
		} else {
			ol.PropertyDefinitionList = append(ol.PropertyDefinitionList, javascript.PropertyDefinition{
				PropertyName: &javascript.PropertyName{
					LiteralPropertyName: attr.Identifier,
				},
				AssignmentExpression: attr.AssignmentExpression,
			})
		}
	}

	return ol
}

type importData struct {
	*javascript.ImportDeclaration
	bindings map[string]string
}

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

	for _, from := range imps {
		imp, ok := existing[from]
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

		imp = &importData{
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

		ni, ok := imp.bindings[ident]
		if !ok {
			tk := &javascript.Token{
				Token: parser.Token{
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

		b[ni] = append(b[ni], s.Bindings[binding]...)
		rename = append(rename, ni)
	}

	s.Bindings = b

	renameNewBindings(s, rename)
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
	ErrTooManyStatements     = errors.New("too many statments")
)
