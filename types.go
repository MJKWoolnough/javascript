package javascript

import (
	"io"
	"math"
	"strconv"
)

type Token interface {
	io.WriterTo
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

type Number float64

type String string

type NoSubstitutionTemplate string

type Template Tokens

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
	if math.IsInf(float64(tk), 1) {
		n, err := io.WriteString(w, "Infinity")
		return int64(n), err
	}
	n, err := io.WriteString(w, strconv.FormatFloat(float64(tk), 'f', -1, 64))
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
	return Tokens(tk).WriteTo(w)
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
