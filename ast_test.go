package javascript

import (
	"reflect"
	"testing"

	"vimagination.zapto.org/parser"
)

type sourceFn struct {
	Source string
	Fn     func(*test, Tokens)
}

type test struct {
	Tokens                jsParser
	Yield, Await, In, Def bool
	Output                interface{}
	Err                   error
}

func doTests(t *testing.T, tests []sourceFn, fn func(*test) (interface{}, error)) {
	t.Helper()
	var err error
	for n, tt := range tests {
		var ts test
		ts.Tokens, err = newJSParser(parser.NewStringTokeniser(tt.Source))
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
			continue
		}
		tt.Fn(&ts, Tokens(ts.Tokens[:cap(ts.Tokens)]))
		output, err := fn(&ts)
		if !reflect.DeepEqual(err, ts.Err) {
			t.Errorf("test %d: expecting error: %v, got %v", n+1, ts.Err, err)
		} else if ts.Output != nil && !reflect.DeepEqual(output, ts.Output) {
			t.Errorf("test %d: expecting \n%+v\n...got...\n%+v", n+1, ts.Output, output)
		}
	}
}

func TestIdentifier(t *testing.T) {
	doTests(t, []sourceFn{
		{`hello_world`, func(t *test, tk Tokens) {
			t.Output = Identifier{&tk[0]}
		}},
		{`yield`, func(t *test, tk Tokens) {
			t.Output = Identifier{&tk[0]}
		}},
		{`yield`, func(t *test, tk Tokens) {
			t.Await = true
			t.Output = Identifier{&tk[0]}
		}},
		{`yield`, func(t *test, tk Tokens) {
			t.Yield = true
			t.Err = Error{
				Err:     ErrMissingIdentifier,
				Parsing: "Identifier",
				Token:   tk[0],
			}
		}},
		{`await`, func(t *test, tk Tokens) {
			t.Output = Identifier{&tk[0]}
		}},
		{`await`, func(t *test, tk Tokens) {
			t.Yield = true
			t.Output = Identifier{&tk[0]}
		}},
		{`await`, func(t *test, tk Tokens) {
			t.Await = true
			t.Err = Error{
				Err:     ErrMissingIdentifier,
				Parsing: "Identifier",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (interface{}, error) {
		return t.Tokens.parseIdentifier(t.Yield, t.Await)
	})
}
