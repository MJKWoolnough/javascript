package minify

type Option uint64

const (
	Literals = 1 << iota
	ArrowFn
	IfToConditional
	RemoveDebugger
	RenameIdentifiers
	BlocksToStatement
	Keys
	RemoveExpressionNames
	FunctionExpressionToArrowFunc
	UnwrapParens
	RemoveLastEmptyReturn
	CombineExpressionRuns
)
