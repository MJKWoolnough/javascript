package javascript

import "vimagination.zapto.org/parser"

// Block as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-Block
type Block struct {
	StatementList []StatementListItem
	Comments      [2]Comments
	Tokens        Tokens
}

func (b *Block) parse(j *jsParser, yield, await, ret bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		return j.Error("Block", ErrMissingOpeningBrace)
	}

	b.Comments[0] = j.AcceptRunWhitespaceNoNewlineComments()

	for {
		j.AcceptRunWhitespaceNoComment()

		g := j.NewGoal()

		g.AcceptRunWhitespace()

		if g.Accept(TokenRightBracePunctuator) {
			break
		}

		g = j.NewGoal()
		si := len(b.StatementList)

		b.StatementList = append(b.StatementList, StatementListItem{})
		if err := b.StatementList[si].parse(&g, yield, await, ret); err != nil {
			return j.Error("Block", err)
		}

		j.Score(g)
	}

	b.Comments[1] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()
	j.Skip()

	b.Tokens = j.ToTokens()

	return nil
}

// StatementListItem as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-StatementListItem
// Only one of Statement, or Declaration must be non-nil.
type StatementListItem struct {
	Statement   *Statement
	Declaration *Declaration
	Comments    [2]Comments
	Tokens      Tokens
}

func (si *StatementListItem) parse(j *jsParser, yield, await, ret bool) error {
	si.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if j.SkipType() || j.SkipInterface() || j.SkipDeclare() {
		si.Statement = &Statement{Tokens: j.ToTokens()}
		si.Tokens = j.ToTokens()

		return nil
	}

	g := j.NewGoal()

	var declaration bool

	switch t := g.Peek(); t {
	case parser.Token{Type: TokenIdentifier, Data: "let"}, parser.Token{Type: TokenKeyword, Data: "const"}:
		declaration = true
	case parser.Token{Type: TokenIdentifier, Data: "abstract"}:
		g.Skip()
		g.AcceptRunWhitespaceNoNewLine()

		if g.Peek() != (parser.Token{Type: TokenKeyword, Data: "class"}) {
			break
		}

		fallthrough
	case parser.Token{Type: TokenKeyword, Data: "class"}:
		g.Skip()
		g.AcceptRunWhitespace()

		if g.parseIdentifier(yield, await) != nil {
			declaration = true
		}
	case parser.Token{Type: TokenIdentifier, Data: "async"}:
		g.Skip()
		g.AcceptRunWhitespaceNoNewLine()

		if g.Peek() != (parser.Token{Type: TokenKeyword, Data: "function"}) {
			break
		}

		fallthrough
	case parser.Token{Type: TokenKeyword, Data: "function"}:
		g.Skip()
		g.AcceptRunWhitespace()

		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "*"}) {
			g.AcceptRunWhitespace()
		}

		if g.parseIdentifier(yield, await) != nil {
			declaration = true
		}
	}

	g = j.NewGoal()

	if declaration {
		si.Declaration = new(Declaration)
		if err := si.Declaration.parse(&g, yield, await); err != nil {
			return j.Error("StatementListItem", err)
		}
	} else {
		si.Statement = new(Statement)
		if err := si.Statement.parse(&g, yield, await, ret); err != nil {
			return j.Error("StatementListItem", err)
		}
	}

	j.Score(g)

	si.Comments[1] = j.AcceptRunWhitespaceNoNewlineComments()
	si.Tokens = j.ToTokens()

	return nil
}

// StatementType determines the type of a Statement type
type StatementType uint8

// Valid StatementType's
const (
	StatementNormal StatementType = iota
	StatementContinue
	StatementBreak
	StatementReturn
	StatementThrow
	StatementDebugger
)

// Statement as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-Statement
//
// It is only valid for one of the pointer type to be non-nil.
//
// If LabelIdentifier is non-nil, either one of LabelledItemFunction, or
// LabelledItemStatement must be non-nil, or Type must be StatementContinue or
// StatementBreak.
//
// If Type is StatementThrow, ExpressionStatement must be non-nil.
type Statement struct {
	Type                    StatementType
	BlockStatement          *Block
	VariableStatement       *VariableStatement
	ExpressionStatement     *Expression
	IfStatement             *IfStatement
	IterationStatementDo    *IterationStatementDo
	IterationStatementWhile *IterationStatementWhile
	IterationStatementFor   *IterationStatementFor
	SwitchStatement         *SwitchStatement
	WithStatement           *WithStatement
	LabelIdentifier         *Token
	LabelledItemFunction    *FunctionDeclaration
	LabelledItemStatement   *Statement
	TryStatement            *TryStatement
	Comments                [2]Comments
	Tokens                  Tokens
}

func (s *Statement) parse(j *jsParser, yield, await, ret bool) error {
	g := j.NewGoal()

	switch g.Peek() {
	case parser.Token{Type: TokenPunctuator, Data: "{"}:
		s.BlockStatement = new(Block)
		if err := s.BlockStatement.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{Type: TokenKeyword, Data: "var"}:
		s.VariableStatement = new(VariableStatement)
		if err := s.VariableStatement.parse(&g, yield, await); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{Type: TokenPunctuator, Data: ";"}:
		g.Skip()
	case parser.Token{Type: TokenKeyword, Data: "if"}:
		s.IfStatement = new(IfStatement)
		if err := s.IfStatement.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{Type: TokenKeyword, Data: "do"}:
		s.IterationStatementDo = new(IterationStatementDo)
		if err := s.IterationStatementDo.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{Type: TokenKeyword, Data: "while"}:
		s.IterationStatementWhile = new(IterationStatementWhile)
		if err := s.IterationStatementWhile.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{Type: TokenKeyword, Data: "for"}:
		s.IterationStatementFor = new(IterationStatementFor)
		if err := s.IterationStatementFor.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{Type: TokenKeyword, Data: "switch"}:
		s.SwitchStatement = new(SwitchStatement)
		if err := s.SwitchStatement.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{Type: TokenKeyword, Data: "continue"}, parser.Token{Type: TokenKeyword, Data: "break"}:
		if g.Peek().Data == "continue" {
			s.Type = StatementContinue
		} else {
			s.Type = StatementBreak
		}

		g.Skip()

		s.Comments[0] = g.AcceptRunWhitespaceComments()

		h := g.NewGoal()

		if !h.parseSemicolon() {
			g.AcceptRunWhitespaceNoNewLine()

			if s.LabelIdentifier = g.parseIdentifier(yield, await); s.LabelIdentifier == nil {
				return g.Error("Statement", ErrNoIdentifier)
			}

			s.Comments[1] = g.AcceptRunWhitespaceComments()

			if !g.parseSemicolon() {
				return g.Error("Statement", ErrMissingSemiColon)
			}
		} else {
			g.Score(h)
		}
	case parser.Token{Type: TokenKeyword, Data: "return"}:
		if !ret {
			return g.Error("Statement", ErrInvalidStatement)
		}

		g.Skip()

		s.Type = StatementReturn

		h := g.NewGoal()

		if !h.parseSemicolon() {
			g.AcceptRunWhitespaceNoComment()

			h := g.NewGoal()
			s.ExpressionStatement = new(Expression)

			if err := s.ExpressionStatement.parse(&h, true, yield, await); err != nil {
				return g.Error("Statement", err)
			}

			g.Score(h)

			if !g.parseSemicolon() {
				return g.Error("Statement", ErrMissingSemiColon)
			}
		} else {
			s.Comments[0] = g.AcceptRunWhitespaceComments()
			g.parseSemicolon()
		}
	case parser.Token{Type: TokenKeyword, Data: "with"}:
		s.WithStatement = new(WithStatement)
		if err := s.WithStatement.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{Type: TokenKeyword, Data: "throw"}:
		g.Skip()

		s.Type = StatementThrow
		h := g.NewGoal()

		if h.AcceptRunWhitespaceNoNewLine() == TokenLineTerminator {
			return h.Error("Statement", ErrUnexpectedLineTerminator)
		}

		g.AcceptRunWhitespaceNoNewLineNoComment()

		h = g.NewGoal()
		s.ExpressionStatement = new(Expression)

		if err := s.ExpressionStatement.parse(&h, true, yield, await); err != nil {
			return g.Error("Statement", err)
		}

		g.Score(h)

		if !g.parseSemicolon() {
			return g.Error("Statement", ErrMissingSemiColon)
		}
	case parser.Token{Type: TokenKeyword, Data: "try"}:
		s.TryStatement = new(TryStatement)
		if err := s.TryStatement.parse(&g, yield, await, ret); err != nil {
			return j.Error("Statement", err)
		}
	case parser.Token{Type: TokenKeyword, Data: "debugger"}:
		g.Skip()

		s.Type = StatementDebugger

		s.Comments[0] = g.AcceptRunWhitespaceComments()

		if !g.parseSemicolon() {
			return g.Error("Statement", ErrMissingSemiColon)
		}
	default:
		if i := g.parseIdentifier(yield, await); i != nil {
			h := g.NewGoal()

			h.AcceptRunWhitespace()

			if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
				s.Comments[0] = g.AcceptRunWhitespaceComments()

				g.AcceptRunWhitespace()
				g.Skip()

				s.Comments[1] = g.AcceptRunWhitespaceComments()

				g.AcceptRunWhitespace()

				s.LabelIdentifier = i
				h := g.NewGoal()

				if h.Peek() == (parser.Token{Type: TokenKeyword, Data: "function"}) {
					s.LabelledItemFunction = new(FunctionDeclaration)
					if err := s.LabelledItemFunction.parse(&h, yield, await, false); err != nil {
						return g.Error("Statement", err)
					}
				} else {
					s.LabelledItemStatement = new(Statement)
					if err := s.LabelledItemStatement.parse(&h, yield, await, ret); err != nil {
						return g.Error("Statement", err)
					}
				}

				g.Score(h)
			}
		}

		if s.LabelIdentifier == nil {
			g = j.NewGoal()

			switch g.Peek() {
			case parser.Token{Type: TokenKeyword, Data: "function"}, parser.Token{Type: TokenKeyword, Data: "class"}:
				return j.Error("Statement", ErrInvalidStatement)
			case parser.Token{Type: TokenIdentifier, Data: "async"}:
				h := g.NewGoal()

				h.Skip()
				h.AcceptRunWhitespaceNoNewLine()

				if h.AcceptToken(parser.Token{Type: TokenKeyword, Data: "function"}) {
					return j.Error("Statement", ErrInvalidStatement)
				}
			}

			s.ExpressionStatement = new(Expression)
			if err := s.ExpressionStatement.parse(&g, true, yield, await); err != nil {
				return j.Error("Statement", err)
			}

			if !g.parseSemicolon() {
				return g.Error("Statement", ErrMissingSemiColon)
			}
		}
	}

	j.Score(g)

	s.Tokens = j.ToTokens()

	return nil
}

// IfStatement as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-IfStatement
type IfStatement struct {
	Expression    Expression
	Statement     Statement
	ElseStatement *Statement
	Comments      [6]Comments
	Tokens        Tokens
}

func (is *IfStatement) parse(j *jsParser, yield, await, ret bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "if"}) {
		return j.Error("IfStatement", ErrInvalidIfStatement)
	}

	is.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		return j.Error("IfStatement", ErrMissingOpeningParenthesis)
	}

	is.Comments[1] = j.AcceptRunWhitespaceNoNewlineComments()

	j.AcceptRunWhitespaceNoComment()

	g := j.NewGoal()

	if err := is.Expression.parse(&g, true, yield, await); err != nil {
		return j.Error("IfStatement", err)
	}

	j.Score(g)

	is.Comments[2] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
		return j.Error("IfStatement", ErrMissingClosingParenthesis)
	}

	is.Comments[3] = j.AcceptRunWhitespaceComments()
	j.AcceptRunWhitespace()

	g = j.NewGoal()

	if err := is.Statement.parse(&g, yield, await, ret); err != nil {
		return j.Error("IfStatement", err)
	}

	if is.Statement.LabelledItemFunction != nil {
		return j.Error("IfStatement", ErrLabelledFunction)
	}

	j.Score(g)

	g = j.NewGoal()

	g.AcceptRunWhitespace()

	if g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "else"}) {
		is.Comments[4] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()
		j.Skip()

		is.Comments[5] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()

		g := j.NewGoal()

		is.ElseStatement = new(Statement)
		if err := is.ElseStatement.parse(&g, yield, await, ret); err != nil {
			return j.Error("IfStatement", err)
		}

		if is.ElseStatement.LabelledItemFunction != nil {
			return j.Error("IfStatement", ErrLabelledFunction)
		}

		j.Score(g)
	}

	is.Tokens = j.ToTokens()

	return nil
}

// IterationStatementDo is the do-while part of IterationStatement as defined
// in ECMA-262
// https://262.ecma-international.org/11.0/#prod-IterationStatement
type IterationStatementDo struct {
	Statement  Statement
	Expression Expression
	Tokens     Tokens
}

func (is *IterationStatementDo) parse(j *jsParser, yield, await, ret bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "do"}) {
		return j.Error("IterationStatementDo", ErrInvalidIterationStatementDo)
	}

	j.AcceptRunWhitespace()

	g := j.NewGoal()

	if err := is.Statement.parse(&g, yield, await, ret); err != nil {
		return j.Error("IterationStatementDo", err)
	}

	if is.Statement.LabelledItemFunction != nil {
		return j.Error("IterationStatementDo", ErrLabelledFunction)
	}

	j.Score(g)

	g = j.NewGoal()

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "while"}) {
		return j.Error("IterationStatementDo", ErrInvalidIterationStatementDo)
	}

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		return j.Error("IterationStatementDo", ErrMissingOpeningParenthesis)
	}

	j.AcceptRunWhitespace()

	g = j.NewGoal()

	if err := is.Expression.parse(&g, true, yield, await); err != nil {
		return j.Error("IterationStatementDo", err)
	}

	j.Score(g)
	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
		return j.Error("IterationStatementDo", ErrMissingClosingParenthesis)
	}

	if !j.parseSemicolon() {
		return j.Error("IterationStatementDo", ErrMissingSemiColon)
	}

	is.Tokens = j.ToTokens()

	return nil
}

// IterationStatementWhile is the while part of IterationStatement as defined
// in ECMA-262
// https://262.ecma-international.org/11.0/#prod-IterationStatement
type IterationStatementWhile struct {
	Expression Expression
	Statement  Statement
	Comments   [4]Comments
	Tokens     Tokens
}

func (is *IterationStatementWhile) parse(j *jsParser, yield, await, ret bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "while"}) {
		return j.Error("IterationStatementWhile", ErrInvalidIterationStatementWhile)
	}

	is.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		return j.Error("IterationStatementWhile", ErrMissingOpeningParenthesis)
	}

	is.Comments[1] = j.AcceptRunWhitespaceNoNewlineComments()

	j.AcceptRunWhitespaceNoComment()

	g := j.NewGoal()

	if err := is.Expression.parse(&g, true, yield, await); err != nil {
		return j.Error("IterationStatementWhile", err)
	}

	j.Score(g)

	is.Comments[2] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
		return j.Error("IterationStatementWhile", ErrMissingClosingParenthesis)
	}

	is.Comments[3] = j.AcceptRunWhitespaceComments()
	j.AcceptRunWhitespace()

	g = j.NewGoal()

	if err := is.Statement.parse(&g, yield, await, ret); err != nil {
		return j.Error("IterationStatementWhile", err)
	}

	if is.Statement.LabelledItemFunction != nil {
		return j.Error("IterationStatementWhile", ErrLabelledFunction)
	}

	j.Score(g)

	is.Tokens = j.ToTokens()

	return nil
}

// ForType determines which kind of for-loop is described by IterationStatementFor
type ForType uint8

// Valid ForType's
const (
	ForNormal ForType = iota
	ForNormalVar
	ForNormalLexicalDeclaration
	ForNormalExpression
	ForInLeftHandSide
	ForInVar
	ForInLet
	ForInConst
	ForOfLeftHandSide
	ForOfVar
	ForOfLet
	ForOfConst
	ForAwaitOfLeftHandSide
	ForAwaitOfVar
	ForAwaitOfLet
	ForAwaitOfConst
)

// IterationStatementFor is the for part of IterationStatement as defined
// in ECMA-262
// https://262.ecma-international.org/11.0/#prod-IterationStatement
//
// Includes TC39 proposal for for-await-of
// https://github.com/tc39/proposal-async-iteration#the-async-iteration-statement-for-await-of
//
// The Type determines which fields must be non-nil:
//
//	ForInLeftHandSide: LeftHandSideExpression and In
//	ForInVar, ForInLet, ForInConst: ForBindingIdentifier, ForBindingPatternObject, or ForBindingPatternArray and In
//	ForOfLeftHandSide, ForAwaitOfLeftHandSide: LeftHandSideExpression and Of
//	ForOfVar, ForAwaitOfVar, ForOfLet, ForAwaitOfLet, ForOfConst, ForAwaitOfConst: ForBindingIdentifier, ForBindingPatternObject, or ForBindingPatternArray and Of
type IterationStatementFor struct {
	Type ForType

	InitExpression *Expression
	InitVar        []VariableDeclaration
	InitLexical    *LexicalDeclaration
	Conditional    *Expression
	Afterthought   *Expression

	LeftHandSideExpression  *LeftHandSideExpression
	ForBindingIdentifier    *Token
	ForBindingPatternObject *ObjectBindingPattern
	ForBindingPatternArray  *ArrayBindingPattern
	In                      *Expression
	Of                      *AssignmentExpression

	Statement Statement
	Tokens    Tokens
}

func (is *IterationStatementFor) parse(j *jsParser, yield, await, ret bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "for"}) {
		return j.Error("IterationStatementFor", ErrInvalidIterationStatementFor)
	}

	j.AcceptRunWhitespace()

	var forAwait bool

	if await {
		forAwait = j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "await"})

		j.AcceptRunWhitespace()
	}

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		return j.Error("IterationStatementFor", ErrMissingOpeningParenthesis)
	}

	j.AcceptRunWhitespace()

	is.Type = ForNormal

	switch j.Peek() {
	case parser.Token{Type: TokenPunctuator, Data: ";"}:
	case parser.Token{Type: TokenKeyword, Data: "const"}, parser.Token{Type: TokenIdentifier, Data: "let"}:
		is.Type = 1

		fallthrough
	case parser.Token{Type: TokenKeyword, Data: "var"}:
		is.Type++

		g := j.NewGoal()

		g.Skip()
		g.AcceptRunWhitespace()

		var (
			opener = "{"
			closer = "}"
		)

		switch g.Peek() {
		case parser.Token{Type: TokenPunctuator, Data: "["}:
			opener = "["
			closer = "]"

			fallthrough
		case parser.Token{Type: TokenPunctuator, Data: "{"}:
			var level uint

		Loop:
			for {
				g.ExceptRun(TokenPunctuator, TokenRightBracePunctuator)

				switch g.Peek().Data {
				case opener:
					level++
				case closer:
					if level--; level == 0 {
						g.Skip()

						break Loop
					}
				}

				g.Skip()
			}
		default:
			g.Skip()
		}

		g.AcceptRunWhitespace()

		switch g.Peek() {
		case parser.Token{Type: TokenKeyword, Data: "in"}:
			is.Type += 4
		case parser.Token{Type: TokenIdentifier, Data: "of"}:
			is.Type += 8
		}

		if is.Type > 4 && j.Peek() == (parser.Token{Type: TokenKeyword, Data: "const"}) {
			is.Type++
		}
	default:
		if forAwait {
			is.Type = ForOfLeftHandSide
		} else {
			is.Type = ForNormalExpression
		}
	}

	if forAwait && is.Type < ForOfLeftHandSide {
		return j.Error("IterationStatementFor", ErrInvalidForAwaitLoop)
	}

	switch is.Type {
	case ForNormal:
		j.Skip()
		j.AcceptRunWhitespace()
	case ForNormalVar:
		j.Skip()

		for {
			j.AcceptRunWhitespace()

			g := j.NewGoal()
			vd := len(is.InitVar)

			is.InitVar = append(is.InitVar, VariableDeclaration{})
			if err := is.InitVar[vd].parse(&g, false, yield, await); err != nil {
				return j.Error("IterationStatementFor", err)
			}

			j.Score(g)
			j.AcceptRunWhitespace()

			if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
				break
			}
		}

		if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) {
			return j.Error("IterationStatementFor", ErrMissingSemiColon)
		}

		j.AcceptRunWhitespace()
	case ForNormalLexicalDeclaration:
		g := j.NewGoal()

		is.InitLexical = new(LexicalDeclaration)
		if err := is.InitLexical.parse(&g, false, yield, await); err != nil {
			return j.Error("IterationStatementFor", err)
		}

		j.Score(g)

		if is.InitLexical.Tokens[len(is.InitLexical.Tokens)-1].Data != ";" {
			j.AcceptRunWhitespace()

			if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) {
				return j.Error("IterationStatementFor", ErrMissingSemiColon)
			}
		}

		j.AcceptRunWhitespace()
	case ForNormalExpression:
		g := j.NewGoal()

		is.InitExpression = new(Expression)
		if err := is.InitExpression.parse(&g, false, yield, await); err != nil {
			return j.Error("IterationStatementFor", err)
		}

		j.Score(g)
		j.AcceptRunWhitespace()

		if len(is.InitExpression.Expressions) == 1 && is.InitExpression.Expressions[0].ConditionalExpression != nil && is.InitExpression.Expressions[0].ConditionalExpression.LogicalORExpression != nil {
			if lhs := is.InitExpression.Expressions[0].ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression; lhs != nil && len(is.InitExpression.Tokens) == len(lhs.Tokens) {
				switch j.Peek() {
				case parser.Token{Type: TokenKeyword, Data: "in"}:
					is.Type = ForInLeftHandSide
					is.InitExpression = nil
					is.LeftHandSideExpression = lhs
				case parser.Token{Type: TokenIdentifier, Data: "of"}:
					is.Type = ForOfLeftHandSide
					is.InitExpression = nil
					is.LeftHandSideExpression = lhs
				}
			}
		}

		if is.InitExpression != nil {
			if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) {
				return j.Error("IterationStatementFor", ErrMissingSemiColon)
			}

			j.AcceptRunWhitespace()
		}
	case ForInVar, ForInLet, ForInConst, ForOfVar, ForOfLet, ForOfConst:
		j.Skip()
		j.AcceptRunWhitespace()

		switch j.Peek() {
		case parser.Token{Type: TokenPunctuator, Data: "{"}:
			g := j.NewGoal()

			is.ForBindingPatternObject = new(ObjectBindingPattern)
			if err := is.ForBindingPatternObject.parse(&g, yield, await); err != nil {
				return j.Error("IterationStatementFor", err)
			}

			j.Score(g)
		case parser.Token{Type: TokenPunctuator, Data: "["}:
			g := j.NewGoal()

			is.ForBindingPatternArray = new(ArrayBindingPattern)
			if err := is.ForBindingPatternArray.parse(&g, yield, await); err != nil {
				return j.Error("IterationStatementFor", err)
			}

			j.Score(g)
		default:
			if is.ForBindingIdentifier = j.parseIdentifier(yield, await); is.ForBindingIdentifier == nil {
				return j.Error("IterationStatementFor", ErrNoIdentifier)
			}
		}

		j.AcceptRunWhitespace()
	case ForOfLeftHandSide:
		g := j.NewGoal()

		is.LeftHandSideExpression = new(LeftHandSideExpression)
		if err := is.LeftHandSideExpression.parse(&g, yield, await); err != nil {
			return j.Error("IterationStatementFor", err)
		}

		j.Score(g)
		j.AcceptRunWhitespace()
	}
	switch is.Type {
	case ForNormal, ForNormalVar, ForNormalLexicalDeclaration, ForNormalExpression:
		if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) {
			g := j.NewGoal()

			is.Conditional = new(Expression)
			if err := is.Conditional.parse(&g, true, yield, await); err != nil {
				return j.Error("IterationStatementFor", err)
			}

			j.Score(g)
			j.AcceptRunWhitespace()

			if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) {
				return j.Error("IterationStatementFor", ErrMissingSemiColon)
			}
		}

		j.AcceptRunWhitespace()

		if j.Peek() != (parser.Token{Type: TokenPunctuator, Data: ")"}) {
			g := j.NewGoal()

			is.Afterthought = new(Expression)
			if err := is.Afterthought.parse(&g, true, yield, await); err != nil {
				return j.Error("IterationStatementFor", err)
			}

			j.Score(g)
		}
	case ForInLeftHandSide, ForInVar, ForInLet, ForInConst:
		j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "in"})
		j.AcceptRunWhitespace()

		g := j.NewGoal()

		is.In = new(Expression)
		if err := is.In.parse(&g, true, yield, await); err != nil {
			return j.Error("IterationStatementFor", err)
		}

		j.Score(g)
	case ForOfLeftHandSide, ForOfVar, ForOfLet, ForOfConst:
		if !j.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "of"}) {
			return j.Error("IterationStatementFor", ErrInvalidForAwaitLoop)
		}

		j.AcceptRunWhitespaceNoComment()

		g := j.NewGoal()

		is.Of = new(AssignmentExpression)
		if err := is.Of.parse(&g, true, yield, await); err != nil {
			return j.Error("IterationStatementFor", err)
		}

		j.Score(g)
	}

	if forAwait {
		is.Type += 4
	}

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
		return j.Error("IterationStatementFor", ErrMissingClosingParenthesis)
	}

	j.AcceptRunWhitespace()

	g := j.NewGoal()

	if err := is.Statement.parse(&g, yield, await, ret); err != nil {
		return j.Error("IterationStatementFor", err)
	}

	if is.Statement.LabelledItemFunction != nil {
		return j.Error("IterationStatementFor", ErrLabelledFunction)
	}

	j.Score(g)

	is.Tokens = j.ToTokens()

	return nil
}

// SwitchStatement as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-SwitchStatement
type SwitchStatement struct {
	Expression             Expression
	CaseClauses            []CaseClause
	DefaultClause          []StatementListItem
	PostDefaultCaseClauses []CaseClause
	Comments               [9]Comments
	Tokens                 Tokens
}

func (ss *SwitchStatement) parse(j *jsParser, yield, await, ret bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "switch"}) {
		return j.Error("SwitchStatement", ErrInvalidSwitchStatement)
	}

	ss.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		return j.Error("SwitchStatement", ErrMissingOpeningParenthesis)
	}

	ss.Comments[1] = j.AcceptRunWhitespaceNoNewlineComments()

	j.AcceptRunWhitespaceNoComment()

	g := j.NewGoal()

	if err := ss.Expression.parse(&g, true, yield, await); err != nil {
		return j.Error("SwitchStatement", err)
	}

	j.Score(g)

	ss.Comments[2] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
		return j.Error("SwitchStatement", ErrMissingClosingParenthesis)
	}

	ss.Comments[3] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		return j.Error("SwitchStatement", ErrMissingOpeningBrace)
	}

	ss.Comments[4] = j.AcceptRunWhitespaceNoNewlineComments()

	for {
		j.AcceptRunWhitespaceNoComment()

		g := j.NewGoal()

		g.AcceptRunWhitespace()

		if g.Accept(TokenRightBracePunctuator) {
			break
		} else if g.Peek() == (parser.Token{Type: TokenKeyword, Data: "default"}) {
			if ss.DefaultClause != nil {
				return j.Error("SwitchStatement", ErrDuplicateDefaultClause)
			}

			ss.Comments[5] = j.AcceptRunWhitespaceComments()

			j.AcceptRunWhitespace()
			j.Skip()

			ss.Comments[6] = j.AcceptRunWhitespaceComments()
			j.AcceptRunWhitespace()

			if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
				return j.Error("SwitchStatement", ErrMissingColon)
			}

			ss.Comments[7] = j.AcceptRunWhitespaceNoNewlineComments()
			ss.DefaultClause = []StatementListItem{}

			for {
				g = j.NewGoal()

				g.AcceptRunWhitespace()

				if pt := g.Peek(); pt == (parser.Token{Type: TokenKeyword, Data: "case"}) || pt.Type == TokenRightBracePunctuator {
					break
				}

				j.AcceptRunWhitespaceNoComment()

				g = j.NewGoal()
				sl := len(ss.DefaultClause)

				ss.DefaultClause = append(ss.DefaultClause, StatementListItem{})
				if err := ss.DefaultClause[sl].parse(&g, yield, await, ret); err != nil {
					return j.Error("SwitchStatement", err)
				}

				j.Score(g)
			}
		} else {
			g = j.NewGoal()

			var cc *CaseClause

			if ss.DefaultClause == nil {
				ss.CaseClauses = append(ss.CaseClauses, CaseClause{})
				cc = &ss.CaseClauses[len(ss.CaseClauses)-1]
			} else {
				ss.PostDefaultCaseClauses = append(ss.PostDefaultCaseClauses, CaseClause{})
				cc = &ss.PostDefaultCaseClauses[len(ss.PostDefaultCaseClauses)-1]
			}

			if err := cc.parse(&g, yield, await, ret); err != nil {
				return j.Error("SwitchStatement", err)
			}

			j.Score(g)
		}
	}

	ss.Comments[8] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()
	j.Skip()

	ss.Tokens = j.ToTokens()

	return nil
}

// CaseClause as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-CaseClauses
type CaseClause struct {
	Expression    Expression
	StatementList []StatementListItem
	Comments      [2]Comments
	Tokens        Tokens
}

func (cc *CaseClause) parse(j *jsParser, yield, await, ret bool) error {
	cc.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "case"}) {
		return j.Error("CaseClause", ErrMissingCaseClause)
	}

	j.AcceptRunWhitespaceNoComment()

	g := j.NewGoal()

	if err := cc.Expression.parse(&g, true, yield, await); err != nil {
		return j.Error("CaseClause", err)
	}

	j.Score(g)
	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
		return j.Error("CaseClause", ErrMissingColon)
	}

	cc.Comments[1] = j.AcceptRunWhitespaceNoNewlineComments()

	for {
		g := j.NewGoal()
		g.AcceptRunWhitespace()

		if tk := g.Peek(); tk == (parser.Token{Type: TokenKeyword, Data: "case"}) || tk == (parser.Token{Type: TokenKeyword, Data: "default"}) || tk.Type == TokenRightBracePunctuator || tk.Type == parser.TokenDone {
			break
		}

		j.AcceptRunWhitespaceNoComment()

		g = j.NewGoal()
		sl := len(cc.StatementList)

		cc.StatementList = append(cc.StatementList, StatementListItem{})
		if err := cc.StatementList[sl].parse(&g, yield, await, ret); err != nil {
			return g.Error("CaseClause", err)
		}

		j.Score(g)
	}

	cc.Tokens = j.ToTokens()

	return nil
}

// WithStatement as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-WithStatement
type WithStatement struct {
	Expression Expression
	Statement  Statement
	Comments   [4]Comments
	Tokens     Tokens
}

func (ws *WithStatement) parse(j *jsParser, yield, await, ret bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "with"}) {
		return j.Error("WithStatement", ErrInvalidWithStatement)
	}

	ws.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		return j.Error("WithStatement", ErrMissingOpeningParenthesis)
	}

	ws.Comments[1] = j.AcceptRunWhitespaceNoNewlineComments()

	j.AcceptRunWhitespaceNoComment()

	g := j.NewGoal()

	if err := ws.Expression.parse(&g, true, yield, await); err != nil {
		return j.Error("WithStatement", err)
	}

	j.Score(g)

	ws.Comments[2] = j.AcceptRunWhitespaceComments()
	j.AcceptRunWhitespace()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
		return j.Error("WithStatement", ErrMissingClosingParenthesis)
	}

	ws.Comments[3] = j.AcceptRunWhitespaceComments()
	j.AcceptRunWhitespace()

	g = j.NewGoal()

	if err := ws.Statement.parse(&g, yield, await, ret); err != nil {
		return j.Error("WithStatement", err)
	}

	if ws.Statement.LabelledItemFunction != nil {
		return j.Error("WithStatement", ErrLabelledFunction)
	}

	j.Score(g)

	ws.Tokens = j.ToTokens()

	return nil
}

// TryStatement as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-TryStatement
//
// Only one of CatchParameterBindingIdentifier,
// CatchParameterObjectBindingPattern, and CatchParameterArrayBindingPattern can
// be non-nil, and must be so if CatchBlock is non-nil.
//
// If one of CatchParameterBindingIdentifier,
// CatchParameterObjectBindingPattern, CatchParameterArrayBindingPattern is
// non-nil, then CatchBlock must be non-nil.
type TryStatement struct {
	TryBlock                           Block
	CatchParameterBindingIdentifier    *Token
	CatchParameterObjectBindingPattern *ObjectBindingPattern
	CatchParameterArrayBindingPattern  *ArrayBindingPattern
	CatchBlock                         *Block
	FinallyBlock                       *Block
	Tokens                             Tokens
}

func (ts *TryStatement) parse(j *jsParser, yield, await, ret bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "try"}) {
		return j.Error("TryStatement", ErrInvalidTryStatement)
	}

	j.AcceptRunWhitespace()

	g := j.NewGoal()

	if err := ts.TryBlock.parse(&g, yield, await, ret); err != nil {
		return j.Error("TryStatement", err)
	}

	j.Score(g)
	j.AcceptRunWhitespace()

	if j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "catch"}) {
		j.AcceptRunWhitespace()

		if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
			j.AcceptRunWhitespace()

			g = j.NewGoal()

			switch g.Peek() {
			case parser.Token{Type: TokenPunctuator, Data: "{"}:
				ts.CatchParameterObjectBindingPattern = new(ObjectBindingPattern)
				if err := ts.CatchParameterObjectBindingPattern.parse(&g, yield, await); err != nil {
					return j.Error("TryStatement", err)
				}
			case parser.Token{Type: TokenPunctuator, Data: "["}:
				ts.CatchParameterArrayBindingPattern = new(ArrayBindingPattern)
				if err := ts.CatchParameterArrayBindingPattern.parse(&g, yield, await); err != nil {
					return j.Error("TryStatement", err)
				}
			default:
				if ts.CatchParameterBindingIdentifier = g.parseIdentifier(yield, await); ts.CatchParameterBindingIdentifier == nil {
					return j.Error("TryStatement", ErrNoIdentifier)
				}
			}

			j.Score(g)
			j.AcceptRunWhitespace()

			if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
				return j.Error("TryStatement", ErrMissingClosingParenthesis)
			}

			j.AcceptRunWhitespace()
		}

		g = j.NewGoal()

		ts.CatchBlock = new(Block)
		if err := ts.CatchBlock.parse(&g, yield, await, ret); err != nil {
			return j.Error("TryStatement", err)
		}

		j.Score(g)
	}

	g = j.NewGoal()

	g.AcceptRunWhitespace()

	if g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "finally"}) {
		g.AcceptRunWhitespace()

		h := g.NewGoal()

		ts.FinallyBlock = new(Block)
		if err := ts.FinallyBlock.parse(&h, yield, await, ret); err != nil {
			return g.Error("TryStatement", err)
		}

		g.Score(h)
		j.Score(g)
	}

	if ts.CatchBlock == nil && ts.FinallyBlock == nil {
		return j.Error("TryStatement", ErrMissingCatchFinally)
	}

	ts.Tokens = j.ToTokens()

	return nil
}

// VariableStatement as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-VariableStatement
//
// VariableDeclarationList must have a length or at least one.
type VariableStatement struct {
	VariableDeclarationList []VariableDeclaration
	Tokens                  Tokens
}

func (vs *VariableStatement) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "var"}) {
		return j.Error("VariableStatement", ErrInvalidVariableStatement)
	}

	for {
		j.AcceptRunWhitespaceNoComment()

		g := j.NewGoal()
		vd := len(vs.VariableDeclarationList)

		vs.VariableDeclarationList = append(vs.VariableDeclarationList, VariableDeclaration{})
		if err := vs.VariableDeclarationList[vd].parse(&g, true, yield, await); err != nil {
			return j.Error("VariableStatement", err)
		}

		j.Score(g)

		g = j.NewGoal()

		g.AcceptRunWhitespace()

		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			if j.parseSemicolon() {
				break
			}

			return j.Error("VariableStatement", ErrMissingComma)
		}

		j.Score(g)
	}

	vs.Tokens = j.ToTokens()

	return nil
}
