package javascript

import (
	"errors"
	"strconv"

	"vimagination.zapto.org/parser"
)

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
	var ret string
Loop:
	for {
		switch s.ExceptRun(chars) {
		case '"', '\'':
			ret += s.Get()
			return ret, nil
		case '\\':
			ret += s.Get()
			s.Accept("\\")
			s.Get()
			if s.Accept("x") {
				s.Get()
				if !s.Accept(hexDigit) || !s.Accept(hexDigit) {
					break Loop
				}
				c, _ := strconv.ParseUint(s.Get(), 16, 8)
				ret += string(rune(c))
			} else if s.Accept("u") {
				s.Get()
				if s.Accept("{") {
					s.Get()
					if !s.Accept(hexDigit) {
						break Loop
					}
					s.AcceptRun(hexDigit)
					c, _ := strconv.ParseUint(s.Get(), 16, 8)
					ret += string(rune(c))
					if !s.Accept("}") {
						break Loop
					}
				} else if !s.Accept(hexDigit) || !s.Accept(hexDigit) || !s.Accept(hexDigit) || !s.Accept(hexDigit) {
					break Loop
				} else {
					c, _ := strconv.ParseUint(s.Get(), 16, 8)
					ret += string(rune(c))
				}
			} else if s.Accept("0") {
				if s.Accept(decimalDigit) {
					break Loop
				}
				s.Get()
				ret += "\000"
			} else if s.Accept(singleEscapeChar) {
				switch s.Get() {
				case "'":
					ret += singleEscapeChar[0:1]
				case "\"":
					ret += singleEscapeChar[1:2]
				case "\\":
					ret += singleEscapeChar[2:3]
				case "b":
					ret += "\b"
				case "f":
					ret += "\f"
				case "n":
					ret += "\n"
				case "r":
					ret += "\r"
				case "t":
					ret += "\t"
				case "v":
					ret += "\v"
				default:
					break Loop
				}
			} else {
				break Loop
			}
		default:
			break Loop
		}
	}
	return "", ErrInvalidQuoted
}

// WrapConditional takes one of many types and wraps it in a
// *ConditionalExpression.
//
// The accepted types/pointers are as follows:
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
//    NewExpression
//    MemberExpression
//    PrimaryExpression
//
// Any other type will result in a panic.
func WrapConditional(p interface{}) *ConditionalExpression {
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
		panic("invalid conditional type")
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
//    *LeftHandSideExpression
//    *CallExpression
//    *NewExpression
//    *MemberExpression
//    *PrimaryExpression
func UnwrapConditional(c *ConditionalExpression) interface{} {
	if c.True != nil || c.LogicalORExpression == nil {
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
	} else if lhs.NewExpression.News > 0 {
		return lhs.NewExpression
	} else if lhs.NewExpression.MemberExpression.PrimaryExpression == nil {
		return &lhs.NewExpression.MemberExpression
	} else {
		return lhs.NewExpression.MemberExpression.PrimaryExpression
	}
}

var (
	ErrInvalidQuoted = errors.New("invalid quoted string")
)