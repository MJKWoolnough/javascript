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
	RemoveDeadCode

	Safe = Literals | ArrowFn | IfToConditional | RemoveDebugger | RenameIdentifiers | BlocksToStatement | Keys | RemoveExpressionNames | FunctionExpressionToArrowFunc | UnwrapParens | RemoveLastEmptyReturn | CombineExpressionRuns | RemoveDeadCode
)

func (o Option) Has(opt Option) bool {
	return o&opt != 0
}
