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
