package javascript

import (
	"strconv"
	"strings"

	"vimagination.zapto.org/parser"
)

func unquoteEscape(ret *strings.Builder, s *parser.Tokeniser) bool {
	ret.WriteString(s.Get())
	s.Accept("\\")
	s.Get()

	if s.Accept("x") {
		s.Get()

		if !s.Accept(hexDigit) || !s.Accept(hexDigit) {
			return false
		}

		c, _ := strconv.ParseUint(s.Get(), 16, 8)

		ret.WriteString(string(rune(c)))
	} else if s.Accept("u") {
		s.Get()

		if s.Accept("{") {
			s.Get()

			if !s.Accept(hexDigit) {
				return false
			}

			s.AcceptRun(hexDigit)

			c, _ := strconv.ParseUint(s.Get(), 16, 8)

			ret.WriteString(string(rune(c)))

			if !s.Accept("}") {
				return false
			}

			s.Get()
		} else if !s.Accept(hexDigit) || !s.Accept(hexDigit) || !s.Accept(hexDigit) || !s.Accept(hexDigit) {
			return false
		} else {
			c, _ := strconv.ParseUint(s.Get(), 16, 8)

			ret.WriteString(string(rune(c)))
		}
	} else if s.Accept("0") {
		if s.Accept(decimalDigit) {
			return false
		}

		s.Get()
		ret.WriteString("\000")
	} else if s.Accept(singleEscapeChar) {
		switch s.Get() {
		case "'":
			ret.WriteString(singleEscapeChar[0:1])
		case "\"":
			ret.WriteString(singleEscapeChar[1:2])
		case "\\":
			ret.WriteString(singleEscapeChar[2:3])
		case "b":
			ret.WriteString("\b")
		case "f":
			ret.WriteString("\f")
		case "n":
			ret.WriteString("\n")
		case "r":
			ret.WriteString("\r")
		case "t":
			ret.WriteString("\t")
		case "v":
			ret.WriteString("\v")
		}
	} else {
		s.Get()
		s.Except("")
		ret.WriteString(s.Get())

		return true
	}

	return true
}

// Unquote parses a javascript quoted string and produces the unquoted version
func Unquote(str string) (string, error) {
	s := parser.NewStringTokeniser(str)

	var chars string

	if s.Accept("\"") {
		chars = doubleStringChars
	} else if s.Accept("'") {
		chars = singleStringChars
	} else {
		return "", ErrInvalidQuoted
	}

	s.Get()

	var ret strings.Builder

	ret.Grow(len(str))

Loop:
	for {
		switch s.ExceptRun(chars) {
		case '"', '\'':
			ret.WriteString(s.Get())

			return ret.String(), nil
		case '\\':
			if !unquoteEscape(&ret, &s) {
				break Loop
			}
		default:
			break Loop
		}
	}

	return "", ErrInvalidQuoted
}

// UnquoteTemplate parses a javascript template (either NoSubstitution Template, or any template part),
// and produces the unquoted version.
func UnquoteTemplate(t string) (string, error) {
	if strings.HasPrefix(t, "`") || strings.HasPrefix(t, "}") {
		t = t[1:]
	} else {
		return "", ErrInvalidQuoted
	}

	if strings.HasSuffix(t, "`") {
		t = t[:len(t)-1]
	} else if strings.HasSuffix(t, "${") {
		t = t[:len(t)-2]
	} else {
		return "", ErrInvalidQuoted
	}

	s := parser.NewStringTokeniser(t)

	var ret strings.Builder

	ret.Grow(len(t))

Loop:
	for {
		switch s.ExceptRun("$`\\") {
		case '$':
			s.Except("")

			if s.Accept("{") {
				break Loop
			}
		case '\\':
			if !unquoteEscape(&ret, &s) {
				break Loop
			}
		case '`':
			break Loop
		default:
			ret.WriteString(s.Get())

			return ret.String(), nil
		}
	}

	return "", ErrInvalidQuoted
}

// TemplateType determines the type of Template used in QuoteTemplate.
type TemplateType byte

const (
	TemplateNoSubstitution TemplateType = iota
	TemplateHead
	TemplateMiddle
	TemplateTail
)

// TokenTypeToTemplateType converts from a parser.TokenType to the appropriate
// TemplateType.
//
// Invalid TokenTypes return 255.
func TokenTypeToTemplateType(tokenType parser.TokenType) TemplateType {
	switch tokenType {
	case TokenNoSubstitutionTemplate:
		return TemplateNoSubstitution
	case TokenTemplateHead:
		return TemplateHead
	case TokenTemplateMiddle:
		return TemplateMiddle
	case TokenTemplateTail:
		return TemplateTail
	}

	return 255
}

// QuoteTemplate creates a minimally quoted template string.
//
// templateType determines the prefix and suffix.
//
// |  Template Type         |  Prefix  |  Suffix  |
// |------------------------|----------|----------|
// | TemplateNoSubstitution | "`"      | "`"      |
// | TemplateHead           | "`"      | "${"     |
// | TemplateMiddle         | "}"      | "}"      |
// | TemplateTail           | "}"      | "`"      |
func QuoteTemplate(t string, templateType TemplateType) string {
	l := len(t) + 2

	if templateType == TemplateHead || templateType == TemplateMiddle {
		l++
	}

	for n, r := range t {
		switch r {
		case '$':
			if len(t) > n && t[n+1] != '{' {
				break
			}

			fallthrough
		case '`', '\\':
			l++
		}
	}

	var ret strings.Builder

	ret.Grow(l)

	if templateType == TemplateMiddle || templateType == TemplateTail {
		ret.WriteByte('}')
	} else {
		ret.WriteByte('`')
	}

	for n, r := range t {
		switch r {
		case '$':
			if len(t) > n && t[n+1] != '{' {
				break
			}

			fallthrough
		case '`', '\\':
			ret.WriteByte('\\')
		}

		ret.WriteRune(r)
	}

	if templateType == TemplateHead || templateType == TemplateMiddle {
		ret.WriteString("${")
	} else {
		ret.WriteByte('`')
	}

	return ret.String()
}

// WrapConditional takes one of many types and wraps it in a
// *ConditionalExpression.
//
// The accepted types/pointers are as follows:
//    ConditionalExpression
//    LogicalORExpression
//    LogicalANDExpression
//    BitwiseORExpression
//    BitwiseXORExpression
//    BitwiseANDExpression
//    EqualityExpression
//    RelationalExpression
//    ShiftExpression
//    AdditiveExpression
//    MultiplicativeExpression
//    ExponentiationExpression
//    UnaryExpression
//    UpdateExpression
//    LeftHandSideExpression
//    CallExpression
//    OptionalExpression
//    NewExpression
//    MemberExpression
//    PrimaryExpression
//    ArrayLiteral
//    ObjectLiteral
//    FunctionDeclaration (FunctionExpression)
//    ClassDeclaration (ClassExpression)
//    TemplateLiteral
//    ParenthesizedExpression
func WrapConditional(p ConditionalWrappable) *ConditionalExpression {
	if c, ok := p.(*ConditionalExpression); ok {
		return c
	} else if c, ok := p.(ConditionalExpression); ok {
		return &c
	}

	c := &ConditionalExpression{
		LogicalORExpression: new(LogicalORExpression),
	}

	switch p := p.(type) {
	case *LogicalORExpression:
		c.LogicalORExpression = p

		goto logicalORExpression
	case LogicalORExpression:
		c.LogicalORExpression = &p

		goto logicalORExpression
	case *LogicalANDExpression:
		c.LogicalORExpression.LogicalANDExpression = *p

		goto logicalANDExpression
	case LogicalANDExpression:
		c.LogicalORExpression.LogicalANDExpression = p

		goto logicalANDExpression
	case *BitwiseORExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression = *p

		goto bitwiseORExpression
	case BitwiseORExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression = p

		goto bitwiseORExpression
	case *BitwiseXORExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression = *p

		goto bitwiseXORExpression
	case BitwiseXORExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression = p

		goto bitwiseXORExpression
	case *BitwiseANDExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression = *p

		goto bitwiseANDExpression
	case BitwiseANDExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression = p

		goto bitwiseANDExpression
	case *EqualityExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression = *p

		goto equalityExpression
	case EqualityExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression = p

		goto equalityExpression
	case *RelationalExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression = *p

		goto relationalExpression
	case RelationalExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression = p

		goto relationalExpression
	case *ShiftExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression = *p

		goto shiftExpression
	case ShiftExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression = p

		goto shiftExpression
	case *AdditiveExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression = *p

		goto additiveExpression
	case AdditiveExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression = p

		goto additiveExpression
	case *MultiplicativeExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression = *p

		goto multiplicativeExpression
	case MultiplicativeExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression = p

		goto multiplicativeExpression
	case *ExponentiationExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression = *p

		goto exponentiationExpression
	case ExponentiationExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression = p

		goto exponentiationExpression
	case *UnaryExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression = *p

		goto unaryExpression
	case UnaryExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression = p

		goto unaryExpression
	case *UpdateExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression = *p

		goto updateExpression
	case UpdateExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression = p

		goto updateExpression
	case *LeftHandSideExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression = p

	case LeftHandSideExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression = &p
	case *CallExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression = &LeftHandSideExpression{
			CallExpression: p,
			Tokens:         p.Tokens,
		}
	case CallExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression = &LeftHandSideExpression{
			CallExpression: &p,
			Tokens:         p.Tokens,
		}
	case *OptionalExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression = &LeftHandSideExpression{
			OptionalExpression: p,
			Tokens:             p.Tokens,
		}
	case OptionalExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression = &LeftHandSideExpression{
			OptionalExpression: &p,
			Tokens:             p.Tokens,
		}
	case *NewExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression = &LeftHandSideExpression{
			NewExpression: p,
			Tokens:        p.Tokens,
		}
	case NewExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression = &LeftHandSideExpression{
			NewExpression: &p,
			Tokens:        p.Tokens,
		}
	case *MemberExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression = &LeftHandSideExpression{
			NewExpression: &NewExpression{
				MemberExpression: *p,
				Tokens:           p.Tokens,
			},
			Tokens: p.Tokens,
		}
	case MemberExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression = &LeftHandSideExpression{
			NewExpression: &NewExpression{
				MemberExpression: p,
				Tokens:           p.Tokens,
			},
			Tokens: p.Tokens,
		}
	case *PrimaryExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression = &LeftHandSideExpression{
			NewExpression: &NewExpression{
				MemberExpression: MemberExpression{
					PrimaryExpression: p,
					Tokens:            p.Tokens,
				},
				Tokens: p.Tokens,
			},
			Tokens: p.Tokens,
		}
	case PrimaryExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression = &LeftHandSideExpression{
			NewExpression: &NewExpression{
				MemberExpression: MemberExpression{
					PrimaryExpression: &p,
					Tokens:            p.Tokens,
				},
				Tokens: p.Tokens,
			},
			Tokens: p.Tokens,
		}
	default:
		pe := new(PrimaryExpression)

		switch p := p.(type) {
		case *ArrayLiteral:
			pe.ArrayLiteral = p
			pe.Tokens = p.Tokens
		case ArrayLiteral:
			pe.ArrayLiteral = &p
			pe.Tokens = p.Tokens
		case *ObjectLiteral:
			pe.ObjectLiteral = p
			pe.Tokens = p.Tokens
		case ObjectLiteral:
			pe.ObjectLiteral = &p
			pe.Tokens = p.Tokens
		case *FunctionDeclaration:
			pe.FunctionExpression = p
			pe.Tokens = p.Tokens
		case FunctionDeclaration:
			pe.FunctionExpression = &p
			pe.Tokens = p.Tokens
		case *ClassDeclaration:
			pe.ClassExpression = p
			pe.Tokens = p.Tokens
		case ClassDeclaration:
			pe.ClassExpression = &p
			pe.Tokens = p.Tokens
		case *TemplateLiteral:
			pe.TemplateLiteral = p
			pe.Tokens = p.Tokens
		case TemplateLiteral:
			pe.TemplateLiteral = &p
			pe.Tokens = p.Tokens
		case *ParenthesizedExpression:
			pe.ParenthesizedExpression = p
			pe.Tokens = p.Tokens
		case ParenthesizedExpression:
			pe.ParenthesizedExpression = &p
			pe.Tokens = p.Tokens
		}

		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression = &LeftHandSideExpression{
			NewExpression: &NewExpression{
				MemberExpression: MemberExpression{
					PrimaryExpression: pe,
					Tokens:            pe.Tokens,
				},
				Tokens: pe.Tokens,
			},
			Tokens: pe.Tokens,
		}
	}

	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.Tokens
updateExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.Tokens
unaryExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.Tokens
exponentiationExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.Tokens
multiplicativeExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.Tokens
additiveExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.Tokens
shiftExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.Tokens
relationalExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.Tokens
equalityExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.Tokens
bitwiseANDExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.Tokens
bitwiseXORExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.Tokens
bitwiseORExpression:
	c.LogicalORExpression.LogicalANDExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.Tokens
logicalANDExpression:
	c.LogicalORExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.Tokens
logicalORExpression:
	c.Tokens = c.LogicalORExpression.Tokens

	return c
}

// UnwrapConditional returns the first value up the ConditionalExpression chain
// that contains all of the information required to rebuild the lower chain.
//
// Possible returns types are as follows:
//    *ConditionalExpression
//    *LogicalORExpression
//    *LogicalANDExpression
//    *BitwiseORExpression
//    *BitwiseXORExpression
//    *BitwiseANDExpression
//    *EqualityExpression
//    *RelationalExpression
//    *ShiftExpression
//    *AdditiveExpression
//    *MultiplicativeExpression
//    *ExponentiationExpression
//    *UnaryExpression
//    *UpdateExpression
//    *CallExpression
//    *OptionalExpression
//    *NewExpression
//    *MemberExpression
//    *PrimaryExpression
//    *ArrayLiteral
//    *ObjectLiteral
//    *FunctionDeclaration (FunctionExpression)
//    *ClassDeclaration (ClassExpression)
//    *TemplateLiteral
//    *ParenthesizedExpression
func UnwrapConditional(c *ConditionalExpression) ConditionalWrappable {
	if c == nil {
		return nil
	} else if c.True != nil || c.LogicalORExpression == nil {
		return c
	} else if c.LogicalORExpression.LogicalORExpression != nil {
		return c.LogicalORExpression
	} else if c.LogicalORExpression.LogicalANDExpression.LogicalANDExpression != nil {
		return &c.LogicalORExpression.LogicalANDExpression
	} else if c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseORExpression != nil {
		return &c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression
	} else if c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseXORExpression != nil {
		return &c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression
	} else if c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.BitwiseANDExpression != nil {
		return &c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression
	} else if c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.EqualityExpression != nil {
		return &c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression
	} else if c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.RelationalExpression != nil {
		return &c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression
	} else if c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.ShiftExpression != nil {
		return &c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression
	} else if c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.AdditiveExpression != nil {
		return &c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression
	} else if c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.MultiplicativeExpression != nil {
		return &c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression
	} else if c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.ExponentiationExpression != nil {
		return &c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression
	} else if len(c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UnaryOperators) > 0 {
		return &c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression
	} else if c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression == nil || c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.UpdateOperator != UpdateNone {
		return &c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression
	} else if lhs := c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression; lhs.CallExpression != nil {
		return lhs.CallExpression
	} else if lhs.OptionalExpression != nil {
		return lhs.OptionalExpression
	} else if lhs.NewExpression.News > 0 {
		return lhs.NewExpression
	} else if lhs.NewExpression.MemberExpression.PrimaryExpression == nil {
		return &lhs.NewExpression.MemberExpression
	} else {
		pe := lhs.NewExpression.MemberExpression.PrimaryExpression

		if pe.ArrayLiteral != nil {
			return pe.ArrayLiteral
		} else if pe.ObjectLiteral != nil {
			return pe.ObjectLiteral
		} else if pe.FunctionExpression != nil {
			return pe.FunctionExpression
		} else if pe.ClassExpression != nil {
			return pe.ClassExpression
		} else if pe.TemplateLiteral != nil {
			return pe.TemplateLiteral
		} else if pe.ParenthesizedExpression != nil {
			return pe.ParenthesizedExpression
		}

		return pe
	}
}
