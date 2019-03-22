package javascript

type ClassDeclaration struct {
	Tokens []TokenPos
}

func (j *jsParser) parseClassDeclaration(yield, await, def bool) (ClassDeclaration, error) {
	var cd ClassDeclaration
	return cd, nil
}
