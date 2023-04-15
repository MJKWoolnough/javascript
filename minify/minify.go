package minify

type Minifier struct{}

func New(opts ...Option) *Minifier {
	m := new(Minifier)
	for _, opt := range opts {
		opt(m)
	}
	return m
}
