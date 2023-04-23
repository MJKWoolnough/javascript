package minify

type Option func(*minifier)

func Literals(m *minifier) {
	m.literals = true
}

func Numbers(m *minifier) {
	m.numbers = true
}

func ArrowFn(m *minifier) {
	m.arrowFn = true
}

func IfToConditional(m *minifier) {
	m.ifToConditional = true
}

func RemoveDebugger(m *minifier) {
	m.rmDebugger = true
}
