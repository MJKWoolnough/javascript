package minify

import (
	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/scope"
)

func removeDeadCode(m *javascript.Module) {
	for {
		s, err := scope.ModuleScope(m, nil)
		if err != nil {
			return
		}

		clearSinglesFromScope(s)

	}
}

func clearSinglesFromScope(s *scope.Scope) bool {
	changed := false
	for name, bindings := range s.Bindings {
		if name == "this" || name == "arguments" || len(bindings) != 1 {
			continue
		}
		bindings[0].Token.Data = ""
		changed = true
	}

	for _, cs := range s.Scopes {
		if clearSinglesFromScope(cs) {
			changed = true
		}
	}

	return changed
}
