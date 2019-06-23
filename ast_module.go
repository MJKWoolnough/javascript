package javascript

import (
	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

type Module struct {
	ModuleListItems []ModuleListItem
	Tokens          Tokens
}

func ParseModule(t parser.Tokeniser) (Module, error) {
	j, err := newJSParser(t)
	if err != nil {
		return Module{}, err
	}
	return j.parseModule()
}

func (j *jsParser) parseModule() (Module, error) {
	var m Module
	for j.AcceptRunWhitespace() != parser.TokenDone {
		g := j.NewGoal()
		ml, err := g.parseModuleStatement()
		if err != nil {
			return m, j.Error(err)
		}
		j.Score(g)
		m.ModuleListItems = append(m.ModuleListItems, ml)
	}
	m.Tokens = j.ToTokens()
	return m, nil
}

type ModuleListItem struct {
	ImportDeclaration *ImportDeclaration
	StatementListItem *StatementListItem
	ExportDeclaration *ExportDeclaration
	Tokens            Tokens
}

func (j *jsParser) parseModuleStatement() (ModuleListItem, error) {
	var ml ModuleListItem
	g := j.NewGoal()
	switch g.Peek() {
	case parser.Token{TokenKeyword, "import"}:
		i, err := g.parseImportDeclaration()
		if err != nil {
			return ml, j.Error(err)
		}
		ml.ImportDeclaration = &i
	case parser.Token{TokenKeyword, "export"}:
		e, err := g.parseExportDeclaration()
		if err != nil {
			return ml, j.Error(err)
		}
		ml.ExportDeclaration = &e
	default:
		s, err := g.parseStatementListItem(false, false, false)
		if err != nil {
			return ml, j.Error(err)
		}
		ml.StatementListItem = &s
	}
	j.Score(g)
	ml.Tokens = j.ToTokens()
	return ml, nil
}

type ImportDeclaration struct {
	*ImportClause
	FromClause
	Tokens Tokens
}

func (j *jsParser) parseImportDeclaration() (ImportDeclaration, error) {
	var id ImportDeclaration
	j.AcceptToken(parser.Token{TokenKeyword, "import"})
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	if g.Accept(TokenStringLiteral) {
		id.FromClause.Tokens = g.ToTokens()
		id.ModuleSpecifier = &id.FromClause.Tokens[0]
	} else {
		ic, err := g.parseImportClause()
		if err != nil {
			return id, j.Error(err)
		}
		id.ImportClause = &ic
		j.Score(g)
		j.AcceptRunWhitespace()
		g = j.NewGoal()
		id.FromClause, err = g.parseFromClause()
		if err != nil {
			return id, j.Error(err)
		}
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
		return id, j.Error(ErrMissingSemiColon)
	}
	id.Tokens = j.ToTokens()
	return id, nil
}

type ImportClause struct {
	ImportedDefaultBinding *ImportedBinding
	NameSpaceImport        *ImportedBinding
	NamedImports           *NamedImports
	Tokens                 Tokens
}

func (j *jsParser) parseImportClause() (ImportClause, error) {
	var ic ImportClause
	if t := j.Peek().Type; t == TokenIdentifier || t == TokenKeyword {
		g := j.NewGoal()
		ib, err := g.parseImportedBinding()
		if err != nil {
			return ic, j.Error(err)
		}
		j.Score(g)
		ic.ImportedDefaultBinding = &ib
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			ic.Tokens = j.ToTokens()
			return ic, nil
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
			return ic, j.Error(ErrInvalidNameSpaceImport)
		}
		g.AcceptRunWhitespace()
		h := g.NewGoal()
		ib, err := h.parseImportedBinding()
		if err != nil {
			return ic, g.Error(err)
		}
		g.Score(h)
		j.Score(g)
		ic.NameSpaceImport = &ib
	} else if j.Peek() == (parser.Token{TokenPunctuator, "{"}) {
		g := j.NewGoal()
		ni, err := g.parseNamedImports()
		if err != nil {
			return ic, j.Error(err)
		}
		j.Score(g)
		ic.NamedImports = &ni
	} else {
		return ic, j.Error(ErrInvalidImport)
	}
	ic.Tokens = j.ToTokens()
	return ic, nil
}

type ImportedBinding BindingIdentifier

func (j *jsParser) parseImportedBinding() (ImportedBinding, error) {
	b, err := j.parseBindingIdentifier(false, false)
	return ImportedBinding(b), err
}

type FromClause struct {
	ModuleSpecifier *Token
	Tokens          Tokens
}

func (j *jsParser) parseFromClause() (FromClause, error) {
	var fc FromClause
	if !j.AcceptToken(parser.Token{TokenIdentifier, "from"}) {
		return fc, j.Error(ErrMissingFrom)
	}
	j.AcceptRunWhitespace()
	if !j.Accept(TokenStringLiteral) {
		return fc, j.Error(ErrMissingModuleSpecifier)
	}
	fc.Tokens = j.ToTokens()
	fc.ModuleSpecifier = &fc.Tokens[len(fc.Tokens)-1]
	return fc, nil
}

type NamedImports struct {
	ImportList []ImportSpecifier
	Tokens     Tokens
}

func (j *jsParser) parseNamedImports() (NamedImports, error) {
	var ni NamedImports
	j.AcceptToken(parser.Token{TokenPunctuator, "{"})
	for {
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		}
		g := j.NewGoal()
		is, err := g.parseImportSpecifier()
		if err != nil {
			return ni, j.Error(err)
		}
		ni.ImportList = append(ni.ImportList, is)
		j.Score(g)
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return ni, j.Error(ErrInvalidNamedImport)
		}
	}
	ni.Tokens = j.ToTokens()
	return ni, nil
}

type ImportSpecifier struct {
	IdentifierName  *Token
	ImportedBinding ImportedBinding
	Tokens          Tokens
}

func (j *jsParser) parseImportSpecifier() (ImportSpecifier, error) {
	var is ImportSpecifier
	if err := j.FindGoal(func(j *jsParser) error {
		if !j.Accept(TokenIdentifier) {
			return errNotApplicable
		}
		in := j.GetLastToken()
		j.AcceptRunWhitespace()
		if !j.AcceptToken(parser.Token{TokenIdentifier, "as"}) {
			return ErrInvalidImportSpecifier
		}
		j.AcceptRunWhitespace()
		ib, err := j.parseImportedBinding()
		if err != nil {
			return err
		}
		is.IdentifierName = in
		is.ImportedBinding = ib
		return nil
	}, func(j *jsParser) error {
		ib, err := j.parseImportedBinding()
		if err != nil {
			return err
		}
		is.ImportedBinding = ib
		return nil
	}); err != nil {
		return is, err
	}
	is.Tokens = j.ToTokens()
	return is, nil
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

func (j *jsParser) parseExportDeclaration() (ExportDeclaration, error) {
	var ed ExportDeclaration
	j.AcceptToken(parser.Token{TokenKeyword, "export"})
	j.AcceptRunWhitespace()
	if j.AcceptToken(parser.Token{TokenKeyword, "default"}) {
		j.AcceptRunWhitespace()
		tk := j.Peek()
		g := j.NewGoal()
		switch tk.Data {
		case "async", "function":
			fd, err := g.parseFunctionDeclaration(false, false, true)
			if err != nil {
				return ed, j.Error(err)
			}
			j.Score(g)
			ed.DefaultFunction = &fd
		case "class":
			cd, err := g.parseClassDeclaration(false, false, true)
			if err != nil {
				return ed, j.Error(err)
			}
			j.Score(g)
			ed.DefaultClass = &cd
		default:
			ae, err := g.parseAssignmentExpression(true, false, false)
			if err != nil {
				return ed, j.Error(err)
			}
			j.Score(g)
			ed.DefaultAssignmentExpression = &ae
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
				return ed, j.Error(ErrMissingSemiColon)
			}
		}
	} else if j.AcceptToken(parser.Token{TokenPunctuator, "*"}) {
		j.AcceptRunWhitespace()
		g := j.NewGoal()
		fc, err := g.parseFromClause()
		if err != nil {
			return ed, j.Error(err)
		}
		j.Score(g)
		ed.FromClause = &fc
		if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			return ed, j.Error(ErrMissingSemiColon)
		}
	} else if g := j.NewGoal(); g.AcceptToken(parser.Token{TokenKeyword, "var"}) {
		v, err := j.parseVariableStatement(false, false)
		if err != nil {
			return ed, j.Error(err)
		}
		ed.VariableStatement = &v
	} else if g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		ec, err := g.parseExportClause()
		if err != nil {
			return ed, j.Error(err)
		}
		j.Score(g)
		ed.ExportClause = &ec
		j.AcceptRunWhitespace()
		if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			g = j.NewGoal()
			fc, err := g.parseFromClause()
			if err != nil {
				return ed, j.Error(err)
			}
			ed.FromClause = &fc
			j.Score(g)
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
				return ed, j.Error(ErrMissingSemiColon)
			}
		}
	} else {
		d, err := j.parseDeclaration(false, false)
		if err != nil {
			return ed, j.Error(err)
		}
		ed.Declaration = &d
	}
	ed.Tokens = j.ToTokens()
	return ed, nil
}

type ExportClause struct {
	ExportList []ExportSpecifier
	Tokens     Tokens
}

func (j *jsParser) parseExportClause() (ExportClause, error) {
	var ec ExportClause
	for {
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		}
		g := j.NewGoal()
		es, err := g.parseExportSpecifier()
		if err != nil {
			return ec, j.Error(err)
		}
		j.Score(g)
		ec.ExportList = append(ec.ExportList, es)
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return ec, j.Error(ErrInvalidExportClause)
		}
	}
	ec.Tokens = j.ToTokens()
	return ec, nil
}

type ExportSpecifier struct {
	IdentifierName, EIdentifierName *Token
	Tokens                          Tokens
}

func (j *jsParser) parseExportSpecifier() (ExportSpecifier, error) {
	var es ExportSpecifier
	if !j.Accept(TokenIdentifier) {
		return es, j.Error(ErrMissingIdentifier)
	}
	es.IdentifierName = j.GetLastToken()
	g := j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenIdentifier, "as"}) {
		g.AcceptRunWhitespace()
		if !g.Accept(TokenIdentifier) {
			return es, j.Error(ErrMissingIdentifier)
		}
		j.Score(g)
		es.EIdentifierName = j.GetLastToken()
	}
	es.Tokens = j.ToTokens()
	return es, nil
}

const (
	ErrInvalidImport          errors.Error = "invalid import statement"
	ErrInvalidNameSpaceImport errors.Error = "invalid namespace import"
	ErrMissingFrom            errors.Error = "missing from"
	ErrMissingModuleSpecifier errors.Error = "missing module specifier"
	ErrInvalidNamedImport     errors.Error = "invalid named import list"
	ErrInvalidImportSpecifier errors.Error = "invalid import specifier"
	ErrMissingIdentifier      errors.Error = "missing identifier"
	ErrInvalidExportClause    errors.Error = "invalid export clause"
)
