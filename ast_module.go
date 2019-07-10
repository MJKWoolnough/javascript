package javascript

import (
	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

type Module struct {
	ModuleListItems []ModuleListItem
	Tokens          Tokens
}

func ParseModule(t parser.Tokeniser) (*Module, error) {
	j, err := newJSParser(t)
	if err != nil {
		return nil, err
	}
	m := newModule()
	if err := m.parse(&j); err != nil {
		m.clear()
		return nil, err
	}
	return m, nil
}

func (m *Module) parse(j *jsParser) error {
	for j.AcceptRunWhitespace() != parser.TokenDone {
		g := j.NewGoal()
		var ml ModuleListItem
		if err := ml.parse(&g); err != nil {
			m.clear()
			return j.Error("Module", err)
		}
		j.Score(g)
		m.ModuleListItems = append(m.ModuleListItems, ml)
	}
	m.Tokens = j.ToTokens()
	return nil
}

type ModuleListItem struct {
	ImportDeclaration *ImportDeclaration
	StatementListItem *StatementListItem
	ExportDeclaration *ExportDeclaration
	Tokens            Tokens
}

func (ml *ModuleListItem) parse(j *jsParser) error {
	g := j.NewGoal()
	switch g.Peek() {
	case parser.Token{TokenKeyword, "import"}:
		ml.ImportDeclaration = newImportDeclaration()
		if err := ml.ImportDeclaration.parse(&g); err != nil {
			return j.Error("ModuleStatement", err)
		}
	case parser.Token{TokenKeyword, "export"}:
		ml.ExportDeclaration = newExportDeclaration()
		if err := ml.ExportDeclaration.parse(&g); err != nil {
			return j.Error("ModuleStatement", err)
		}
	default:
		ml.StatementListItem = newStatementListItem()
		if err := ml.StatementListItem.parse(&g, false, false, false); err != nil {
			return j.Error("ModuleStatement", err)
		}
	}
	j.Score(g)
	ml.Tokens = j.ToTokens()
	return nil
}

type ImportDeclaration struct {
	*ImportClause
	FromClause
	Tokens Tokens
}

func (id *ImportDeclaration) parse(j *jsParser) error {
	j.AcceptToken(parser.Token{TokenKeyword, "import"})
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	if g.Accept(TokenStringLiteral) {
		id.FromClause.Tokens = g.ToTokens()
		id.ModuleSpecifier = &id.FromClause.Tokens[0]
	} else {
		id.ImportClause = newImportClause()
		if err := id.ImportClause.parse(&g); err != nil {
			return j.Error("ImportDeclaration", err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
		g = j.NewGoal()
		if err := id.FromClause.parse(&g); err != nil {
			return j.Error("ImportDeclaration", err)
		}
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
		return j.Error("ImportDeclaration", ErrMissingSemiColon)
	}
	id.Tokens = j.ToTokens()
	return nil
}

type ImportClause struct {
	ImportedDefaultBinding *Token
	NameSpaceImport        *Token
	NamedImports           *NamedImports
	Tokens                 Tokens
}

func (ic *ImportClause) parse(j *jsParser) error {
	if t := j.Peek().Type; t == TokenIdentifier || t == TokenKeyword {
		g := j.NewGoal()
		ib, err := g.parseIdentifier(false, false)
		if err != nil {
			return j.Error("ImportClause", err)
		}
		j.Score(g)
		ic.ImportedDefaultBinding = ib
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			ic.Tokens = j.ToTokens()
			return nil
		}
		g.AcceptRunWhitespace()
		j.Score(g)
		g = j.NewGoal()
	}
	if j.Peek() == (parser.Token{TokenPunctuator, "*"}) {
		g := j.NewGoal()
		g.Except()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenIdentifier, "as"}) {
			return j.Error("ImportClause", ErrInvalidNameSpaceImport)
		}
		g.AcceptRunWhitespace()
		h := g.NewGoal()
		ib, err := h.parseIdentifier(false, false)
		if err != nil {
			return g.Error("ImportClause", err)
		}
		g.Score(h)
		j.Score(g)
		ic.NameSpaceImport = ib
	} else if j.Peek() == (parser.Token{TokenPunctuator, "{"}) {
		g := j.NewGoal()
		ic.NamedImports = newNamedImports()
		if err := ic.NamedImports.parse(&g); err != nil {
			return j.Error("ImportClause", err)
		}
		j.Score(g)
	} else {
		return j.Error("ImportClause", ErrInvalidImport)
	}
	ic.Tokens = j.ToTokens()
	return nil
}

type FromClause struct {
	ModuleSpecifier *Token
	Tokens          Tokens
}

func (fc *FromClause) parse(j *jsParser) error {
	if !j.AcceptToken(parser.Token{TokenIdentifier, "from"}) {
		return j.Error("FromClause", ErrMissingFrom)
	}
	j.AcceptRunWhitespace()
	if !j.Accept(TokenStringLiteral) {
		return j.Error("FromClause", ErrMissingModuleSpecifier)
	}
	fc.Tokens = j.ToTokens()
	fc.ModuleSpecifier = &fc.Tokens[len(fc.Tokens)-1]
	return nil
}

type NamedImports struct {
	ImportList []ImportSpecifier
	Tokens     Tokens
}

func (ni *NamedImports) parse(j *jsParser) error {
	j.AcceptToken(parser.Token{TokenPunctuator, "{"})
	for {
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		}
		g := j.NewGoal()
		var is ImportSpecifier
		if err := is.parse(&g); err != nil {
			is.clear()
			return j.Error("NamedImports", err)
		}
		ni.ImportList = append(ni.ImportList, is)
		j.Score(g)
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return j.Error("NamedImports", ErrInvalidNamedImport)
		}
	}
	ni.Tokens = j.ToTokens()
	return nil
}

type ImportSpecifier struct {
	IdentifierName  *Token
	ImportedBinding *Token
	Tokens          Tokens
}

func (is *ImportSpecifier) parse(j *jsParser) error {
	if !j.Accept(TokenIdentifier, TokenKeyword) {
		return j.Error("ImportSpecifier", ErrInvalidImportSpecifier)
	}
	is.ImportedBinding = j.GetLastToken()
	if is.ImportedBinding.Type == TokenIdentifier || is.ImportedBinding.Data == "yield" || is.ImportedBinding.Data == "await" {
		g := j.NewGoal()
		g.AcceptRunWhitespace()
		if g.AcceptToken(parser.Token{TokenIdentifier, "as"}) {
			is.IdentifierName = is.ImportedBinding
			g.AcceptRunWhitespace()
			var err error
			if is.ImportedBinding, err = g.parseIdentifier(false, false); err != nil {
				return j.Error("ImportSpecifier", err)
			}
			j.Score(g)
		}
	}
	is.Tokens = j.ToTokens()
	return nil
}

type ExportDeclaration struct {
	ExportClause                *ExportClause
	FromClause                  *FromClause
	VariableStatement           *VariableStatement
	Declaration                 *Declaration
	DefaultFunction             *FunctionDeclaration
	DefaultClass                *ClassDeclaration
	DefaultAssignmentExpression *AssignmentExpression
	Tokens                      Tokens
}

func (ed *ExportDeclaration) parse(j *jsParser) error {
	j.AcceptToken(parser.Token{TokenKeyword, "export"})
	j.AcceptRunWhitespace()
	if j.AcceptToken(parser.Token{TokenKeyword, "default"}) {
		j.AcceptRunWhitespace()
		tk := j.Peek()
		g := j.NewGoal()
		switch tk.Data {
		case "async", "function":
			ed.DefaultFunction = newFunctionDeclaration()
			if err := ed.DefaultFunction.parse(&g, false, false, true); err != nil {
				return j.Error("ExportDeclaration", err)
			}
			j.Score(g)
		case "class":
			ed.DefaultClass = newClassDeclaration()
			if err := ed.DefaultClass.parse(&g, false, false, true); err != nil {
				return j.Error("ExportDeclaration", err)
			}
			j.Score(g)
		default:
			ed.DefaultAssignmentExpression = newAssignmentExpression()
			if err := ed.DefaultAssignmentExpression.parse(&g, true, false, false); err != nil {
				return j.Error("ExportDeclaration", err)
			}
			j.Score(g)
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
				return j.Error("ExportDeclaration", ErrMissingSemiColon)
			}
		}
	} else if j.AcceptToken(parser.Token{TokenPunctuator, "*"}) {
		j.AcceptRunWhitespace()
		g := j.NewGoal()
		ed.FromClause = newFromClause()
		if err := ed.FromClause.parse(&g); err != nil {
			return j.Error("ExportDeclaration", err)
		}
		j.Score(g)
		if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			return j.Error("ExportDeclaration", ErrMissingSemiColon)
		}
	} else if g := j.NewGoal(); g.Peek() == (parser.Token{TokenKeyword, "var"}) {
		ed.VariableStatement = newVariableStatement()
		if err := ed.VariableStatement.parse(&g, false, false); err != nil {
			return j.Error("ExportDeclaration", err)
		}
		j.Score(g)
	} else if g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		ed.ExportClause = newExportClause()
		if err := ed.ExportClause.parse(&g); err != nil {
			return j.Error("ExportDeclaration", err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
		if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			g = j.NewGoal()
			ed.FromClause = newFromClause()
			if err := ed.FromClause.parse(&g); err != nil {
				return j.Error("ExportDeclaration", err)
			}
			j.Score(g)
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
				return j.Error("ExportDeclaration", ErrMissingSemiColon)
			}
		}
	} else {
		ed.Declaration = newDeclaration()
		if err := ed.Declaration.parse(&g, false, false); err != nil {
			return j.Error("ExportDeclaration", err)
		}
		j.Score(g)
	}
	ed.Tokens = j.ToTokens()
	return nil
}

type ExportClause struct {
	ExportList []ExportSpecifier
	Tokens     Tokens
}

func (ec *ExportClause) parse(j *jsParser) error {
	for {
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		}
		g := j.NewGoal()
		var es ExportSpecifier
		if err := es.parse(&g); err != nil {
			ec.clear()
			poolExportSpecifier.Put(ec)
			return j.Error("ExportClause", err)
		}
		j.Score(g)
		ec.ExportList = append(ec.ExportList, es)
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return j.Error("ExportClause", ErrInvalidExportClause)
		}
	}
	ec.Tokens = j.ToTokens()
	return nil
}

type ExportSpecifier struct {
	IdentifierName, EIdentifierName *Token
	Tokens                          Tokens
}

func (es *ExportSpecifier) parse(j *jsParser) error {
	if !j.Accept(TokenIdentifier, TokenKeyword) {
		return j.Error("ExportClause", ErrMissingIdentifier)
	}
	es.IdentifierName = j.GetLastToken()
	g := j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenIdentifier, "as"}) {
		g.AcceptRunWhitespace()
		if !g.Accept(TokenIdentifier, TokenKeyword) {
			return j.Error("ExportClause", ErrMissingIdentifier)
		}
		j.Score(g)
		es.EIdentifierName = j.GetLastToken()
	}
	es.Tokens = j.ToTokens()
	return nil
}

var (
	ErrInvalidImport          = errors.New("invalid import statement")
	ErrInvalidNameSpaceImport = errors.New("invalid namespace import")
	ErrMissingFrom            = errors.New("missing from")
	ErrMissingModuleSpecifier = errors.New("missing module specifier")
	ErrInvalidNamedImport     = errors.New("invalid named import list")
	ErrInvalidImportSpecifier = errors.New("invalid import specifier")
	ErrMissingIdentifier      = errors.New("missing identifier")
	ErrInvalidExportClause    = errors.New("invalid export clause")
)
