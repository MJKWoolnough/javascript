package jsx

import (
	"errors"
	"text/template"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/scope"
	"vimagination.zapto.org/javascript/walk"
)

type jsxWalker struct {
	identifiers map[string]map[string][]scope.Binding
	tmpl        *template.Template
	namespace   string
}

func (j *jsxWalker) Handle(t javascript.Type) error {
	return nil
}

func Process(m *javascript.Module, tmpl *template.Template) error {
	j := &jsxWalker{
		identifiers: make(map[string]map[string][]scope.Binding),
		tmpl:        tmpl,
	}

	if err := walk.Walk(m, j); err != nil {
		return err
	}

	return nil
}

var (
	ErrInvalidTransformation = errors.New("invalid transformation")
	ErrTooManyStatements     = errors.New("too many statments")
)
