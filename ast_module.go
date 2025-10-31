package javascript

import (
	"vimagination.zapto.org/parser"
)

// Module represents the top-level of a parsed javascript module
type Module struct {
	ModuleListItems []ModuleItem
	Comments        [2]Comments
	Tokens          Tokens
}

// ParseModule parses a javascript module
func ParseModule(t Tokeniser) (*Module, error) {
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
	m.Comments[0] = j.AcceptRunWhitespaceNoNewlineComments()
	g := j.NewGoal()

	for g.AcceptRunWhitespace() != parser.TokenDone {
		j.AcceptRunWhitespaceNoComment()

		g = j.NewGoal()
		ml := len(m.ModuleListItems)

		m.ModuleListItems = append(m.ModuleListItems, ModuleItem{})
		if err := m.ModuleListItems[ml].parse(&g); err != nil {
			return err
		}

		j.Score(g)

		g = j.NewGoal()
	}

	m.Comments[1] = j.AcceptRunWhitespaceComments()
	m.Tokens = j.ToTokens()

	return nil
}

// ModuleItem as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ModuleItem
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

	g.AcceptRunWhitespace()

	switch g.Peek() {
	case parser.Token{Type: TokenKeyword, Data: "export"}:
		if g.SkipExportType() {
			ml.StatementListItem = &StatementListItem{
				Statement: &Statement{
					Tokens: g.ToTokens(),
				},
				Tokens: g.ToTokens(),
			}

			break
		}

		g = j.NewGoal()

		ml.ExportDeclaration = new(ExportDeclaration)
		if err := ml.ExportDeclaration.parse(&g); err != nil {
			return j.Error("ModuleItem", err)
		}
	case parser.Token{Type: TokenKeyword, Data: "import"}:
		if g.SkipImportType() {
			ml.StatementListItem = &StatementListItem{
				Statement: &Statement{
					Tokens: g.ToTokens(),
				},
				Tokens: g.ToTokens(),
			}

			break
		}

		h := g.NewGoal()

		h.Skip()
		h.AcceptRunWhitespace()

		if !h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "."}) && !h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
			g = j.NewGoal()
			ml.ImportDeclaration = new(ImportDeclaration)

			if err := ml.ImportDeclaration.parse(&g); err != nil {
				return j.Error("ModuleItem", err)
			}

			break
		}

		fallthrough
	default:
		ml.StatementListItem = new(StatementListItem)
		if err := ml.StatementListItem.parse(&g, false, true, false); err != nil {
			return j.Error("ModuleItem", err)
		}
	}

	j.Score(g)

	ml.Tokens = j.ToTokens()

	return nil
}

// ImportDeclaration as defined in ECMA-262
// https://tc39.es/ecma262/#prod-ImportDeclaration
type ImportDeclaration struct {
	*ImportClause
	FromClause
	*WithClause
	Comments [4]Comments
	Tokens   Tokens
}

func (id *ImportDeclaration) parse(j *jsParser) error {
	id.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "import"}) {
		return j.Error("ImportDeclaration", ErrInvalidImport)
	}

	g := j.NewGoal()

	g.AcceptRunWhitespace()

	if g.Accept(TokenStringLiteral) {
		id.Comments[1] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()

		g = j.NewGoal()

		g.Next()

		id.FromClause.Tokens = g.ToTokens()
		id.ModuleSpecifier = &id.FromClause.Tokens[0]
	} else {
		j.AcceptRunWhitespaceNoComment()

		g = j.NewGoal()

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

	g = j.NewGoal()

	g.AcceptRunWhitespace()

	if g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "with"}) {
		j.AcceptRunWhitespaceNoComment()

		g = j.NewGoal()
		id.WithClause = new(WithClause)

		if err := id.WithClause.parse(&g); err != nil {
			return j.Error("ImportDeclaration", err)
		}

		j.Score(g)
	}

	id.Comments[2] = j.AcceptRunWhitespaceNoNewlineComments()

	if !j.parseSemicolon() {
		return j.Error("ImportDeclaration", ErrMissingSemiColon)
	}

	id.Comments[3] = j.AcceptRunWhitespaceNoNewlineComments()
	id.Tokens = j.ToTokens()

	return nil
}

// WithClause as defined in ECMA-262
// https://tc39.es/ecma262/#prod-WithClause
type WithClause struct {
	WithEntries []WithEntry
	Comments    [4]Comments
	Tokens      Tokens
}

func (w *WithClause) parse(j *jsParser) error {
	w.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()
	j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "with"})

	w.Comments[1] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		return j.Error("WithClause", ErrMissingOpeningBrace)
	}

	w.Comments[2] = j.AcceptRunWhitespaceNoNewlineComments()

	g := j.NewGoal()

	g.AcceptRunWhitespace()

	for !g.Accept(TokenRightBracePunctuator) {
		j.AcceptRunWhitespaceNoComment()

		g = j.NewGoal()

		var we WithEntry

		if err := we.parse(&g); err != nil {
			return j.Error("WithClause", err)
		}

		w.WithEntries = append(w.WithEntries, we)

		j.Score(g)

		g = j.NewGoal()

		g.AcceptRunWhitespace()

		if g.Accept(TokenRightBracePunctuator) {
			break
		} else if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			return g.Error("WithClause", ErrMissingComma)
		}

		j.Score(g)

		g = j.NewGoal()

		g.AcceptRunWhitespace()
	}

	w.Comments[3] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()
	j.Next()

	w.Tokens = j.ToTokens()

	return nil
}

// WithEntry as defined in ECMA-262
// https://tc39.es/ecma262/#prod-WithEntries
type WithEntry struct {
	AttributeKey *Token
	Value        *Token
	Comments     [4]Comments
	Tokens       Tokens
}

func (w *WithEntry) parse(j *jsParser) error {
	w.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	g := j.NewGoal()

	if !j.Accept(TokenIdentifier, TokenStringLiteral) {
		return j.Error("WithEntry", ErrMissingAttributeKey)
	}

	w.AttributeKey = j.GetLastToken()

	j.Score(g)

	w.Comments[1] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
		return j.Error("WithEntry", ErrMissingColon)
	}

	w.Comments[2] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.Accept(TokenStringLiteral) {
		return j.Error("WithEntry", ErrMissingString)
	}

	w.Value = j.GetLastToken()

	g = j.NewGoal()

	g.AcceptRunWhitespace()

	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
		w.Comments[3] = j.AcceptRunWhitespaceComments()
	} else {
		w.Comments[3] = j.AcceptRunWhitespaceNoNewlineComments()
	}

	w.Tokens = j.ToTokens()

	return nil
}

// ImportClause as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ImportClause
//
// At least one of ImportedDefaultBinding, NameSpaceImport, and NamedImports
// must be non-nil.
//
// Both NameSpaceImport and NamedImports can not be non-nil.
type ImportClause struct {
	ImportedDefaultBinding *Token
	NameSpaceImport        *Token
	NamedImports           *NamedImports
	Comments               [6]Comments
	Tokens                 Tokens
}

func (ic *ImportClause) parse(j *jsParser) error {
	ic.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if t := j.Peek().Type; t == TokenIdentifier || t == TokenKeyword {
		g := j.NewGoal()

		if ic.ImportedDefaultBinding = g.parseIdentifier(false, false); ic.ImportedDefaultBinding == nil {
			return j.Error("ImportClause", ErrNoIdentifier)
		}

		j.Score(g)

		ic.Comments[1] = j.AcceptRunWhitespaceComments()

		g = j.NewGoal()

		g.AcceptRunWhitespace()

		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			ic.Comments[5] = ic.Comments[1]
			ic.Comments[1] = nil
			ic.Tokens = j.ToTokens()

			return nil
		}

		j.Score(g)

		ic.Comments[2] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()
	}

	if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "*"}) {
		ic.Comments[3] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()

		if !j.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "as"}) {
			return j.Error("ImportClause", ErrInvalidNameSpaceImport)
		}

		ic.Comments[4] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()

		if ic.NameSpaceImport = j.parseIdentifier(false, false); ic.NameSpaceImport == nil {
			return j.Error("ImportClause", ErrNoIdentifier)
		}
	} else if j.Peek() == (parser.Token{Type: TokenPunctuator, Data: "{"}) {
		g := j.NewGoal()

		ic.NamedImports = new(NamedImports)
		if err := ic.NamedImports.parse(&g); err != nil {
			return j.Error("ImportClause", err)
		}

		j.Score(g)
	} else {
		return j.Error("ImportClause", ErrInvalidImport)
	}

	ic.Comments[5] = j.AcceptRunWhitespaceComments()
	ic.Tokens = j.ToTokens()

	return nil
}

// FromClause as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-FromClause
//
// ModuleSpecifier must be non-nil.
type FromClause struct {
	ModuleSpecifier *Token
	Comments        Comments
	Tokens          Tokens
}

func (fc *FromClause) parse(j *jsParser) error {
	if !j.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "from"}) {
		return j.Error("FromClause", ErrMissingFrom)
	}

	fc.Comments = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.Accept(TokenStringLiteral) {
		return j.Error("FromClause", ErrMissingModuleSpecifier)
	}

	fc.ModuleSpecifier = j.GetLastToken()
	fc.Tokens = j.ToTokens()

	return nil
}

// NamedImports as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-NamedImports
type NamedImports struct {
	ImportList []ImportSpecifier
	Comments   [2]Comments
	Tokens     Tokens
}

func (ni *NamedImports) parse(j *jsParser) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		return j.Error("NamedImports", ErrInvalidNamedImport)
	}

	ni.Comments[0] = j.AcceptRunWhitespaceNoNewlineComments()

	for {
		g := j.NewGoal()

		g.AcceptRunWhitespace()

		if g.Accept(TokenRightBracePunctuator) {
			break
		}

		j.AcceptRunWhitespaceNoComment()

		if !j.SkipTypeImport() {
			g := j.NewGoal()
			is := len(ni.ImportList)

			ni.ImportList = append(ni.ImportList, ImportSpecifier{})
			if err := ni.ImportList[is].parse(&g); err != nil {
				return j.Error("NamedImports", err)
			}

			name := ni.ImportList[is].ImportedBinding.Data
			for _, im := range ni.ImportList[:is] {
				if im.ImportedBinding.Data == name {
					return j.Error("NamedImports", ErrInvalidNamedImport)
				}
			}

			j.Score(g)
		}

		g = j.NewGoal()

		g.AcceptRunWhitespace()

		if g.Accept(TokenRightBracePunctuator) {
			break
		} else if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			return g.Error("NamedImports", ErrInvalidNamedImport)
		}

		j.Score(g)
	}

	ni.Comments[1] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()
	j.Next()

	ni.Tokens = j.ToTokens()

	return nil
}

// ImportSpecifier as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ImportSpecifier
//
// ImportedBinding must be non-nil, and IdentifierName should be non-nil.
type ImportSpecifier struct {
	IdentifierName  *Token
	ImportedBinding *Token
	Comments        [4]Comments
	Tokens          Tokens
}

func (is *ImportSpecifier) parse(j *jsParser) error {
	is.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.Accept(TokenIdentifier, TokenKeyword) {
		return j.Error("ImportSpecifier", ErrInvalidImportSpecifier)
	}

	is.ImportedBinding = j.GetLastToken()
	if is.ImportedBinding.Type == TokenIdentifier || is.ImportedBinding.Data == "yield" || is.ImportedBinding.Data == "await" {
		g := j.NewGoal()

		g.AcceptRunWhitespace()

		if g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "as"}) {
			is.Comments[1] = j.AcceptRunWhitespaceComments()

			j.AcceptRunWhitespace()
			j.Next()

			is.Comments[2] = j.AcceptRunWhitespaceComments()

			j.AcceptRunWhitespace()

			g = j.NewGoal()

			is.IdentifierName = is.ImportedBinding
			if is.ImportedBinding = g.parseIdentifier(false, false); is.ImportedBinding == nil {
				return g.Error("ImportSpecifier", ErrNoIdentifier)
			}

			j.Score(g)
		}
	}

	if is.IdentifierName == nil {
		is.IdentifierName = &Token{
			Token:   is.ImportedBinding.Token,
			Pos:     is.ImportedBinding.Pos,
			Line:    is.ImportedBinding.Line,
			LinePos: is.ImportedBinding.LinePos,
		}
	}

	g := j.NewGoal()

	g.AcceptRunWhitespace()

	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
		is.Comments[3] = j.AcceptRunWhitespaceComments()
	} else {
		is.Comments[3] = j.AcceptRunWhitespaceNoNewlineComments()
	}

	is.Tokens = j.ToTokens()

	return nil
}

// ExportDeclaration as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ExportDeclaration
//
// It is only valid for one of ExportClause, ExportFromClause,
// VariableStatement, Declaration, DefaultFunction, DefaultClass, or
// DefaultAssignmentExpression to be non-nil.
//
// FromClause can be non-nil exclusively or paired with ExportClause.
type ExportDeclaration struct {
	ExportClause                *ExportClause
	ExportFromClause            *Token
	FromClause                  *FromClause
	VariableStatement           *VariableStatement
	Declaration                 *Declaration
	DefaultFunction             *FunctionDeclaration
	DefaultClass                *ClassDeclaration
	DefaultAssignmentExpression *AssignmentExpression
	Comments                    [7]Comments
	Tokens                      Tokens
}

func (ed *ExportDeclaration) parse(j *jsParser) error {
	ed.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "export"}) {
		return j.Error("ExportDeclaration", ErrInvalidExportDeclaration)
	}

	ed.Comments[1] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "default"}) {
		ed.Comments[2] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()

		g := j.NewGoal()

		switch g.Peek().Data {
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
		case "abstract":
			h := g.NewGoal()

			if h.SkipAbstract() {
				h.AcceptRunWhitespaceNoNewLine()
			}

			if h.Peek() == (parser.Token{Type: TokenKeyword, Data: "class"}) {
				ed.DefaultClass = new(ClassDeclaration)
				if err := ed.DefaultClass.parse(&h, false, false, true); err != nil {
					return j.Error("ExportDeclaration", err)
				}

				g.Score(h)
				j.Score(g)

				break
			}

			fallthrough
		default:
			ed.DefaultAssignmentExpression = new(AssignmentExpression)
			if err := ed.DefaultAssignmentExpression.parse(&g, true, false, true); err != nil {
				return j.Error("ExportDeclaration", err)
			}

			j.Score(g)

			ed.Comments[5] = j.AcceptRunWhitespaceComments()

			if !j.parseSemicolon() {
				return j.Error("ExportDeclaration", ErrMissingSemiColon)
			}
		}
	} else if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "*"}) {
		ed.Comments[2] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()

		g := j.NewGoal()

		if g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "as"}) {
			ed.Comments[3] = g.AcceptRunWhitespaceComments()

			g.AcceptRunWhitespace()

			if ed.ExportFromClause = g.parseIdentifier(false, false); ed.ExportFromClause == nil {
				return g.Error("ExportDeclaration", ErrNoIdentifier)
			}

			j.Score(g)

			ed.Comments[4] = j.AcceptRunWhitespaceComments()
			j.AcceptRunWhitespace()

			g = j.NewGoal()
		}

		ed.FromClause = new(FromClause)
		if err := ed.FromClause.parse(&g); err != nil {
			return j.Error("ExportDeclaration", err)
		}

		j.Score(g)

		ed.Comments[5] = j.AcceptRunWhitespaceComments()

		if !j.parseSemicolon() {
			return j.Error("ExportDeclaration", ErrMissingSemiColon)
		}
	} else if g := j.NewGoal(); g.Peek() == (parser.Token{Type: TokenKeyword, Data: "var"}) {
		ed.VariableStatement = new(VariableStatement)
		if err := ed.VariableStatement.parse(&g, false, true); err != nil {
			return j.Error("ExportDeclaration", err)
		}

		j.Score(g)
	} else if g.Peek() == (parser.Token{Type: TokenPunctuator, Data: "{"}) {
		ed.ExportClause = new(ExportClause)
		if err := ed.ExportClause.parse(&g); err != nil {
			return j.Error("ExportDeclaration", err)
		}

		j.Score(g)

		g = j.NewGoal()

		g.AcceptRunWhitespace()

		if g.Peek() == (parser.Token{Type: TokenIdentifier, Data: "from"}) {
			j.AcceptRunWhitespaceNoComment()

			ed.Comments[4] = j.AcceptRunWhitespaceComments()

			j.AcceptRunWhitespace()

			g = j.NewGoal()

			g.Skip()
			g.AcceptRunWhitespace()

			if g.Accept(TokenStringLiteral) {
				g = j.NewGoal()
				ed.FromClause = new(FromClause)
				_ = ed.FromClause.parse(&g)

				j.Score(g)
			}
		}

		ed.Comments[5] = j.AcceptRunWhitespaceComments()

		if !j.parseSemicolon() {
			return j.Error("ExportDeclaration", ErrMissingSemiColon)
		}
	} else {
		ed.Declaration = new(Declaration)
		if err := ed.Declaration.parse(&g, false, true); err != nil {
			return j.Error("ExportDeclaration", err)
		}

		j.Score(g)
	}

	ed.Comments[6] = j.AcceptRunWhitespaceNoNewlineComments()
	ed.Tokens = j.ToTokens()

	return nil
}

// ExportClause as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ExportClause
type ExportClause struct {
	ExportList []ExportSpecifier
	Comments   [2]Comments
	Tokens     Tokens
}

func (ec *ExportClause) parse(j *jsParser) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		return j.Error("ExportClause", ErrInvalidExportClause)
	}

	ec.Comments[0] = j.AcceptRunWhitespaceNoNewlineComments()

	for {
		j.AcceptRunWhitespaceNoComment()

		g := j.NewGoal()

		g.AcceptRunWhitespace()

		if g.Accept(TokenRightBracePunctuator) {
			break
		}

		g = j.NewGoal()
		es := len(ec.ExportList)

		ec.ExportList = append(ec.ExportList, ExportSpecifier{})
		if err := ec.ExportList[es].parse(&g); err != nil {
			return j.Error("ExportClause", err)
		}

		j.Score(g)

		g = j.NewGoal()

		g.AcceptRunWhitespace()

		if g.Accept(TokenRightBracePunctuator) {
			break
		} else if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			return g.Error("ExportClause", ErrInvalidExportClause)
		}

		j.Score(g)
	}

	ec.Comments[1] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()
	j.Next()

	ec.Tokens = j.ToTokens()

	return nil
}

// ExportSpecifier as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ExportSpecifier
//
// IdentifierName must be non-nil, EIdentifierName should be non-nil.
type ExportSpecifier struct {
	IdentifierName  *Token
	EIdentifierName *Token
	Comments        [4]Comments
	Tokens          Tokens
}

func (es *ExportSpecifier) parse(j *jsParser) error {
	es.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.Accept(TokenIdentifier, TokenKeyword) {
		return j.Error("ExportSpecifier", ErrNoIdentifier)
	}

	es.IdentifierName = j.GetLastToken()
	g := j.NewGoal()

	g.AcceptRunWhitespace()

	if g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "as"}) {
		es.Comments[1] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()
		j.Next()

		es.Comments[2] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()

		if !j.Accept(TokenIdentifier, TokenKeyword) {
			return j.Error("ExportSpecifier", ErrNoIdentifier)
		}

		es.EIdentifierName = j.GetLastToken()
	} else {
		es.EIdentifierName = &Token{
			Token:   es.IdentifierName.Token,
			Pos:     es.IdentifierName.Pos,
			Line:    es.IdentifierName.Line,
			LinePos: es.IdentifierName.LinePos,
		}
	}

	g = j.NewGoal()

	g.AcceptRunWhitespace()

	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
		es.Comments[3] = j.AcceptRunWhitespaceComments()
	} else {
		es.Comments[3] = j.AcceptRunWhitespaceNoNewlineComments()
	}

	es.Tokens = j.ToTokens()

	return nil
}
