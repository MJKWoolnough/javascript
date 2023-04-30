package minify

import (
	"sort"

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
		return len(b[i].Scope.Bindings[b[i].Name]) > len(b[j].Scope.Bindings[b[j].Name])
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
