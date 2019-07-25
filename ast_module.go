package javascript

import (
	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

// Module represents the top-level of a parsed javascript module
type Module struct {
	ModuleListItems []ModuleItem
	Tokens          Tokens
}

// ParseModule parses a javascript module
func ParseModule(t parser.Tokeniser) (*Module, error) {
	j, err := newJSParser(t)
	if err != nil {
		return nil, err
	}
	m := new(Module)
	if err := m.parse(&j); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Module) parse(j *jsParser) error {
	for j.AcceptRunWhitespace() != parser.TokenDone {
		g := j.NewGoal()
		ml := len(m.ModuleListItems)
		m.ModuleListItems = append(m.ModuleListItems, ModuleItem{})
		if err := m.ModuleListItems[ml].parse(&g); err != nil {
			return err
		}
		j.Score(g)
	}
	m.Tokens = j.ToTokens()
	return nil
}

// ModuleItem as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-ModuleItem
//
// Only one of ImportDeclaration, StatementListItem, or ExportDeclaration must
// be non-nil.
type ModuleItem struct {
	ImportDeclaration *ImportDeclaration
	StatementListItem *StatementListItem
	ExportDeclaration *ExportDeclaration
	Tokens            Tokens
}

func (ml *ModuleItem) parse(j *jsParser) error {
	g := j.NewGoal()
	switch g.Peek() {
	case parser.Token{TokenKeyword, "import"}:
		ml.ImportDeclaration = new(ImportDeclaration)
		if err := ml.ImportDeclaration.parse(&g); err != nil {
			return j.Error("ModuleItem", err)
		}
	case parser.Token{TokenKeyword, "export"}:
		ml.ExportDeclaration = new(ExportDeclaration)
		if err := ml.ExportDeclaration.parse(&g); err != nil {
			return j.Error("ModuleItem", err)
		}
	default:
		ml.StatementListItem = new(StatementListItem)
		if err := ml.StatementListItem.parse(&g, false, false, false); err != nil {
			return j.Error("ModuleItem", err)
		}
	}
	j.Score(g)
	ml.Tokens = j.ToTokens()
	return nil
}

// ImportDeclaration as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-ImportDeclaration
type ImportDeclaration struct {
	*ImportClause
	FromClause
	Tokens Tokens
}

func (id *ImportDeclaration) parse(j *jsParser) error {
	if !j.AcceptToken(parser.Token{TokenKeyword, "import"}) {
		return j.Error("ImportDeclaration", ErrInvalidImport)
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	if g.Accept(TokenStringLiteral) {
		id.FromClause.Tokens = g.ToTokens()
		id.ModuleSpecifier = &id.FromClause.Tokens[0]
	} else {
		id.ImportClause = new(ImportClause)
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
	if !j.parseSemicolon() {
		return j.Error("ImportDeclaration", ErrMissingSemiColon)
	}
	id.Tokens = j.ToTokens()
	return nil
}

// ImportClause as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-ImportClause
//
// At least one of ImportedDefaultBinding, NameSpaceImport, and NamedImports
// must be non-nil.
//
// Both NameSpaceImport and NamedImports can not be non-nil.
type ImportClause struct {
	ImportedDefaultBinding *Token
	NameSpaceImport        *Token
	NamedImports           *NamedImports
	Tokens                 Tokens
}

func (ic *ImportClause) parse(j *jsParser) error {
	if t := j.Peek().Type; t == TokenIdentifier || t == TokenKeyword {
		g := j.NewGoal()
		if ic.ImportedDefaultBinding = g.parseIdentifier(false, false); ic.ImportedDefaultBinding == nil {
			return j.Error("ImportClause", ErrNoIdentifier)
		}
		j.Score(g)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			ic.Tokens = j.ToTokens()
			return nil
		}
		g.AcceptRunWhitespace()
		j.Score(g)
	}
	if j.Peek() == (parser.Token{TokenPunctuator, "*"}) {
		j.Skip()
		j.AcceptRunWhitespace()
		if !j.AcceptToken(parser.Token{TokenIdentifier, "as"}) {
			return j.Error("ImportClause", ErrInvalidNameSpaceImport)
		}
		j.AcceptRunWhitespace()
		if ic.NameSpaceImport = j.parseIdentifier(false, false); ic.NameSpaceImport == nil {
			return j.Error("ImportClause", ErrNoIdentifier)
		}
	} else if j.Peek() == (parser.Token{TokenPunctuator, "{"}) {
		g := j.NewGoal()
		ic.NamedImports = new(NamedImports)
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

// FromClause as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-FromClause
//
// ModuleSpecifier must be non-nil.
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

// NamedImports as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-NamedImports
type NamedImports struct {
	ImportList []ImportSpecifier
	Tokens     Tokens
}

func (ni *NamedImports) parse(j *jsParser) error {
	if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		return j.Error("NamedImports", ErrInvalidNamedImport)
	}
	for {
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		}
		g := j.NewGoal()
		is := len(ni.ImportList)
		ni.ImportList = append(ni.ImportList, ImportSpecifier{})
		if err := ni.ImportList[is].parse(&g); err != nil {
			return j.Error("NamedImports", err)
		}
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

// ImportSpecifier as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-ImportSpecifier
//
// ImportedBinding must be non-nil.
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
			if is.ImportedBinding = g.parseIdentifier(false, false); is.ImportedBinding == nil {
				return g.Error("ImportSpecifier", ErrNoIdentifier)
			}
			j.Score(g)
		}
	}
	is.Tokens = j.ToTokens()
	return nil
}

// ExportDeclaration as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-ExportDeclaration
//
// It is only valid for one of ExportClause, VariableStatement, Declaration,
// DefaultFunction, DefaultClass, or DefaultAssignmentExpression to be non-nil.
//
// FromClause can be non-nil exclusively or paired with ExportClause.
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
	if !j.AcceptToken(parser.Token{TokenKeyword, "export"}) {
		return j.Error("ExportDeclaration", ErrInvalidExportDeclaration)
	}
	j.AcceptRunWhitespace()
	if j.AcceptToken(parser.Token{TokenKeyword, "default"}) {
		j.AcceptRunWhitespace()
		tk := j.Peek()
		g := j.NewGoal()
		switch tk.Data {
		case "async", "function":
			ed.DefaultFunction = new(FunctionDeclaration)
			if err := ed.DefaultFunction.parse(&g, false, false, true); err != nil {
				return j.Error("ExportDeclaration", err)
			}
			j.Score(g)
		case "class":
			ed.DefaultClass = new(ClassDeclaration)
			if err := ed.DefaultClass.parse(&g, false, false, true); err != nil {
				return j.Error("ExportDeclaration", err)
			}
			j.Score(g)
		default:
			ed.DefaultAssignmentExpression = new(AssignmentExpression)
			if err := ed.DefaultAssignmentExpression.parse(&g, true, false, false); err != nil {
				return j.Error("ExportDeclaration", err)
			}
			j.Score(g)
			if !j.parseSemicolon() {
				return j.Error("ExportDeclaration", ErrMissingSemiColon)
			}
		}
	} else if j.AcceptToken(parser.Token{TokenPunctuator, "*"}) {
		j.AcceptRunWhitespace()
		g := j.NewGoal()
		ed.FromClause = new(FromClause)
		if err := ed.FromClause.parse(&g); err != nil {
			return j.Error("ExportDeclaration", err)
		}
		j.Score(g)
		if !j.parseSemicolon() {
			return j.Error("ExportDeclaration", ErrMissingSemiColon)
		}
	} else if g := j.NewGoal(); g.Peek() == (parser.Token{TokenKeyword, "var"}) {
		ed.VariableStatement = new(VariableStatement)
		if err := ed.VariableStatement.parse(&g, false, false); err != nil {
			return j.Error("ExportDeclaration", err)
		}
		j.Score(g)
	} else if g.Peek() == (parser.Token{TokenPunctuator, "{"}) {
		ed.ExportClause = new(ExportClause)
		if err := ed.ExportClause.parse(&g); err != nil {
			return j.Error("ExportDeclaration", err)
		}
		j.Score(g)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if g.Peek() == (parser.Token{TokenIdentifier, "from"}) {
			h := g.NewGoal()
			h.Skip()
			h.AcceptRunWhitespace()
			if h.Accept(TokenStringLiteral) {
				h = g.NewGoal()
				ed.FromClause = new(FromClause)
				ed.FromClause.parse(&h)
				g.Score(h)
				j.Score(g)
			}
		}
		if !j.parseSemicolon() {
			return j.Error("ExportDeclaration", ErrMissingSemiColon)
		}
	} else {
		ed.Declaration = new(Declaration)
		if err := ed.Declaration.parse(&g, false, false); err != nil {
			return j.Error("ExportDeclaration", err)
		}
		j.Score(g)
	}
	ed.Tokens = j.ToTokens()
	return nil
}

// ExportClause as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-ExportClause
type ExportClause struct {
	ExportList []ExportSpecifier
	Tokens     Tokens
}

func (ec *ExportClause) parse(j *jsParser) error {
	if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		return j.Error("ExportClause", ErrInvalidExportClause)
	}
	for {
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		}
		g := j.NewGoal()
		es := len(ec.ExportList)
		ec.ExportList = append(ec.ExportList, ExportSpecifier{})
		if err := ec.ExportList[es].parse(&g); err != nil {
			return j.Error("ExportClause", err)
		}
		j.Score(g)
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

// ExportSpecifier as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-ExportSpecifier
//
// IdentifierName must be non-nil
type ExportSpecifier struct {
	IdentifierName, EIdentifierName *Token
	Tokens                          Tokens
}

func (es *ExportSpecifier) parse(j *jsParser) error {
	if !j.Accept(TokenIdentifier, TokenKeyword) {
		return j.Error("ExportSpecifier", ErrNoIdentifier)
	}
	es.IdentifierName = j.GetLastToken()
	g := j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenIdentifier, "as"}) {
		j.Score(g)
		j.AcceptRunWhitespace()
		if !j.Accept(TokenIdentifier, TokenKeyword) {
			return j.Error("ExportSpecifier", ErrNoIdentifier)
		}
		es.EIdentifierName = j.GetLastToken()
	}
	es.Tokens = j.ToTokens()
	return nil
}

// Errors
var (
	ErrInvalidImport            = errors.New("invalid import statement")
	ErrInvalidExportDeclaration = errors.New("invalid export declaration")
	ErrInvalidNameSpaceImport   = errors.New("invalid namespace import")
	ErrMissingFrom              = errors.New("missing from")
	ErrMissingModuleSpecifier   = errors.New("missing module specifier")
	ErrInvalidNamedImport       = errors.New("invalid named import list")
	ErrInvalidImportSpecifier   = errors.New("invalid import specifier")
	ErrInvalidExportClause      = errors.New("invalid export clause")
)
