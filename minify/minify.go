package minify

import (
	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/walk"
)

type Minifier struct{}

func New(opts ...Option) *Minifier {
	m := new(Minifier)
	for _, opt := range opts {
		opt(m)
	}
	return m
}

type walker struct {
	*Minifier
}

func (w *walker) Handle(t javascript.Type) error {
	return nil
}

func (m *Minifier) Process(jm *javascript.Module) {
	walk.Walk(jm, &walker{Minifier: m})
}
