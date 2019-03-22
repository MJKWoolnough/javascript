package javascript

import "vimagination.zapto.org/parser"

type FunctionDeclaration struct {
	BindingIdentifier *BindingIdentifier
	FormalParameters  FormalParameters
	FunctionBody      FunctionBody
	Tokens            []TokenPos
}

func (j *jsParser) parseFunctionDeclaration(yield, await, def bool) (FunctionDeclaration, error) {
	var fd FunctionDeclaration
	j.AcceptToken(parser.Token{TokenKeyword, "function"})
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	bi, err := g.parseBindingIdentifier(yield, await)
	if err != nil {
		if !def {
			return fd, j.Error(err)
		}
	} else {
		j.Score(g)
		fd.BindingIdentifier = &bi
		j.AcceptRunWhitespace()
	}
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return fd, j.Error(ErrMissingOpeningParentheses)
	}
	g = j.NewGoal()
	fd.FormalParameters, err = g.parseFormalParameters(false, false)
	if err != nil {
		return fd, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return fd, j.Error(ErrMissingClosingParentheses)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		return fd, j.Error(ErrMissingOpeningBrace)
	}
	g = j.NewGoal()
	fd.FunctionBody, err = j.parseFunctionBody(false, false)
	if err != nil {
		return fd, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.Accept(TokenRightBracePunctuator) {
		return fd, j.Error(ErrMissingClosingBrace)
	}
	fd.Tokens = j.ToTokens()
	return fd, nil
}

type AsyncFunctionDeclaration struct {
	BindingIdentifier *BindingIdentifier
	FormalParameters  FormalParameters
	FunctionBody      FunctionBody
	Tokens            []TokenPos
}

func (j *jsParser) parseAsyncFunctionDeclaration(yield, await, def bool) (AsyncFunctionDeclaration, error) {
	var af AsyncFunctionDeclaration
	j.AcceptToken(parser.Token{TokenKeyword, "async"})
	j.AcceptRunWhitespaceNoNewLine()
	if !j.AcceptToken(parser.Token{TokenKeyword, "function"}) {
		return af, j.Error(ErrMissingFunction)
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	bi, err := g.parseBindingIdentifier(yield, await)
	if err != nil {
		if !def {
			return af, j.Error(err)
		}
	} else {
		j.Score(g)
		af.BindingIdentifier = &bi
		j.AcceptRunWhitespace()
	}
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return af, j.Error(ErrMissingOpeningParentheses)
	}
	g = j.NewGoal()
	af.FormalParameters, err = g.parseFormalParameters(false, await)
	if err != nil {
		return af, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return af, j.Error(ErrMissingClosingParentheses)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		return af, j.Error(ErrMissingOpeningBrace)
	}
	g = j.NewGoal()
	af.FunctionBody, err = j.parseFunctionBody(false, true)
	if err != nil {
		return af, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.Accept(TokenRightBracePunctuator) {
		return af, j.Error(ErrMissingClosingBrace)
	}
	af.Tokens = j.ToTokens()
	return af, nil
}

type GeneratorDeclaration struct {
	BindingIdentifier *BindingIdentifier
	FormalParameters  FormalParameters
	FunctionBody      FunctionBody
	Tokens            []TokenPos
}

func (j *jsParser) parseGeneratorDeclaration(yield, await, def bool) (GeneratorDeclaration, error) {
	var gd GeneratorDeclaration
	j.AcceptToken(parser.Token{TokenKeyword, "function"})
	j.AcceptRunWhitespace()
	j.AcceptToken(parser.Token{TokenPunctuator, "*"})
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	bi, err := g.parseBindingIdentifier(yield, await)
	if err != nil {
		if !def {
			return gd, j.Error(err)
		}
	} else {
		j.Score(g)
		gd.BindingIdentifier = &bi
		j.AcceptRunWhitespace()
	}
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return gd, j.Error(ErrMissingOpeningParentheses)
	}
	g = j.NewGoal()
	gd.FormalParameters, err = g.parseFormalParameters(true, false)
	if err != nil {
		return gd, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return gd, j.Error(ErrMissingClosingParentheses)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		return gd, j.Error(ErrMissingOpeningBrace)
	}
	g = j.NewGoal()
	gd.FunctionBody, err = j.parseFunctionBody(true, false)
	if err != nil {
		return gd, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.Accept(TokenRightBracePunctuator) {
		return gd, j.Error(ErrMissingClosingBrace)
	}
	gd.Tokens = j.ToTokens()
	return gd, nil
}

type FormalParameters struct {
	FormalParameterList   []BindingElement
	FunctionRestParameter *FunctionRestParameter
	Tokens                []TokenPos
}

func (j *jsParser) parseFormalParameters(yield, await bool) (FormalParameters, error) {
	var fp FormalParameters
	for {
		g := j.NewGoal()
		g.AcceptRunWhitespace()
		if g.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
			break
		}
		if g.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
			fr, err := g.parseFunctionRestParameter(yield, await)
			if err != nil {
				return fp, j.Error(err)
			}
			j.Score(g)
			fp.FunctionRestParameter = &fr
			break
		}
		be, err := g.parseBindingElement(yield, await)
		if err != nil {
			return fp, err
		}
		j.Score(g)
		fp.FormalParameterList = append(fp.FormalParameterList, be)
		j.AcceptRunWhitespace()
		if j.Peek().Token == (parser.Token{TokenPunctuator, ")"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return fp, j.Error(ErrInvalidFormalParameterList)
		}
	}
	fp.Tokens = j.ToTokens()
	return fp, nil
}

type BindingElement struct {
	SingleNameBinding    *BindingIdentifier
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Initializer          *AssignmentExpression
	Tokens               []TokenPos
}

func (j *jsParser) parseBindingElement(yield, await bool) (BindingElement, error) {
	var be BindingElement
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{TokenPunctuator, "["}) {
		ab, err := g.parseArrayBindingPattern(yield, await)
		if err != nil {
			return be, j.Error(err)
		}
		be.ArrayBindingPattern = &ab
	} else if g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		ob, err := g.parseObjectBindingPattern(yield, await)
		if err != nil {
			return be, j.Error(err)
		}
		be.ObjectBindingPattern = &ob
	} else {
		bi, err := g.parseBindingIdentifier(yield, await)
		if err != nil {
			return be, j.Error(err)
		}
		be.SingleNameBinding = &bi
	}
	j.Score(g)
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenPunctuator, "="}) {
		g.AcceptRunWhitespace()
		j.Score(g)
		g = j.NewGoal()
		ae, err := g.parseAssignmentExpression(true, yield, await)
		if err != nil {
			return be, j.Error(err)
		}
		j.Score(g)
		be.Initializer = &ae
	}
	be.Tokens = j.ToTokens()
	return be, nil
}

type FunctionRestParameter struct {
	BindingIdentifier    *BindingIdentifier
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Tokens               []TokenPos
}

func (j *jsParser) parseFunctionRestParameter(yield, await bool) (FunctionRestParameter, error) {
	var fr FunctionRestParameter
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{TokenPunctuator, "["}) {
		ab, err := g.parseArrayBindingPattern(yield, await)
		if err != nil {
			return fr, j.Error(err)
		}
		fr.ArrayBindingPattern = &ab
	} else if g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		ob, err := g.parseObjectBindingPattern(yield, await)
		if err != nil {
			return fr, j.Error(err)
		}
		fr.ObjectBindingPattern = &ob
	} else {
		bi, err := g.parseBindingIdentifier(yield, await)
		if err != nil {
			return fr, j.Error(err)
		}
		fr.BindingIdentifier = &bi
	}
	j.Score(g)
	fr.Tokens = j.ToTokens()
	return fr, nil
}

type FunctionBody struct {
	Tokens []TokenPos
}

func (j *jsParser) parseFunctionBody(yield, await bool) (FunctionBody, error) {
	var fb FunctionBody
	return fb, nil
}
