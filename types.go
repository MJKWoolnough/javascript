package javascript

import (
	"io"
	"math"
	"strconv"
	"strings"

	"vimagination.zapto.org/parser"
)

type Token interface {
	io.WriterTo
	Data() interface{}
}

type Tokens []Token

type Whitespace string

type LineTerminators string

type SingleLineComment string

type MultiLineComment string

type Identifier string

type Boolean bool

type Keyword string

type Punctuator string

type Number string

type NumberBinary string

type NumberOctal string

type NumberHexadecimal string

type String string

type NoSubstitutionTemplate string

type Template struct {
	TemplateStart
	TemplateMiddle Tokens
	TemplateEnd
}

type TemplateStart string

type TemplateMiddle string

type TemplateEnd string

type Regex string

func (t Tokens) WriteTo(w io.Writer) (int64, error) {
	var total int64
	for _, tk := range t {
		n, err := tk.WriteTo(w)
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

func (tk Whitespace) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

func (tk LineTerminators) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

func (tk SingleLineComment) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

func (tk MultiLineComment) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

func (tk Identifier) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

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

func (tk Keyword) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

func (tk Punctuator) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

func (tk Number) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

func (tk NumberBinary) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

func (tk NumberOctal) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

func (tk NumberHexadecimal) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

func (tk String) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

func (tk NoSubstitutionTemplate) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

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

func (tk TemplateStart) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

func (tk TemplateMiddle) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

func (tk TemplateEnd) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

func (tk Regex) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(tk))
	return int64(n), err
}

func (tk Tokens) Data() interface{} {
	return tk
}

func (tk Whitespace) Data() interface{} {
	return string(tk)
}

func (tk LineTerminators) Data() interface{} {
	return string(tk)
}

func (tk SingleLineComment) Data() interface{} {
	return tk.Comment()
}

func (tk SingleLineComment) Comment() string {
	return string(tk[2:])
}

func (tk MultiLineComment) Data() interface{} {
	return tk.Comment()
}

func (tk MultiLineComment) Comment() string {
	return string(tk[2 : len(tk)-2])
}

func (tk Identifier) Data() interface{} {
	return tk.String()
}

func (tk Identifier) String() string {
	return unescape(string(tk))
}

func (tk Boolean) Data() interface{} {
	return bool(tk)
}

func (tk Keyword) Data() interface{} {
	return string(tk)
}

func (tk Punctuator) Data() interface{} {
	return string(tk)
}

func (tk Number) Data() interface{} {
	return tk.Number()
}

func (tk Number) Number() float64 {
	if pos := strings.IndexAny(string(tk), "Ee"); pos > 0 {
		co, _ := strconv.ParseFloat(string(tk[:pos]), 64)
		ex, _ := strconv.ParseInt(string(tk[pos+1:]), 10, 64)
		return math.Pow(co, math.Pow10(int(ex)))
	}
	co, _ := strconv.FormatFloat(string(tk), 64)
	return co
}

func (tk NumberBinary) Data() interface{} {
	return tk.Number()
}

func (tk NumberBinary) Number() float64 {
	n, _ := strconv.ParseUint(string(tk[2:]), 2, 64)
	return float64(n)
}

func (tk NumberOctal) Data() interface{} {
	return tk.Number()
}

func (tk NumberOctal) Number() float64 {
	n, _ := strconv.ParseUint(string(tk[2:]), 8, 64)
	return float64(n)
}

func (tk NumberHexadecimal) Data() interface{} {
	return tk.Number()
}

func (tk NumberHexadecimal) Number() float64 {
	n, _ := strconv.ParseUint(string(tk[2:]), 16, 64)
	return float64(n)
}

func (tk String) Data() interface{} {
	return tk.String()
}

func (tk String) String() string {
	if len(tk) < 2 {
		return ""
	}
	return unescape(string(tk[1 : len(tk)-1]))
}

func (tk NoSubstitutionTemplate) Data() interface{} {
	return tk.String()
}

func (tk NoSubstitutionTemplate) String() string {
	return unescape(string(tk))
}

func (tk Template) Data() interface{} {
	return tk
}

func (tk TemplateStart) Data() interface{} {
	return tk.String()
}

func (tk TemplateStart) String() string {
	if len(tk) < 3 {
		return ""
	}
	return unescape(string(tk[1 : len(tk)-2]))
}

func (tk TemplateMiddle) Data() interface{} {
	return tk.String()
}

func (tk TemplateMiddle) String() string {
	if len(tk) < 3 {
		return ""
	}
	return unescape(string(tk[1 : len(tk)-2]))
}

func (tk TemplateEnd) Data() interface{} {
	return tk.String()
}

func (tk TemplateEnd) String() string {
	if len(tk) < 2 {
		return ""
	}
	return unescape(string(tk[1 : len(tk)-1]))
}

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
