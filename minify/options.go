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

func RenameIdentifiers(m *minifier) {
	m.rename = true
}

func BlocksToStatement(m *minifier) {
	m.blocks = true
}

func Keys(m *minifier) {
	m.keys = true
}

func RemoveExpressionNames(m *minifier) {
	m.nonHoistableNames = true
}

func FunctionExpressionToArrowFunc(m *minifier) {
	m.replaceFEWithAF = true
}

func UnwrapParens(m *minifier) {
	m.unwrapParens = true
}

func RemoveLastEmptyReturn(m *minifier) {
	m.removeLastReturn = true
}

func CombineExpressionRuns(m *minifier) {
	m.combineExpressions = true
}
