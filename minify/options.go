package minify

type Option func(*Minifier)

func Literals() func(m *Minifier) {
	return func(m *Minifier) {
		m.literals = true
	}
}
