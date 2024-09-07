# scope
--
    import "vimagination.zapto.org/javascript/scope"

Package scope parses out a scope tree for a javascript module or script.

## Usage

#### type Binding

```go
type Binding struct {
	BindingType
	*Scope
	*javascript.Token
}
```

Binding represents a single instance of a bound name.

#### type BindingType

```go
type BindingType uint8
```

BindingType indicates where the binding came from.

```go
const (
	BindingRef BindingType = iota
	BindingBare
	BindingVar
	BindingHoistable
	BindingLexicalLet
	BindingLexicalConst
	BindingImport
	BindingFunctionParam
	BindingCatch
)
```
Binding Types.

#### type ErrDuplicateDeclaration

```go
type ErrDuplicateDeclaration struct {
	Declaration, Duplicate *javascript.Token
}
```

ErrDuplicateDeclaration is an error when a binding is declared more than once
with a scope.

#### func (ErrDuplicateDeclaration) Error

```go
func (ErrDuplicateDeclaration) Error() string
```

#### type Scope

```go
type Scope struct {
	IsLexicalScope bool
	Parent         *Scope
	Scopes         map[javascript.Type]*Scope
	Bindings       map[string][]Binding
}
```

Scope represents a single level of variable scope.

#### func  ModuleScope

```go
func ModuleScope(m *javascript.Module, global *Scope) (*Scope, error)
```
ModuleScope parses out the scope tree for a javascript Module

#### func  NewScope

```go
func NewScope() *Scope
```
NewScope returns a init'd Scope type.

#### func  ScriptScope

```go
func ScriptScope(s *javascript.Script, global *Scope) (*Scope, error)
```
ScriptScope parses out the scope tree for a javascript script
