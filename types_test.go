package javascript

import (
	"reflect"
	"testing"
)

func TestTypesData(t *testing.T) {
	for n, test := range [...]struct {
		Token
		Data interface{}
	}{
		{
			String("\"Hello\""),
			"Hello",
		},
		{
			String("'Hello,\n\"World\"'"),
			"Hello,\n\"World\"",
		},
		{
			String("\"\\x48\\u0065llo,\tWorld\""),
			"Hello,	World",
		},
	} {
		data := test.Token.Data()
		if !reflect.DeepEqual(test.Data, data) {
			t.Errorf("test %d: expecting %v, got %v", n+1, test.Data, data)
		}
	}
}
