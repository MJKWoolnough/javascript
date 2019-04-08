package javascript

type ConditionalExpression struct {
	LogicalORExpression LogicalORExpression
	True                *AssignmentExpression
	False               *AssignmentExpression
	Tokens              []TokenPos
}

func (j *jsParser) parseConditionalExpression(in, yield, await bool) (ConditionalExpression, error) {
	var ce ConditionalExpression
	return ce, nil
}

type LogicalORExpression struct {
	LogicalORExpression  *LogicalORExpression
	LogicalANDExpression LogicalANDExpression
	Tokens               []TokenPos
}

type LogicalANDExpression struct {
	LogicalANDExpression *LogicalANDExpression
	BitwiseORExpression  BitwiseORExpression
	Tokens               []TokenPos
}

type BitwiseORExpression struct {
	BitwiseORExpression  *BitwiseORExpression
	BitwiseXORExpression BitwiseXORExpression
	Tokens               []TokenPos
}

type BitwiseXORExpression struct {
	BitwiseXORExpression *BitwiseXORExpression
	BitwiseANDExpression BitwiseANDExpression
	Tokens               []TokenPos
}

type BitwiseANDExpression struct {
	BitwiseANDExpression *BitwiseANDExpression
	EqualityExpression   EqualityExpression
	Tokens               []TokenPos
}

type EqualityOperator int

const (
	EqualityNone EqualityOperator = iota
	EqualityEqual
	EqualityNotEqual
	EqualityStrictEqual
	EqualityStrictNotEqual
)

type EqualityExpression struct {
	EqualityExpression   *EqualityExpression
	EqualityOperator     EqualityOperator
	RelationalExpression RelationalExpression
	Tokens               []TokenPos
}

type RelationshipOperator int

const (
	RelationshipNone RelationshipOperator = iota
	RelationshipLessThan
	RelationshipGreaterThan
	RelationshipLessThanEqual
	RelationshipGreaterThanEqual
	RelationshipInstanceOf
	RelationshipIn
)

type RelationalExpression struct {
	RelationalExpression *RelationalExpression
	RelationshipOperator RelationshipOperator
	ShiftExpression      ShiftExpression
	Tokens               []TokenPos
}

type ShiftOperator int

const (
	ShiftNone ShiftOperator = iota
	ShiftLeft
	ShiftRight
	ShiftUnsignedRight
)

type ShiftExpression struct {
	ShiftExpression    *ShiftExpression
	ShiftOperator      ShiftOperator
	AdditiveExpression AdditiveExpression
	Tokens             []TokenPos
}

type AdditiveOperator int

const (
	AdditiveNone AdditiveOperator = iota
	AdditiveAdd
	AdditiveMinus
)

type AdditiveExpression struct {
	AdditiveExpression       *AdditiveExpression
	AdditiveOperator         AdditiveOperator
	MultiplicativeExpression MultiplicativeExpression
	Tokens                   []TokenPos
}

type MultiplicativeOperator int

const (
	MultiplicativeNone MultiplicativeOperator = iota
	MultiplicativeMultiply
	MultiplicativeDivide
	MultiplicativeRemainder
)

type MultiplicativeExpression struct {
	MultiplicativeExpression *MultiplicativeExpression
	MultiplicativeOperator   MultiplicativeOperator
	ExponentiationExpression ExponentiationExpression
	Tokens                   []TokenPos
}

type ExponentiationExpression struct {
	UnaryExpression          *UnaryExpression
	UpdateExpression         *UpdateExpression
	ExponentiationExpression *ExponentiationExpression
	Tokens                   []TokenPos
}

type UnaryOperator int

const (
	UnaryNone UnaryOperator = iota
	UnaryDelete
	UnaryVoid
	UnaryTypeof
	UnaryAdd
	UnaryMinus
	UnaryBitwiseNot
	UnaryLogicalNot
)

type UnaryExpression struct {
	UpdateExpression *UpdateExpression
	UnaryOperator    UnaryOperator
	AwaitExpression  *AwaitExpression
	Tokens           []TokenPos
}

type UpdateOperator int

const (
	UpdateNone UpdateOperator = iota
	UpdatePostIncrement
	UpdatePostDecrement
	UpdatePreIncrement
	UpdatePreDecrement
)

type UpdateExpression struct {
	LeftHandSideExpression *LeftHandSideExpression
	UnaryExpression        *UnaryExpression
	UpdateOperator         UpdateOperator
	Tokens                 []TokenPos
}

type AwaitExpression struct {
	UnaryExpression UnaryExpression
	Tokens          []TokenPos
}
