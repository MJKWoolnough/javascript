package minify

import (
	"sort"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/scope"
)

type binding struct {
	Name    string
	Scope   *scope.Scope
	NameSet bool
}

func orderedScope(s *scope.Scope) []binding {
	b := walkScope(s, nil)
	sort.Slice(b, func(i, j int) bool {
		if b[i].NameSet {
			if !b[j].NameSet {
				return false
			}
		} else if b[j].NameSet {
			return true
		}
		il := len(b[i].Scope.Bindings[b[i].Name])
		jl := len(b[j].Scope.Bindings[b[j].Name])

		if il == jl {
			return b[i].Name < b[j].Name
		}
		return il > jl
	})
	return b
}

func walkScope(s *scope.Scope, b []binding) []binding {
	for name := range s.Bindings {
		if name == "this" || name == "arguments" {
			continue
		}
		b = append(b, binding{
			Name:    name,
			Scope:   s,
			NameSet: s.Bindings[name][0].BindingType == scope.BindingRef,
		})
	}
	for _, cs := range s.Scopes {
		b = walkScope(cs, b)
	}
	return b
}

func renameIdentifiers(m *javascript.Module) error {
	s, err := scope.ModuleScope(m, nil)
	if err != nil {
		return err
	}
	bindings := orderedScope(s)
	for n, binding := range bindings {
		if binding.NameSet {
			break
		}
		identifiersInScope := make(map[string]struct{})
		for _, checkBinding := range bindings {
			if !checkBinding.NameSet || checkBinding == binding {
				continue
			}
			if binding.Scope == checkBinding.Scope {
				identifiersInScope[checkBinding.Name] = struct{}{}
			} else if isParentScope(binding.Scope, checkBinding.Scope) {
				for _, scope := range binding.Scope.Bindings[binding.Name] {
					if isParentScope(checkBinding.Scope, scope.Scope) {
						identifiersInScope[checkBinding.Name] = struct{}{}
						break
					}
				}
			} else {
				for _, scope := range checkBinding.Scope.Bindings[checkBinding.Name] {
					if isParentScope(binding.Scope, scope.Scope) {
						identifiersInScope[checkBinding.Name] = struct{}{}
						break
					}
				}
			}
		}
		name := makeUniqueName(identifiersInScope)
		for _, b := range binding.Scope.Bindings[binding.Name] {
			b.Data = name
		}
		bindings[n].Name = name
		bindings[n].NameSet = true
	}
	return nil
}

func isParentScope(a, b *scope.Scope) bool {
	for b != nil {
		if b == a {
			return true
		}
		b = b.Parent
	}
	return false
}

var (
	extraChars = []byte("0123456789_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ$")
	startChars = extraChars[10:]
)

func makeUniqueName(exclude map[string]struct{}) string {
	name := make([]byte, 8)
	parts := make([][]byte, 8)
	for l := 0; ; l++ {
		parts = parts[:1]
		parts[0] = startChars
		for n := 0; n < l; n++ {
			parts = append(parts, extraChars)
		}
	L:
		for {
			name = name[:0]
			for i := 0; i <= l; i++ {
				name = append(name, parts[i][0])
			}
			if _, ok := exclude[string(name)]; !ok {
				return string(name)
			}
			for i := l; i >= 0; i-- {
				parts[i] = parts[i][1:]
				if len(parts[i]) > 0 {
					break
				}
				if i == 0 {
					break L
				}
				parts[i] = extraChars
			}
		}
	}
}
