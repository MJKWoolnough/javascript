package javascript

import (
	"io"
	"math"
	"strconv"
	"strings"

	"vimagination.zapto.org/parser"
)

// Token represents a parsed Token in the tree
type Token interface {
	io.WriterTo
	Data() interface{}
}

// Tokens represents a list of Tokens
type Tokens []Token

// Whitespace is sequential tabs, vertical tabs, form feeds, spaces, no-break
// spaces and zero width no-break spaces
type Whitespace string

// LineTerminators is sequential line feeds, carriage returns, line separators
// and paragraph separators
type LineTerminators string

// SingleLineComment is a comment started by two slashes
type SingleLineComment string

// MultiLineComment is a comment started by a slash and an asterix
type MultiLineComment string

// Identifier represents a non-keyword identified, usually a variable name or
// global
type Identifier string

// Boolean represents a literal boolean value (true, false)
type Boolean bool

// Keyword represents one of the following reserved keywords:
//
// await, break, case, catch, class, const, continue, debugger, default, delete,
// do, else, export, extends, finally, for, function, if, import, in,
// instanceof, new, return, super, switch, this, throw, try, typeof, var, void,
// while, with, yield
type Keyword string

// Punctuator represents one of the following:
//
// {, }, [, ], (, ), ;, ,, ?, :, ~, ., <, >, >=, <=, >>=, <<=, >>>, =, ==, ===,
// ==>, !, !=, !==, +, +=, -, -=, *, **, *=, /, /= &, &&, |, ||, &=, |=, %, %=,
// ^, ^=
type Punctuator string

// Number represents a number literal
type Number string

// NumberBinary represents a binary literal
type NumberBinary string

// NumberOctal represents an octal literal
type NumberOctal string

// NumberHexadecimal represents a hexadecimal literal
type NumberHexadecimal string

// String represents a string literal
type String string

// NoSubstitutionTemplate represents a template literal
type NoSubstitutionTemplate string

// Template represents a template with substitutions
type Template struct {
	TemplateStart
	TemplateMiddle Tokens
	TemplateEnd
}

// TemplateStart represents the opening of a template
type TemplateStart string

//TemplateMiddle represents a middle chunk of a template
type TemplateMiddle string

// TemplateEnd represents the end of a template
type TemplateEnd string

// Regex represents a regular expression literal
type Regex string

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the tokens
func (tk Tokens) WriteTo(w io.Writer) (int64, error) {
	var total int64
	for _, t := range tk {
		n, err := t.WriteTo(w)
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk Whitespace) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk LineTerminators) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk SingleLineComment) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk MultiLineComment) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk Identifier) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk Boolean) WriteTo(w io.Writer) (int64, error) {
	var (
		n   int
		err error
	)
	if tk {
		n, err = io.WriteString(w, "true")
	} else {
		n, err = io.WriteString(w, "false")
	}
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk Keyword) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk Punctuator) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk Number) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk NumberBinary) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk NumberOctal) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk NumberHexadecimal) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk String) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk NoSubstitutionTemplate) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk Template) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk.TemplateStart))
	if err != nil {
		return int64(n), err
	}
	m, err := tk.TemplateMiddle.WriteTo(w)
	m += int64(n)
	if err != nil {
		return m, err
	}
	n, err = io.WriteString(w, string(tk.TemplateEnd))
	return m + int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk TemplateStart) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk TemplateMiddle) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk TemplateEnd) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// WriteTo implements the io.WriterTo interface and writes to the Writer the
// original contents of the token
func (tk Regex) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

// Data returns typed interpreted data
func (tk Tokens) Data() interface{} {
	return tk
}

// Data returns typed interpreted data
func (tk Whitespace) Data() interface{} {
	return string(tk)
}

// Data returns typed interpreted data
func (tk LineTerminators) Data() interface{} {
	return string(tk)
}

// Data returns typed interpreted data
func (tk SingleLineComment) Data() interface{} {
	return tk.Comment()
}

// Comment returns the contents of the comment
func (tk SingleLineComment) Comment() string {
	return string(tk[2:])
}

// Data returns typed interpreted data
func (tk MultiLineComment) Data() interface{} {
	return tk.Comment()
}

// Comment returns the contents of the comment
func (tk MultiLineComment) Comment() string {
	return string(tk[2 : len(tk)-2])
}

// Data returns typed interpreted data
func (tk Identifier) Data() interface{} {
	return tk.String()
}

// String returns the unescaped string
func (tk Identifier) String() string {
	return unescape(string(tk))
}

// Data returns typed interpreted data
func (tk Boolean) Data() interface{} {
	return bool(tk)
}

// Data returns typed interpreted data
func (tk Keyword) Data() interface{} {
	return string(tk)
}

// Data returns typed interpreted data
func (tk Punctuator) Data() interface{} {
	return string(tk)
}

// Data returns typed interpreted data
func (tk Number) Data() interface{} {
	return tk.Number()
}

var pInf = math.Inf(1)

// Number returns a a numeric interpretation of the data
func (tk Number) Number() float64 {
	if tk == "Infinity" {
		return pInf
	}
	if pos := strings.IndexAny(string(tk), "Ee"); pos > 0 {
		co, _ := strconv.ParseFloat(string(tk[:pos]), 64)
		ex, _ := strconv.ParseInt(string(tk[pos+1:]), 10, 64)
		return math.Pow(co, math.Pow10(int(ex)))
	}
	co, _ := strconv.ParseFloat(string(tk), 64)
	return co
}

// Data returns typed interpreted data
func (tk NumberBinary) Data() interface{} {
	return tk.Number()
}

// Number returns a a numeric interpretation of the data
func (tk NumberBinary) Number() float64 {
	n, _ := strconv.ParseUint(string(tk[2:]), 2, 64)
	return float64(n)
}

// Data returns typed interpreted data
func (tk NumberOctal) Data() interface{} {
	return tk.Number()
}

// Number returns a a numeric interpretation of the data
func (tk NumberOctal) Number() float64 {
	n, _ := strconv.ParseUint(string(tk[2:]), 8, 64)
	return float64(n)
}

// Data returns typed interpreted data
func (tk NumberHexadecimal) Data() interface{} {
	return tk.Number()
}

// Number returns a a numeric interpretation of the data
func (tk NumberHexadecimal) Number() float64 {
	n, _ := strconv.ParseUint(string(tk[2:]), 16, 64)
	return float64(n)
}

// Data returns typed interpreted data
func (tk String) Data() interface{} {
	return tk.String()
}

// String returns the unescaped string
func (tk String) String() string {
	if len(tk) < 2 {
		return ""
	}
	return unescape(string(tk[1 : len(tk)-1]))
}

// Data returns typed interpreted data
func (tk NoSubstitutionTemplate) Data() interface{} {
	return tk.String()
}

// String returns the unescaped string
func (tk NoSubstitutionTemplate) String() string {
	return unescape(string(tk))
}

// Data returns typed interpreted data
func (tk Template) Data() interface{} {
	return tk
}

// Data returns typed interpreted data
func (tk TemplateStart) Data() interface{} {
	return tk.String()
}

// String returns the unescaped string
func (tk TemplateStart) String() string {
	if len(tk) < 3 {
		return ""
	}
	return unescape(string(tk[1 : len(tk)-2]))
}

// Data returns typed interpreted data
func (tk TemplateMiddle) Data() interface{} {
	return tk.String()
}

// String returns the unescaped string
func (tk TemplateMiddle) String() string {
	if len(tk) < 3 {
		return ""
	}
	return unescape(string(tk[1 : len(tk)-2]))
}

// Data returns typed interpreted data
func (tk TemplateEnd) Data() interface{} {
	return tk.String()
}

// String returns the unescaped string
func (tk TemplateEnd) String() string {
	if len(tk) < 2 {
		return ""
	}
	return unescape(string(tk[1 : len(tk)-1]))
}

// Data returns typed interpreted data
func (tk Regex) Data() interface{} {
	return tk
}

func unescape(str string) string {
	if !strings.ContainsRune(str, '\\') {
		return str
	}
	s := make([]byte, 0, len(str))
	p := parser.NewStringTokeniser(str)
	for {
		switch p.ExceptRun("\\") {
		case -1:
			s = append(s, p.Get()...)
			return string(s)
		case '\\':
			s = append(s, p.Get()...)
			p.Accept("\\")
			p.Get()
			p.Except("")
			switch b := p.Get(); b {
			case "'":
				s = append(s, '\'')
			case "\"":
				s = append(s, '"')
			case "\\":
				s = append(s, '\\')
			case "`":
				s = append(s, '`')
			case "b":
				s = append(s, '\b')
			case "f":
				s = append(s, '\f')
			case "n":
				s = append(s, '\n')
			case "r":
				s = append(s, '\r')
			case "t":
				s = append(s, '	')
			case "u":
				var n string
				if p.Accept("{") {
					p.Get()
					p.ExceptRun("}")
					n = p.Get()
					p.Accept("}")
				} else {
					p.Get()
					p.Accept(hexDigit)
					p.Accept(hexDigit)
					p.Accept(hexDigit)
					p.Accept(hexDigit)
					n = p.Get()
				}
				c, _ := strconv.ParseUint(n, 16, 64)
				s = append(s, string(rune(c))...)
			case "v":
				s = append(s, '\v')
			case "x":
				p.Accept(hexDigit)
				p.Accept(hexDigit)
				n, _ := strconv.ParseUint(p.Get(), 16, 8)
				s = append(s, byte(n))
			case "0":
				s = append(s, 0)
			default:
				s = append(s, '\\')
				s = append(s, b...)
			}
		}
	}
}
