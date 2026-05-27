package javascript_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/walk"
	"vimagination.zapto.org/parser"
)

func TestPrintingOriginal(t *testing.T) {
	for n, test := range [...]struct {
		Input, Output string
	}{} {
		tk := parser.NewStringTokeniser(test.Input)

		m, err := javascript.ParseModule(&tk)
		if err != nil {
			t.Errorf("test %d.1: unexpected error: %s", n+1, err)

			continue
		}

		var h walk.Handler

		h = walk.HandlerFunc(func(t javascript.Type) error {
			v := reflect.ValueOf(t)

			if v.Type().Kind() != reflect.Pointer || v.Type().Elem().Kind() != reflect.Struct {
				return nil
			}

			if f, ok := v.Type().Elem().FieldByName("Tokens"); ok {
				v.Elem().FieldByIndex(f.Index).SetZero()
			}

			return walk.Walk(t, h)
		})

		if err = h.Handle(m); err != nil {
			t.Errorf("test %d.2: unexpected error: %s", n+1, err)

			continue
		}

		var buf bytes.Buffer

		fmt.Fprintf(&buf, "%#s", m)

		if str := buf.String(); str != test.Output {
			t.Errorf("test %d.3: expecting %q, got %q", n+1, test.Output, str)
		}
	}
}
