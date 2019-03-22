package javascript

import (
	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

type Module struct {
	Imports    []ImportDeclaration
	Statements []StatementListItem
	Exports    []ExportDeclaration
	Tokens     []TokenPos
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
	for {
		if j.AcceptRunWhitespace() == parser.TokenDone {
			m.Tokens = j.ToTokens()
			return m, nil
		} else {
			var err error
			g := j.NewGoal()
			if g.AcceptToken(parser.Token{TokenKeyword, "import"}) {
				var i ImportDeclaration
				i, err = g.parseImportDeclaration()
				m.Imports = append(m.Imports, i)
			} else if g.AcceptToken(parser.Token{TokenKeyword, "export"}) {
				var e ExportDeclaration
				e, err = g.parseExportDeclaration()
				m.Exports = append(m.Exports, e)
			} else {
				var s StatementListItem
				s, err = g.parseStatementListItem(false, false, false)
				m.Statements = append(m.Statements, s)
			}
			if err != nil {
				return m, j.Error(err)
			}
			j.Score(g)
		}
	}
}

type ImportDeclaration struct {
	*ImportClause
	FromClause
	Tokens []TokenPos
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
	ImportedDefaultBinding *TokenPos
	NamespaceImport        *ImportedBinding
	NamedImports           *NamedImports
	Tokens                 []TokenPos
}

func (j *jsParser) parseImportClause() (ImportClause, error) {
	var ic ImportClause
	g := j.NewGoal()
	if g.Accept(TokenIdentifier) {
		ic.ImportedDefaultBinding = j.GetLastToken()
		j.Score(g)
		j.AcceptRunWhitespace()
		if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return ic, nil
		}
		j.AcceptRunWhitespace()
		g = j.NewGoal()
	}
	if g.AcceptToken(parser.Token{TokenPunctuator, "*"}) {
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenIdentifier, "as"}) {
			return ic, j.Error(ErrInvalidNameSpaceImport)
		}
		g.AcceptRunWhitespace()
		if !g.Accept(TokenIdentifier) {
			ib, err := g.parseImportedBinding()
			if err != nil {
				return ic, j.Error(err)
			}
			ic.NamespaceImport = &ib
		}
	} else if g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		ni, err := g.parseNamedImports()
		if err != nil {
			return ic, j.Error(err)
		}
		ic.NamedImports = &ni
	} else {
		return ic, j.Error(ErrInvalidImport)
	}
	j.Score(g)
	ic.Tokens = j.ToTokens()
	return ic, nil
}

type ImportedBinding BindingIdentifier

func (j *jsParser) parseImportedBinding() (ImportedBinding, error) {
	b, err := j.parseBindingIdentifier(false, false)
	return ImportedBinding(b), err
}

type FromClause struct {
	ModuleSpecifier *TokenPos
	Tokens          []TokenPos
}

func (j *jsParser) parseFromClause() (FromClause, error) {
	var fc FromClause
	j.AcceptRunWhitespace()
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
	Tokens     []TokenPos
}

func (j *jsParser) parseNamedImports() (NamedImports, error) {
	var ni NamedImports
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, "}"}) {
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
		if j.AcceptToken(parser.Token{TokenPunctuator, "}"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return ni, j.Error(ErrInvalidNamedImport)
		}
	}
	ni.Tokens = j.ToTokens()
	return ni, nil
}

type ImportSpecifier struct {
	IdentifierName  *TokenPos
	ImportedBinding ImportedBinding
	Tokens          []TokenPos
}

func (j *jsParser) parseImportSpecifier() (ImportSpecifier, error) {
	var is ImportSpecifier
	g := j.NewGoal()
	if !g.Accept(TokenIdentifier) {
		return is, j.Error(ErrInvalidImportSpecifier)
	}
	g.AcceptRunWhitespace()
	var err error
	if g.AcceptToken(parser.Token{TokenIdentifier, "as"}) { // No IdentifierName
		g.AcceptRunWhitespace()
		if is.ImportedBinding, err = g.parseImportedBinding(); err != nil {
			return is, j.Error(err)
		}
		j.Score(g)
	} else {
		if is.ImportedBinding, err = j.parseImportedBinding(); err != nil {
			return is, j.Error(err)
		}
	}
	is.Tokens = j.ToTokens()
	return is, nil
}

type ExportDeclaration struct {
	ExportClause      *ExportClause
	FromClause        *FromClause
	VariableStatement *VariableStatement
	Declaration       *Declaration
	Default           Token
	Tokens            []TokenPos
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
		case "function":
			fd, err := g.parseFunctionDeclaration(false, false, true)
			if err != nil {
				return ed, j.Error(err)
			}
			j.Score(g)
			ed.Default = &fd
		case "async":
			af, err := g.parseAsyncFunctionDeclaration(false, false, true)
			if err != nil {
				return ed, j.Error(err)
			}
			j.Score(g)
			ed.Default = &af
		case "class":
			cd, err := g.parseClassDeclaration(false, false, true)
			if err != nil {
				return ed, j.Error(err)
			}
			j.Score(g)
			ed.Default = &cd
		default:
			ae, err := g.parseAssignmentExpression(true, false, false)
			if err != nil {
				return ed, j.Error(err)
			}
			j.Score(g)
			ed.Default = &ae
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
	Tokens     []TokenPos
}

func (j *jsParser) parseExportClause() (ExportClause, error) {
	var ec ExportClause
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, "}"}) {
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
		if j.AcceptToken(parser.Token{TokenPunctuator, "}"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return ec, j.Error(ErrInvalidExportClause)
		}
	}
	ec.Tokens = j.ToTokens()
	return ec, nil
}

type ExportSpecifier struct {
	IdentifierName, EIdentifierName *TokenPos
	Tokens                          []TokenPos
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
