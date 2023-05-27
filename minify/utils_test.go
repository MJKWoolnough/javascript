package minify

import (
	"reflect"
	"testing"

	"vimagination.zapto.org/javascript"
)

func TestBlockAsModule(t *testing.T) {
	for n, test := range [...]struct {
		Input    *javascript.Block
		Callback func(*javascript.Module)
		Output   *javascript.Block
	}{
		{ // 1
			&javascript.Block{},
			func(m *javascript.Module) {
			},
			&javascript.Block{},
		},
		{ // 2
			&javascript.Block{
				StatementList: []javascript.StatementListItem{
					{
						Statement: &javascript.Statement{
							Type: javascript.StatementDebugger,
						},
					},
					{
						Statement: &javascript.Statement{
							Type: javascript.StatementContinue,
						},
					},
					{
						Statement: &javascript.Statement{
							Type: javascript.StatementDebugger,
						},
					},
					{
						Statement: &javascript.Statement{
							Type: javascript.StatementContinue,
						},
					},
				},
			},
			func(m *javascript.Module) {
				for i := 0; i < len(m.ModuleListItems); i++ {
					if m.ModuleListItems[i].StatementListItem.Statement.Type == javascript.StatementDebugger {
						m.ModuleListItems = append(m.ModuleListItems[:i], m.ModuleListItems[i+1:]...)
						i--
					}
				}
			},
			&javascript.Block{
				StatementList: []javascript.StatementListItem{
					{
						Statement: &javascript.Statement{
							Type: javascript.StatementContinue,
						},
					},
					{
						Statement: &javascript.Statement{
							Type: javascript.StatementContinue,
						},
					},
				},
			},
		},
	} {
		blockAsModule(test.Input, test.Callback)
		if !reflect.DeepEqual(test.Input, test.Output) {
			t.Errorf("test %d: expecting output %v, got %v", n+1, test.Output, test.Input)
		}
	}
}
