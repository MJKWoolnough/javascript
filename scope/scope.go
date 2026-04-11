// Package scope parses out a scope tree for a JavaScript module or script.
package scope // import "vimagination.zapto.org/javascript/scope"

import (
	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/walk"
)

// ErrDuplicateDeclaration is an error when a binding is declared more than once with a scope.
type ErrDuplicateDeclaration struct {
	Declaration, Duplicate *javascript.Token
}

func (ErrDuplicateDeclaration) Error() string {
	return "duplicate declaration"
}

// BindingType indicates where the binding came from.
type BindingType uint8

// Binding Types.
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

// Binding represents a single instance of a bound name.
type Binding struct {
	BindingType
	*Scope
	*javascript.Token
}

// Scope represents a single level of variable scope.
type Scope struct {
	IsLexicalScope bool
	Parent         *Scope
	Scopes         map[javascript.Type]*Scope
	Bindings       map[string][]Binding
}

func (s *Scope) setBinding(t *javascript.Token, bindingType BindingType) error {
	name := t.Data
	binding := Binding{BindingType: bindingType, Token: t, Scope: s}

	if b, ok := s.Bindings[name]; ok {
		if bindingType == BindingVar && len(b) > 0 && (b[0].BindingType == BindingVar || b[0].BindingType == BindingCatch) {
			s.Bindings[name] = append(b, binding)

			if b[0].BindingType == BindingCatch && bindingType == BindingVar {
				return nil
			}
		} else {
			var bd *javascript.Token

			if len(b) > 0 {
				bd = b[0].Token
			}

			return ErrDuplicateDeclaration{
				Declaration: bd,
				Duplicate:   t,
			}
		}
	} else {
		s.Bindings[name] = []Binding{binding}
	}

	if s.IsLexicalScope && (bindingType == BindingHoistable || bindingType == BindingVar) {
	Loop:
		for s.IsLexicalScope && s.Parent != nil {
			s = s.Parent

			if bindingType == BindingVar {
				if b, ok := s.Bindings[name]; ok && len(b) > 0 {
					switch b[0].BindingType {
					case BindingCatch:
						break Loop
					case BindingVar, BindingBare:
					default:
						return ErrDuplicateDeclaration{
							Declaration: b[0].Token,
							Duplicate:   t,
						}
					}
				}
			}
		}

		if b, ok := s.Bindings[name]; !ok {
			s.Bindings[name] = []Binding{binding}
		} else if bindingType == BindingVar {
			s.Bindings[name] = append(b, binding)
		}
	}

	return nil
}

func (s *Scope) addBinding(t *javascript.Token, bindingType BindingType) {
	name := t.Data
	binding := Binding{BindingType: bindingType, Token: t, Scope: s}

	for {
		if bs, ok := s.Bindings[name]; ok {
			s.Bindings[name] = append(bs, binding)

			if !s.IsLexicalScope || len(bs) == 0 || bs[0].BindingType != BindingVar {
				return
			}
		}

		if s.Parent == nil {
			s.Bindings[name] = []Binding{binding}

			return
		}

		s = s.Parent
	}
}

// NewScope returns a init'd Scope type.
func NewScope() *Scope {
	return &Scope{
		Scopes:   make(map[javascript.Type]*Scope),
		Bindings: make(map[string][]Binding),
	}
}

func (s *Scope) newFunctionScope(js javascript.Type) *Scope {
	if ns, ok := s.Scopes[js]; ok {
		return ns
	}

	ns := &Scope{
		Parent: s,
		Scopes: make(map[javascript.Type]*Scope),
		Bindings: map[string][]Binding{
			"this":      {},
			"arguments": {},
		},
	}

	s.Scopes[js] = ns

	return ns
}

func (s *Scope) newArrowFunctionScope(js javascript.Type) *Scope {
	if ns, ok := s.Scopes[js]; ok {
		return ns
	}

	ns := &Scope{
		Parent:   s,
		Scopes:   make(map[javascript.Type]*Scope),
		Bindings: make(map[string][]Binding),
	}

	s.Scopes[js] = ns

	return ns
}

func (s *Scope) newLexicalScope(js javascript.Type) *Scope {
	if ns, ok := s.Scopes[js]; ok {
		return ns
	}

	ns := &Scope{
		Parent:         s,
		IsLexicalScope: true,
		Scopes:         make(map[javascript.Type]*Scope),
		Bindings:       make(map[string][]Binding),
	}

	s.Scopes[js] = ns

	return ns
}

// ModuleScope parses out the scope tree for a JavaScript Module
func ModuleScope(m *javascript.Module, global *Scope) (*Scope, error) {
	if global == nil {
		global = NewScope()
	}

	if err := walk.Walk(m, &scoper{
		scope: global,
		set:   true,
	}); err != nil {
		return nil, err
	}

	walk.Walk(m, &scoper{scope: global})

	return global, nil
}

// ScriptScope parses out the scope tree for a JavaScript script
func ScriptScope(s *javascript.Script, global *Scope) (*Scope, error) {
	if global == nil {
		global = NewScope()
	}

	if err := walk.Walk(s, &scoper{
		scope: global,
		set:   true,
	}); err != nil {
		return nil, err
	}

	walk.Walk(s, &scoper{scope: global})

	return global, nil
}

// FindIdentifier look up the Scope chain to find the first that contains the
// specified identifier.
//
// If the current scope, and no parent, contains the identifier, nil is
// returned.
func (s *Scope) FindIdentifier(name string) *Scope {
	if s == nil {
		return nil
	} else if _, ok := s.Bindings[name]; !ok {
		return s.Parent.FindIdentifier(name)
	}

	return s
}

// Rename will rename an identifier in the current scope, returning true on a
// success.
//
// If the new identifier is already declared within the scope, this function
// will do nothing and return false.
//
// It is recommended to check whether renaming the identifier will break child
// scopes by using IdentifierInUse.
func (s *Scope) Rename(from, to string) bool {
	if b, ok := s.Bindings[to]; ok && b[0].BindingType != BindingBare {
		return false
	}

	for _, b := range s.Bindings[from] {
		b.Data = to
	}

	s.Bindings[to] = append(s.Bindings[to], s.Bindings[from]...)
	delete(s.Bindings, from)

	return true
}

// IdentifierInUse returns true if the given identifier is declared in this, or
// a parent scope, and is used in this or a child scope.
//
// Can be used to determine whether renaming, or adding an identifier
// declaration will break a child scope.
func (s *Scope) IdentifierInUse(identifier string) bool {
	for t := s; t != nil; t = t.Parent {
		for _, binding := range t.Bindings[identifier] {
			for p := binding.Scope; p != nil; p = p.Parent {
				if p == s {
					return true
				}
			}
		}
	}

	return false
}
