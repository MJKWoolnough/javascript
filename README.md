# javascript
--
    import "vimagination.zapto.org/javascript"


## Usage

```go
const (
	TokenWhitespace parser.TokenType = iota
	TokenLineTerminator
	TokenSingleLineComment
	TokenMultiLineComment
	TokenIdentifier
	TokenBooleanLiteral
	TokenKeyword
	TokenPunctuator
	TokenNumericLiteral
	TokenStringLiteral
	TokenNoSubstitutionTemplate
	TokenTemplateHead
	TokenTemplateMiddle
	TokenTemplateTail
	TokenDivPunctuator
	TokenRightBracePunctuator
	TokenRegularExpressionLiteral
	TokenNullLiteral
	TokenFutureReservedWord
)
```
Javascript Token values

```go
var (
	ErrReservedIdentifier        = errors.New("reserved identifier")
	ErrNoIdentifier              = errors.New("missing identifier")
	ErrMissingFunction           = errors.New("missing function")
	ErrMissingOpeningParenthesis = errors.New("missing opening parenthesis")
	ErrMissingClosingParenthesis = errors.New("missing closing parenthesis")
	ErrMissingOpeningBrace       = errors.New("missing opening brace")
	ErrMissingClosingBrace       = errors.New("missing closing brace")
	ErrMissingOpeningBracket     = errors.New("missing opening bracket")
	ErrMissingClosingBracket     = errors.New("missing closing bracket")
	ErrMissingComma              = errors.New("missing comma")
	ErrMissingArrow              = errors.New("missing arrow")
	ErrMissingCaseClause         = errors.New("missing case clause")
	ErrMissingExpression         = errors.New("missing expression")
	ErrMissingCatchFinally       = errors.New("missing catch/finally block")
	ErrMissingSemiColon          = errors.New("missing semi-colon")
	ErrMissingColon              = errors.New("missing colon")
	ErrMissingInitializer        = errors.New("missing initializer")
	ErrInvalidStatementList      = errors.New("invalid statement list")
	ErrInvalidStatement          = errors.New("invalid statement")
	ErrInvalidDeclaration        = errors.New("invalid declaration")
	ErrInvalidLexicalDeclaration = errors.New("invalid lexical declaration")
	ErrInvalidAssignment         = errors.New("invalid assignment operator")
	ErrInvalidSuperProperty      = errors.New("invalid super property")
	ErrInvalidMetaProperty       = errors.New("invalid meta property")
	ErrInvalidTemplate           = errors.New("invalid template")
	ErrInvalidAsyncArrowFunction = errors.New("invalid async arrow function")
)
```
Errors

```go
var (
	ErrInvalidMethodName       = errors.New("invalid method name")
	ErrInvalidPropertyName     = errors.New("invalid property name")
	ErrInvalidClassDeclaration = errors.New("invalid class declaration")
)
```
Errors

```go
var (
	ErrInvalidCallExpression = errors.New("invalid CallExpression")
	ErrMissingOptional       = errors.New("missing optional chain punctuator")
	ErrInvalidOptionalChain  = errors.New("invalid OptionalChain")
)
```
Errors

```go
var (
	ErrInvalidImport            = errors.New("invalid import statement")
	ErrInvalidExportDeclaration = errors.New("invalid export declaration")
	ErrInvalidNameSpaceImport   = errors.New("invalid namespace import")
	ErrMissingFrom              = errors.New("missing from")
	ErrMissingModuleSpecifier   = errors.New("missing module specifier")
	ErrInvalidNamedImport       = errors.New("invalid named import list")
	ErrInvalidImportSpecifier   = errors.New("invalid import specifier")
	ErrInvalidExportClause      = errors.New("invalid export clause")
)
```
Errors

```go
var (
	ErrDuplicateDefaultClause         = errors.New("duplicate default clause")
	ErrInvalidIterationStatementDo    = errors.New("invalid do interation statement")
	ErrInvalidIterationStatementWhile = errors.New("invalid while interation statement")
	ErrInvalidIterationStatementFor   = errors.New("invalid for interation statement")
	ErrInvalidForLoop                 = errors.New("invalid for loop")
	ErrInvalidForAwaitLoop            = errors.New("invalid for await loop")
	ErrInvalidIfStatement             = errors.New("invalid if statement")
	ErrInvalidSwitchStatement         = errors.New("invalid switch statement")
	ErrInvalidWithStatement           = errors.New("invalid with statement")
	ErrInvalidTryStatement            = errors.New("invalid try statement")
	ErrInvalidVariableStatement       = errors.New("invalid variabl statement")
)
```
Errors

```go
var (
	ErrInvalidCharacter         = errors.New("invalid character")
	ErrInvalidSequence          = errors.New("invalid character sequence")
	ErrInvalidRegexpCharacter   = errors.New("invalid regexp character")
	ErrInvalidRegexpSequence    = errors.New("invalid regexp sequence")
	ErrInvalidNumber            = errors.New("invalid number")
	ErrUnexpectedBackslash      = errors.New("unexpected backslash")
	ErrInvalidUnicode           = errors.New("invalid unicode escape sequence")
	ErrInvalidEscapeSequence    = errors.New("invalid escape sequence")
	ErrUnexpectedLineTerminator = errors.New("line terminator in string")
)
```
Errors

```go
var (
	ErrInvalidFunction = errors.New("invalid function")
)
```
Errors

```go
var (
	ErrInvalidQuoted = errors.New("invalid quoted string")
)
```

#### func  SetTokeniser

```go
func SetTokeniser(t *parser.Tokeniser) *parser.Tokeniser
```
SetTokeniser provides javascript parsing functions to a Tokeniser

#### func  Unquote

```go
func Unquote(str string) (string, error)
```
Unquote parses a javascript quoted string and produces the unquoted version

#### func  UnwrapConditional

```go
func UnwrapConditional(c *ConditionalExpression) interface{}
```
UnwrapConditional returns the first value up the ConditionalExpression chain
that contains all of the information required to rebuild the lower chain.

Possible returns types are as follows:

    *ConditionalExpression
    *LogicalORExpression
    *LogicalANDExpression
    *BitwiseORExpression
    *BitwiseXORExpression
    *BitwiseANDExpression
    *EqualityExpression
    *RelationalExpression
    *ShiftExpression
    *AdditiveExpression
    *MultiplicativeExpression
    *ExponentiationExpression
    *UnaryExpression
    *UpdateExpression
    *LeftHandSideExpression
    *CallExpression
    *NewExpression
    *MemberExpression
    *PrimaryExpression

#### type AdditiveExpression

```go
type AdditiveExpression struct {
	AdditiveExpression       *AdditiveExpression
	AdditiveOperator         AdditiveOperator
	MultiplicativeExpression MultiplicativeExpression
	Tokens                   Tokens
}
```

AdditiveExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-AdditiveExpression

If AdditiveOperator is not AdditiveNone then AdditiveExpression must be non-nil,
and vice-versa.

#### func (AdditiveExpression) Format

```go
func (f AdditiveExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type AdditiveOperator

```go
type AdditiveOperator int
```

AdditiveOperator determines the additive type for AdditiveExpression

```go
const (
	AdditiveNone AdditiveOperator = iota
	AdditiveAdd
	AdditiveMinus
)
```
Valid AdditiveOperator's

#### func (AdditiveOperator) String

```go
func (a AdditiveOperator) String() string
```
String implements the fmt.Stringer interface

#### type Arguments

```go
type Arguments struct {
	ArgumentList   []AssignmentExpression
	SpreadArgument *AssignmentExpression
	Tokens         Tokens
}
```

Arguments as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-Arguments

#### func (Arguments) Format

```go
func (f Arguments) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ArrayBindingPattern

```go
type ArrayBindingPattern struct {
	BindingElementList []BindingElement
	BindingRestElement *BindingElement
	Tokens             Tokens
}
```

ArrayBindingPattern as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-ArrayBindingPattern

#### func (ArrayBindingPattern) Format

```go
func (f ArrayBindingPattern) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ArrayLiteral

```go
type ArrayLiteral struct {
	ElementList   []AssignmentExpression
	SpreadElement *AssignmentExpression
	Tokens        Tokens
}
```

ArrayLiteral as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-ArrayLiteral

#### func (ArrayLiteral) Format

```go
func (f ArrayLiteral) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ArrowFunction

```go
type ArrowFunction struct {
	Async                                             bool
	BindingIdentifier                                 *Token
	CoverParenthesizedExpressionAndArrowParameterList *CoverParenthesizedExpressionAndArrowParameterList
	FormalParameters                                  *FormalParameters
	AssignmentExpression                              *AssignmentExpression
	FunctionBody                                      *Block
	Tokens                                            Tokens
}
```

ArrowFunction as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-ArrowFunction

Also includes AsyncArrowFunction.

It is only valid for one of BindingIdentifier,
CoverParenthesizedExpressionAndArrowParameterList (!Async), and FormalParameters
(+Async) to be non-nil.

Only one of AssignmentExpression or FunctionBody must be non-nil.

#### func (ArrowFunction) Format

```go
func (f ArrowFunction) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type AssignmentExpression

```go
type AssignmentExpression struct {
	ConditionalExpression  *ConditionalExpression
	ArrowFunction          *ArrowFunction
	LeftHandSideExpression *LeftHandSideExpression
	Yield                  bool
	Delegate               bool
	AssignmentOperator     AssignmentOperator
	AssignmentExpression   *AssignmentExpression
	Tokens                 Tokens
}
```

AssignmentExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-AssignmentExpression

It is only valid for one of ConditionalExpression, ArrowFunction,
LeftHandSideExpression to be non-nil.

If LeftHandSideExpression is non-nil, then AssignmentOperator must not be
AssignmentNone and AssignmentExpression must be non-nil.

If Yield is true, AssignmentExpression must be non-nil.

If AssignmentOperator is AssignmentNone LeftHandSideExpression must be nil.

If LeftHandSideExpression is nil and Yield is false, AssignmentExpression must
be nil.

#### func (AssignmentExpression) Format

```go
func (f AssignmentExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type AssignmentOperator

```go
type AssignmentOperator uint8
```

AssignmentOperator specifies the type of assignment in AssignmentExpression

```go
const (
	AssignmentNone AssignmentOperator = iota
	AssignmentAssign
	AssignmentMultiply
	AssignmentDivide
	AssignmentRemainder
	AssignmentAdd
	AssignmentSubtract
	AssignmentLeftShift
	AssignmentSignPropagatinRightShift
	AssignmentZeroFillRightShift
	AssignmentBitwiseAND
	AssignmentBitwiseXOR
	AssignmentBitwiseOR
	AssignmentExponentiation
)
```
Valid AssignmentOperator's

#### func (AssignmentOperator) String

```go
func (a AssignmentOperator) String() string
```
String implements the fmt.Stringer interface

#### type BindingElement

```go
type BindingElement struct {
	SingleNameBinding    *Token
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Initializer          *AssignmentExpression
	Tokens               Tokens
}
```

BindingElement as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-BindingElement

Only one of SingleNameBinding, ArrayBindingPattern, or ObjectBindingPattern must
be non-nil.

The Initializer is optional.

#### func (BindingElement) Format

```go
func (f BindingElement) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type BindingProperty

```go
type BindingProperty struct {
	PropertyName   PropertyName
	BindingElement BindingElement
	Tokens         Tokens
}
```

BindingProperty as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-BindingProperty

A SingleNameBinding, with or without an initializer, is cloned into the Property
Name and Binding Element. This allows the Binding Element Identifier to be
modified while keeping the correct Property Name

#### func (BindingProperty) Format

```go
func (f BindingProperty) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type BitwiseANDExpression

```go
type BitwiseANDExpression struct {
	BitwiseANDExpression *BitwiseANDExpression
	EqualityExpression   EqualityExpression
	Tokens               Tokens
}
```

BitwiseANDExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-BitwiseANDExpression

#### func (BitwiseANDExpression) Format

```go
func (f BitwiseANDExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type BitwiseORExpression

```go
type BitwiseORExpression struct {
	BitwiseORExpression  *BitwiseORExpression
	BitwiseXORExpression BitwiseXORExpression
	Tokens               Tokens
}
```

BitwiseORExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-BitwiseORExpression

#### func (BitwiseORExpression) Format

```go
func (f BitwiseORExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type BitwiseXORExpression

```go
type BitwiseXORExpression struct {
	BitwiseXORExpression *BitwiseXORExpression
	BitwiseANDExpression BitwiseANDExpression
	Tokens               Tokens
}
```

BitwiseXORExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-BitwiseXORExpression

#### func (BitwiseXORExpression) Format

```go
func (f BitwiseXORExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type Block

```go
type Block struct {
	StatementList []StatementListItem
	Tokens        Tokens
}
```

Block as defined in ECMA-262 https://262.ecma-international.org/11.0/#prod-Block

#### func (Block) Format

```go
func (f Block) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type CallExpression

```go
type CallExpression struct {
	MemberExpression *MemberExpression
	SuperCall        bool
	ImportCall       *AssignmentExpression
	CallExpression   *CallExpression
	Arguments        *Arguments
	Expression       *Expression
	IdentifierName   *Token
	TemplateLiteral  *TemplateLiteral
	Tokens           Tokens
}
```

CallExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-CallExpression

Includes the TC39 proposal for the dynamic import function call
https://github.com/tc39/proposal-dynamic-import/#import

It is only valid for one of MemberExpression, ImportCall, or CallExpression to
be non-nil or SuperCall to be true.

If MemberExpression is non-nil, or SuperCall is true, Arguments must be non-nil.

If CallExpression is non-nil, only one of Arguments, Expression, IdentifierName,
or TemplateLiteral must be non-nil.

#### func (CallExpression) Format

```go
func (f CallExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type CaseClause

```go
type CaseClause struct {
	Expression    Expression
	StatementList []StatementListItem
	Tokens        Tokens
}
```

CaseClause as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-CaseClauses

#### func (CaseClause) Format

```go
func (f CaseClause) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ClassDeclaration

```go
type ClassDeclaration struct {
	BindingIdentifier *Token
	ClassHeritage     *LeftHandSideExpression
	ClassBody         []MethodDefinition
	Tokens            Tokens
}
```

ClassDeclaration as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-ClassDeclaration

Also covers ClassExpression when BindingIdentifier is nil.

#### func (ClassDeclaration) Format

```go
func (f ClassDeclaration) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type CoalesceExpression

```go
type CoalesceExpression struct {
	CoalesceExpressionHead *CoalesceExpression
	BitwiseORExpression    BitwiseORExpression
	Tokens                 Tokens
}
```

CoalesceExpression as defined in TC39
https://tc39.es/ecma262/#prod-CoalesceExpression

#### func (CoalesceExpression) Format

```go
func (f CoalesceExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ConditionalExpression

```go
type ConditionalExpression struct {
	LogicalORExpression *LogicalORExpression
	CoalesceExpression  *CoalesceExpression
	True                *AssignmentExpression
	False               *AssignmentExpression
	Tokens              Tokens
}
```

ConditionalExpression as defined in TC39
https://tc39.es/ecma262/#prod-ConditionalExpression

One, and only one, of LogicalORExpression or CoalesceExpression must be non-nil

If True is non-nil, False must be non-nil also.

#### func  WrapConditional

```go
func WrapConditional(p interface{}) *ConditionalExpression
```
WrapConditional takes one of many types and wraps it in a
*ConditionalExpression.

The accepted types/pointers are as follows:

    LogicalORExpression
    LogicalANDExpression
    BitwiseORExpression
    BitwiseXORExpression
    BitwiseANDExpression
    EqualityExpression
    RelationalExpression
    ShiftExpression
    AdditiveExpression
    MultiplicativeExpression
    ExponentiationExpression
    UnaryExpression
    UpdateExpression
    LeftHandSideExpression
    CallExpression
    NewExpression
    MemberExpression
    PrimaryExpression

Any other type will result in a panic.

#### func (ConditionalExpression) Format

```go
func (f ConditionalExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type CoverParenthesizedExpressionAndArrowParameterList

```go
type CoverParenthesizedExpressionAndArrowParameterList struct {
	Expressions          []AssignmentExpression
	BindingIdentifier    *Token
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Tokens               Tokens
}
```

CoverParenthesizedExpressionAndArrowParameterList as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-CoverParenthesizedExpressionAndArrowParameterList

It is valid for only one of BindingIdentifier, ArrayBindingPattern, and
ObjectBindingPattern to be non-nil

#### func (CoverParenthesizedExpressionAndArrowParameterList) Format

```go
func (f CoverParenthesizedExpressionAndArrowParameterList) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type Declaration

```go
type Declaration struct {
	ClassDeclaration    *ClassDeclaration
	FunctionDeclaration *FunctionDeclaration
	LexicalDeclaration  *LexicalDeclaration
	Tokens              Tokens
}
```

Declaration as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-Declaration

Only one of ClassDeclaration, FunctionDeclaration or LexicalDeclaration must be
non-nil

#### func (Declaration) Format

```go
func (f Declaration) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type EqualityExpression

```go
type EqualityExpression struct {
	EqualityExpression   *EqualityExpression
	EqualityOperator     EqualityOperator
	RelationalExpression RelationalExpression
	Tokens               Tokens
}
```

EqualityExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-EqualityExpression

If EqualityOperator is not EqualityNone, then EqualityExpression must be
non-nil, and vice-versa.

#### func (EqualityExpression) Format

```go
func (f EqualityExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type EqualityOperator

```go
type EqualityOperator int
```

EqualityOperator determines the type of EqualityExpression

```go
const (
	EqualityNone EqualityOperator = iota
	EqualityEqual
	EqualityNotEqual
	EqualityStrictEqual
	EqualityStrictNotEqual
)
```
Valid EqualityOperator's

#### func (EqualityOperator) String

```go
func (e EqualityOperator) String() string
```
String implements the fmt.Stringer interface

#### type Error

```go
type Error struct {
	Err     error
	Parsing string
	Token   Token
}
```

Error is a parsing error with trace details

#### func (Error) Error

```go
func (e Error) Error() string
```
Error returns the error string

#### func (Error) Unwrap

```go
func (e Error) Unwrap() error
```
Unwrap returns the wrapped error

#### type ExponentiationExpression

```go
type ExponentiationExpression struct {
	ExponentiationExpression *ExponentiationExpression
	UnaryExpression          UnaryExpression
	Tokens                   Tokens
}
```

ExponentiationExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-ExponentiationExpression

#### func (ExponentiationExpression) Format

```go
func (f ExponentiationExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ExportClause

```go
type ExportClause struct {
	ExportList []ExportSpecifier
	Tokens     Tokens
}
```

ExportClause as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-ExportClause

#### func (ExportClause) Format

```go
func (f ExportClause) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ExportDeclaration

```go
type ExportDeclaration struct {
	ExportClause                *ExportClause
	FromClause                  *FromClause
	VariableStatement           *VariableStatement
	Declaration                 *Declaration
	DefaultFunction             *FunctionDeclaration
	DefaultClass                *ClassDeclaration
	DefaultAssignmentExpression *AssignmentExpression
	Tokens                      Tokens
}
```

ExportDeclaration as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-ExportDeclaration

It is only valid for one of ExportClause, VariableStatement, Declaration,
DefaultFunction, DefaultClass, or DefaultAssignmentExpression to be non-nil.

FromClause can be non-nil exclusively or paired with ExportClause.

#### func (ExportDeclaration) Format

```go
func (f ExportDeclaration) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ExportSpecifier

```go
type ExportSpecifier struct {
	IdentifierName, EIdentifierName *Token
	Tokens                          Tokens
}
```

ExportSpecifier as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-ExportSpecifier

IdentifierName must be non-nil

#### func (ExportSpecifier) Format

```go
func (f ExportSpecifier) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type Expression

```go
type Expression struct {
	Expressions []AssignmentExpression
	Tokens      Tokens
}
```

Expression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-Expression

Expressions must have a length of at least one to be valid.

#### func (Expression) Format

```go
func (f Expression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ForType

```go
type ForType uint8
```

ForType determines which kind of for-loop is described by IterationStatementFor

```go
const (
	ForNormal ForType = iota
	ForNormalVar
	ForNormalLexicalDeclaration
	ForNormalExpression
	ForInLeftHandSide
	ForInVar
	ForInLet
	ForInConst
	ForOfLeftHandSide
	ForOfVar
	ForOfLet
	ForOfConst
	ForAwaitOfLeftHandSide
	ForAwaitOfVar
	ForAwaitOfLet
	ForAwaitOfConst
)
```
Valid ForType's

#### func (ForType) String

```go
func (ft ForType) String() string
```
String implements the fmt.Stringer interface

#### type FormalParameters

```go
type FormalParameters struct {
	FormalParameterList   []BindingElement
	FunctionRestParameter *FunctionRestParameter
	Tokens                Tokens
}
```

FormalParameters as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-FormalParameters

#### func (FormalParameters) Format

```go
func (f FormalParameters) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type FromClause

```go
type FromClause struct {
	ModuleSpecifier *Token
	Tokens          Tokens
}
```

FromClause as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-FromClause

ModuleSpecifier must be non-nil.

#### func (FromClause) Format

```go
func (f FromClause) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type FunctionDeclaration

```go
type FunctionDeclaration struct {
	Type              FunctionType
	BindingIdentifier *Token
	FormalParameters  FormalParameters
	FunctionBody      Block
	Tokens            Tokens
}
```

FunctionDeclaration as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-FunctionDeclaration

Also parses FunctionExpression, for when BindingIdentifier is nil.

Include TC39 proposal for async generator functions
https://github.com/tc39/proposal-async-iteration#async-generator-functions

#### func (FunctionDeclaration) Format

```go
func (f FunctionDeclaration) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type FunctionRestParameter

```go
type FunctionRestParameter struct {
	BindingIdentifier    *Token
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Tokens               Tokens
}
```

FunctionRestParameter as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-FunctionRestParameter

Only one of BindingIdentifier, ArrayBindingPattern, or ObjectBindingPattern must
be non-nil.

#### func (FunctionRestParameter) Format

```go
func (f FunctionRestParameter) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type FunctionType

```go
type FunctionType uint8
```

FunctionType determines which type of function is specified by
FunctionDeclaration

```go
const (
	FunctionNormal FunctionType = iota
	FunctionGenerator
	FunctionAsync
	FunctionAsyncGenerator
)
```
Valid FunctionType's

#### func (FunctionType) String

```go
func (ft FunctionType) String() string
```
Format implements the fmt.Formatter interface

#### type IfStatement

```go
type IfStatement struct {
	Expression    Expression
	Statement     Statement
	ElseStatement *Statement
	Tokens        Tokens
}
```

IfStatement as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-IfStatement

#### func (IfStatement) Format

```go
func (f IfStatement) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ImportClause

```go
type ImportClause struct {
	ImportedDefaultBinding *Token
	NameSpaceImport        *Token
	NamedImports           *NamedImports
	Tokens                 Tokens
}
```

ImportClause as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-ImportClause

At least one of ImportedDefaultBinding, NameSpaceImport, and NamedImports must
be non-nil.

Both NameSpaceImport and NamedImports can not be non-nil.

#### func (ImportClause) Format

```go
func (f ImportClause) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ImportDeclaration

```go
type ImportDeclaration struct {
	*ImportClause
	FromClause
	Tokens Tokens
}
```

ImportDeclaration as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-ImportDeclaration

#### func (ImportDeclaration) Format

```go
func (f ImportDeclaration) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ImportSpecifier

```go
type ImportSpecifier struct {
	IdentifierName  *Token
	ImportedBinding *Token
	Tokens          Tokens
}
```

ImportSpecifier as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-ImportSpecifier

ImportedBinding must be non-nil.

#### func (ImportSpecifier) Format

```go
func (f ImportSpecifier) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type IterationStatementDo

```go
type IterationStatementDo struct {
	Statement  Statement
	Expression Expression
	Tokens     Tokens
}
```

IterationStatementDo is the do-while part of IterationStatement as defined in
ECMA-262 https://262.ecma-international.org/11.0/#prod-IterationStatement

#### func (IterationStatementDo) Format

```go
func (f IterationStatementDo) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type IterationStatementFor

```go
type IterationStatementFor struct {
	Type ForType

	InitExpression *Expression
	InitVar        []VariableDeclaration
	InitLexical    *LexicalDeclaration
	Conditional    *Expression
	Afterthought   *Expression

	LeftHandSideExpression  *LeftHandSideExpression
	ForBindingIdentifier    *Token
	ForBindingPatternObject *ObjectBindingPattern
	ForBindingPatternArray  *ArrayBindingPattern
	In                      *Expression
	Of                      *AssignmentExpression

	Statement Statement
	Tokens    Tokens
}
```

IterationStatementFor is the for part of IterationStatement as defined in
ECMA-262 https://262.ecma-international.org/11.0/#prod-IterationStatement

Includes TC39 proposal for for-await-of
https://github.com/tc39/proposal-async-iteration#the-async-iteration-statement-for-await-of

The Type determines which fields must be non-nil:

    ForInLeftHandSide: LeftHandSideExpression and In
    ForInVar, ForInLet, ForInConst: ForBindingIdentifier, ForBindingPatternObject, or ForBindingPatternArray and In
    ForOfLeftHandSide, ForAwaitOfLeftHandSide: LeftHandSideExpression and Of
    ForOfVar, ForAwaitOfVar, ForOfLet, ForAwaitOfLet, ForOfConst, ForAwaitOfConst: ForBindingIdentifier, ForBindingPatternObject, or ForBindingPatternArray and Of

#### func (IterationStatementFor) Format

```go
func (f IterationStatementFor) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type IterationStatementWhile

```go
type IterationStatementWhile struct {
	Expression Expression
	Statement  Statement
	Tokens     Tokens
}
```

IterationStatementWhile is the while part of IterationStatement as defined in
ECMA-262 https://262.ecma-international.org/11.0/#prod-IterationStatement

#### func (IterationStatementWhile) Format

```go
func (f IterationStatementWhile) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type LeftHandSideExpression

```go
type LeftHandSideExpression struct {
	NewExpression      *NewExpression
	CallExpression     *CallExpression
	OptionalExpression *OptionalExpression
	Tokens             Tokens
}
```

LeftHandSideExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-LeftHandSideExpression

It is only valid for one of NewExpression, CallExpression or OptionalExpression
to be non-nil.

Includes OptionalExpression as per TC39 (2020-03)

#### func (LeftHandSideExpression) Format

```go
func (f LeftHandSideExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type LetOrConst

```go
type LetOrConst bool
```

LetOrConst specifies whether a LexicalDeclaration is a let or const declaration

```go
const (
	Let   LetOrConst = false
	Const LetOrConst = true
)
```
Valid LetOrConst values

#### func (LetOrConst) String

```go
func (l LetOrConst) String() string
```
String implements the fmt.Stringer interface

#### type LexicalBinding

```go
type LexicalBinding struct {
	BindingIdentifier    *Token
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Initializer          *AssignmentExpression
	Tokens               Tokens
}
```

LexicalBinding as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-LexicalBinding

Only one of BindingIdentifier, ArrayBindingPattern or ObjectBindingPattern must
be non-nil. The Initializer is optional only for a BindingIdentifier.

#### func (LexicalBinding) Format

```go
func (f LexicalBinding) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type LexicalDeclaration

```go
type LexicalDeclaration struct {
	LetOrConst
	BindingList []LexicalBinding
	Tokens      Tokens
}
```

LexicalDeclaration as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-LexicalDeclaration

#### func (LexicalDeclaration) Format

```go
func (f LexicalDeclaration) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type LogicalANDExpression

```go
type LogicalANDExpression struct {
	LogicalANDExpression *LogicalANDExpression
	BitwiseORExpression  BitwiseORExpression
	Tokens               Tokens
}
```

LogicalANDExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-LogicalANDExpression

#### func (LogicalANDExpression) Format

```go
func (f LogicalANDExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type LogicalORExpression

```go
type LogicalORExpression struct {
	LogicalORExpression  *LogicalORExpression
	LogicalANDExpression LogicalANDExpression
	Tokens               Tokens
}
```

LogicalORExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-LogicalORExpression

#### func (LogicalORExpression) Format

```go
func (f LogicalORExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type MemberExpression

```go
type MemberExpression struct {
	MemberExpression  *MemberExpression
	PrimaryExpression *PrimaryExpression
	Expression        *Expression
	IdentifierName    *Token
	TemplateLiteral   *TemplateLiteral
	SuperProperty     bool
	NewTarget         bool
	ImportMeta        bool
	Arguments         *Arguments
	Tokens            Tokens
}
```

MemberExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-MemberExpression

If PrimaryExpression is nil, SuperProperty is true, NewTarget is true, or
ImportMeta is true, Expression, IdentifierName, TemplateLiteral, and Arguments
must be nil.

If Expression, IdentifierName, TemplateLiteral, or Arguments is non-nil, then
MemberExpression must be non-nil.

It is only valid if one of Expression, IdentifierName, TemplateLiteral, and
Arguments is non-nil.

#### func (MemberExpression) Format

```go
func (f MemberExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type MethodDefinition

```go
type MethodDefinition struct {
	Type         MethodType
	PropertyName PropertyName
	Params       FormalParameters
	FunctionBody Block
	Tokens       Tokens
}
```

MethodDefinition as specified in ECMA-262
https://262.ecma-international.org/11.0/#prod-MethodDefinition

Static methods from ClassElement are parsed here with the `static` prefix

#### func (MethodDefinition) Format

```go
func (f MethodDefinition) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type MethodType

```go
type MethodType uint8
```

MethodType determines the prefixes for MethodDefinition

```go
const (
	MethodNormal MethodType = iota
	MethodGenerator
	MethodAsync
	MethodAsyncGenerator
	MethodGetter
	MethodSetter
	MethodStatic
	MethodStaticGenerator
	MethodStaticAsync
	MethodStaticAsyncGenerator
	MethodStaticGetter
	MethodStaticSetter
)
```
Valid MethodType's

#### func (MethodType) String

```go
func (mt MethodType) String() string
```
String implements the fmt.Stringer interface

#### type Module

```go
type Module struct {
	ModuleListItems []ModuleItem
	Tokens          Tokens
}
```

Module represents the top-level of a parsed javascript module

#### func  ParseModule

```go
func ParseModule(t parser.Tokeniser) (*Module, error)
```
ParseModule parses a javascript module

#### func (Module) Format

```go
func (f Module) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ModuleItem

```go
type ModuleItem struct {
	ImportDeclaration *ImportDeclaration
	StatementListItem *StatementListItem
	ExportDeclaration *ExportDeclaration
	Tokens            Tokens
}
```

ModuleItem as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-ModuleItem

Only one of ImportDeclaration, StatementListItem, or ExportDeclaration must be
non-nil.

#### func (ModuleItem) Format

```go
func (f ModuleItem) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type MultiplicativeExpression

```go
type MultiplicativeExpression struct {
	MultiplicativeExpression *MultiplicativeExpression
	MultiplicativeOperator   MultiplicativeOperator
	ExponentiationExpression ExponentiationExpression
	Tokens                   Tokens
}
```

MultiplicativeExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-MultiplicativeExpression

If MultiplicativeOperator is not MultiplicativeNone then
MultiplicativeExpression must be non-nil, and vice-versa.

#### func (MultiplicativeExpression) Format

```go
func (f MultiplicativeExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type MultiplicativeOperator

```go
type MultiplicativeOperator int
```

MultiplicativeOperator determines the multication type for
MultiplicativeExpression

```go
const (
	MultiplicativeNone MultiplicativeOperator = iota
	MultiplicativeMultiply
	MultiplicativeDivide
	MultiplicativeRemainder
)
```
Valid MultiplicativeOperator's

#### func (MultiplicativeOperator) String

```go
func (m MultiplicativeOperator) String() string
```
String implements the fmt.Stringer interface

#### type NamedImports

```go
type NamedImports struct {
	ImportList []ImportSpecifier
	Tokens     Tokens
}
```

NamedImports as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-NamedImports

#### func (NamedImports) Format

```go
func (f NamedImports) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type NewExpression

```go
type NewExpression struct {
	News             uint
	MemberExpression MemberExpression
	Tokens           Tokens
}
```

NewExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-NewExpression

The News field is a count of the number of 'new' keywords that proceed the
MemberExpression

#### func (NewExpression) Format

```go
func (f NewExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ObjectBindingPattern

```go
type ObjectBindingPattern struct {
	BindingPropertyList []BindingProperty
	BindingRestProperty *Token
	Tokens              Tokens
}
```

ObjectBindingPattern as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-ObjectBindingPattern

#### func (ObjectBindingPattern) Format

```go
func (f ObjectBindingPattern) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ObjectLiteral

```go
type ObjectLiteral struct {
	PropertyDefinitionList []PropertyDefinition
	Tokens                 Tokens
}
```

ObjectLiteral as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-ObjectLiteral

#### func (ObjectLiteral) Format

```go
func (f ObjectLiteral) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type OptionalChain

```go
type OptionalChain struct {
	OptionalChain   *OptionalChain
	Arguments       *Arguments
	Expression      *Expression
	IdentifierName  *Token
	TemplateLiteral *TemplateLiteral
	Tokens          Tokens
}
```

OptionalChain as defined in TC39
https://tc39.es/ecma262/#prod-OptionalExpression

It is only valid for one of Arguments, Expression, IdentifierName, or
TemplateLiteral to be non-nil.

#### func (OptionalChain) Format

```go
func (f OptionalChain) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type OptionalExpression

```go
type OptionalExpression struct {
	MemberExpression   *MemberExpression
	CallExpression     *CallExpression
	OptionalExpression *OptionalExpression
	OptionalChain      OptionalChain
	Tokens             Tokens
}
```

OptionalExpression as defined in TC39
https://tc39.es/ecma262/#prod-OptionalExpression

It is only valid for one of NewExpression, CallExpression or OptionalExpression
to be non-nil.

#### func (OptionalExpression) Format

```go
func (f OptionalExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type PrimaryExpression

```go
type PrimaryExpression struct {
	This                                              bool
	IdentifierReference                               *Token
	Literal                                           *Token
	ArrayLiteral                                      *ArrayLiteral
	ObjectLiteral                                     *ObjectLiteral
	FunctionExpression                                *FunctionDeclaration
	ClassExpression                                   *ClassDeclaration
	TemplateLiteral                                   *TemplateLiteral
	CoverParenthesizedExpressionAndArrowParameterList *CoverParenthesizedExpressionAndArrowParameterList
	Tokens                                            Tokens
}
```

PrimaryExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-PrimaryExpression

It is only valid is one IdentifierReference, Literal, ArrayLiteral,
ObjectLiteral, FunctionExpression, ClassExpression, TemplateLiteral, or
CoverParenthesizedExpressionAndArrowParameterList is non-nil or This is true.

#### func (PrimaryExpression) Format

```go
func (f PrimaryExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type PropertyDefinition

```go
type PropertyDefinition struct {
	IsCoverInitializedName bool
	PropertyName           *PropertyName
	AssignmentExpression   *AssignmentExpression
	MethodDefinition       *MethodDefinition
	Tokens                 Tokens
}
```

PropertyDefinition as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-PropertyDefinition

One, and only one, of AssignmentExpression or MethodDefinition must be non-nil.

It is only valid for PropertyName to be non-nil when AssignmentExpression is
also non-nil.

The IdentifierReference is stored within PropertyName.

#### func (PropertyDefinition) Format

```go
func (f PropertyDefinition) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type PropertyName

```go
type PropertyName struct {
	LiteralPropertyName  *Token
	ComputedPropertyName *AssignmentExpression
	Tokens               Tokens
}
```

PropertyName as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-PropertyName

Only one of LiteralPropertyName or ComputedPropertyName must be non-nil.

#### func (PropertyName) Format

```go
func (f PropertyName) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type RelationalExpression

```go
type RelationalExpression struct {
	RelationalExpression *RelationalExpression
	RelationshipOperator RelationshipOperator
	ShiftExpression      ShiftExpression
	Tokens               Tokens
}
```

RelationalExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-RelationalExpression

If RelationshipOperator is not RelationshipNone then RelationalExpression must
be non-nil, and vice-verse.

#### func (RelationalExpression) Format

```go
func (f RelationalExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type RelationshipOperator

```go
type RelationshipOperator int
```

RelationshipOperator determines the relationship type for RelationalExpression

```go
const (
	RelationshipNone RelationshipOperator = iota
	RelationshipLessThan
	RelationshipGreaterThan
	RelationshipLessThanEqual
	RelationshipGreaterThanEqual
	RelationshipInstanceOf
	RelationshipIn
)
```
Valid RelationshipOperator's

#### func (RelationshipOperator) String

```go
func (r RelationshipOperator) String() string
```
String implements the fmt.Stringer interface

#### type Script

```go
type Script struct {
	StatementList []StatementListItem
	Tokens        Tokens
}
```

Script represents the top-level of a parsed javascript text

#### func  ParseScript

```go
func ParseScript(t parser.Tokeniser) (*Script, error)
```
ParseScript parses a javascript input into an AST.

It is recommended to use ParseModule instead of this function.

#### func (Script) Format

```go
func (f Script) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ShiftExpression

```go
type ShiftExpression struct {
	ShiftExpression    *ShiftExpression
	ShiftOperator      ShiftOperator
	AdditiveExpression AdditiveExpression
	Tokens             Tokens
}
```

ShiftExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-ShiftExpression

If ShiftOperator is not ShiftNone then ShiftExpression must be non-nil, and
vice-versa.

#### func (ShiftExpression) Format

```go
func (f ShiftExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type ShiftOperator

```go
type ShiftOperator int
```

ShiftOperator determines the shift tyoe for ShiftExpression

```go
const (
	ShiftNone ShiftOperator = iota
	ShiftLeft
	ShiftRight
	ShiftUnsignedRight
)
```
Valid ShiftOperator's

#### func (ShiftOperator) String

```go
func (s ShiftOperator) String() string
```
String implements the fmt.Stringer interface

#### type Statement

```go
type Statement struct {
	Type                    StatementType
	BlockStatement          *Block
	VariableStatement       *VariableStatement
	ExpressionStatement     *Expression
	IfStatement             *IfStatement
	IterationStatementDo    *IterationStatementDo
	IterationStatementWhile *IterationStatementWhile
	IterationStatementFor   *IterationStatementFor
	SwitchStatement         *SwitchStatement
	WithStatement           *WithStatement
	LabelIdentifier         *Token
	LabelledItemFunction    *FunctionDeclaration
	LabelledItemStatement   *Statement
	TryStatement            *TryStatement
	Tokens                  Tokens
}
```

Statement as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-Statement

It is only valid for one of the pointer type to be non-nil.

If LabelIdentifier is non-nil, either one of LabelledItemFunction, or
LabelledItemStatement must be non-nil, or Type must be StatementContinue or
StatementBreak.

If Type is StatementThrow, ExpressionStatement must be non-nil.

#### func (Statement) Format

```go
func (f Statement) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type StatementListItem

```go
type StatementListItem struct {
	Statement   *Statement
	Declaration *Declaration
	Tokens      Tokens
}
```

StatementListItem as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-StatementListItem Only one of
Statement, or Declaration must be non-nil.

#### func (StatementListItem) Format

```go
func (f StatementListItem) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type StatementType

```go
type StatementType uint8
```

StatementType determines the type of a Statement type

```go
const (
	StatementNormal StatementType = iota
	StatementContinue
	StatementBreak
	StatementReturn
	StatementThrow
	StatementDebugger
)
```
Valid StatementType's

#### func (StatementType) String

```go
func (st StatementType) String() string
```
String implements the fmt.Stringer interface

#### type SwitchStatement

```go
type SwitchStatement struct {
	Expression             Expression
	CaseClauses            []CaseClause
	DefaultClause          []StatementListItem
	PostDefaultCaseClauses []CaseClause
	Tokens                 Tokens
}
```

SwitchStatement as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-SwitchStatement

#### func (SwitchStatement) Format

```go
func (f SwitchStatement) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type TemplateLiteral

```go
type TemplateLiteral struct {
	NoSubstitutionTemplate *Token
	TemplateHead           *Token
	Expressions            []Expression
	TemplateMiddleList     []*Token
	TemplateTail           *Token
	Tokens                 Tokens
}
```

TemplateLiteral as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-TemplateLiteral

If NoSubstitutionTemplate is non-nil it is only valid for TemplateHead,
Expressions, TemplateMiddleList, and TemplateTail to be nil.

If NoSubstitutionTemplate is nil, the TemplateHead, Expressions, and
TemplateTail must be non-nil. TemplateMiddleList must have a length of one less
than the length of Expressions.

#### func (TemplateLiteral) Format

```go
func (f TemplateLiteral) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type Token

```go
type Token struct {
	parser.Token
	Pos, Line, LinePos uint64
}
```

Token represents a single parsed token with source positioning

#### func (Token) Format

```go
func (t Token) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type Tokens

```go
type Tokens []Token
```

Tokens is a collection of Token values

#### func (Tokens) Format

```go
func (t Tokens) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type TryStatement

```go
type TryStatement struct {
	TryBlock                           Block
	CatchParameterBindingIdentifier    *Token
	CatchParameterObjectBindingPattern *ObjectBindingPattern
	CatchParameterArrayBindingPattern  *ArrayBindingPattern
	CatchBlock                         *Block
	FinallyBlock                       *Block
	Tokens                             Tokens
}
```

TryStatement as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-TryStatement

Only one of CatchParameterBindingIdentifier, CatchParameterObjectBindingPattern,
and CatchParameterArrayBindingPattern can be non-nil, and must be so if
CatchBlock is non-nil.

If one of CatchParameterBindingIdentifier, CatchParameterObjectBindingPattern,
CatchParameterArrayBindingPattern is non-nil, then CatchBlock must be non-nil.

#### func (TryStatement) Format

```go
func (f TryStatement) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type UnaryExpression

```go
type UnaryExpression struct {
	UnaryOperators   []UnaryOperator
	UpdateExpression UpdateExpression
	Tokens           Tokens
}
```

UnaryExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-UnaryExpression

#### func (UnaryExpression) Format

```go
func (f UnaryExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type UnaryOperator

```go
type UnaryOperator byte
```

UnaryOperator determines a unary operator within UnaryExpression

```go
const (
	UnaryNone UnaryOperator = iota
	UnaryDelete
	UnaryVoid
	UnaryTypeOf
	UnaryAdd
	UnaryMinus
	UnaryBitwiseNot
	UnaryLogicalNot
	UnaryAwait
)
```
Valid UnaryOperator's

#### func (UnaryOperator) String

```go
func (u UnaryOperator) String() string
```
String implements the fmt.Stringer interface

#### type UpdateExpression

```go
type UpdateExpression struct {
	LeftHandSideExpression *LeftHandSideExpression
	UpdateOperator         UpdateOperator
	UnaryExpression        *UnaryExpression
	Tokens                 Tokens
}
```

UpdateExpression as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-UpdateExpression

If UpdateOperator is UpdatePreIncrement or UpdatePreDecrement UnaryExpression
must be non-nil, and vice-versa. In all other cases, LeftHandSideExpression must
be non-nil.

#### func (UpdateExpression) Format

```go
func (f UpdateExpression) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type UpdateOperator

```go
type UpdateOperator int
```

UpdateOperator determines the type of update operation for UpdateExpression

```go
const (
	UpdateNone UpdateOperator = iota
	UpdatePostIncrement
	UpdatePostDecrement
	UpdatePreIncrement
	UpdatePreDecrement
)
```
Valid UpdateOperator's

#### func (UpdateOperator) String

```go
func (u UpdateOperator) String() string
```
String implements the fmt.Stringer interface

#### type VariableDeclaration

```go
type VariableDeclaration LexicalBinding
```

VariableDeclaration as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-VariableDeclaration

#### func (VariableDeclaration) Format

```go
func (f VariableDeclaration) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type VariableStatement

```go
type VariableStatement struct {
	VariableDeclarationList []VariableDeclaration
	Tokens                  Tokens
}
```

VariableStatement as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-VariableStatement

VariableDeclarationList must have a length or at least one.

#### func (VariableStatement) Format

```go
func (f VariableStatement) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface

#### type WithStatement

```go
type WithStatement struct {
	Expression Expression
	Statement  Statement
	Tokens     Tokens
}
```

WithStatement as defined in ECMA-262
https://262.ecma-international.org/11.0/#prod-WithStatement

#### func (WithStatement) Format

```go
func (f WithStatement) Format(s fmt.State, v rune)
```
Format implements the fmt.Formatter interface
