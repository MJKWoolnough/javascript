package javascript

const marker = "TS"

func ParseTypescript(t Tokeniser) (*Script, error) {
	j, err := newJSParser(t)
	if err != nil {
		return nil, err
	}
	j[len(j)-1].Data = marker
	s := new(Script)
	if err := s.parse(&j); err != nil {
		return nil, err
	}
	return s, nil
}

func ParseTypescriptModule(t Tokeniser) (*Module, error) {
	j, err := newJSParser(t)
	if err != nil {
		return nil, err
	}
	j[:cap(j)][cap(j)-1].Data = marker
	m := new(Module)
	if err := m.parse(&j); err != nil {
		return nil, err
	}
	return m, nil
}

func (j *jsParser) IsTypescript() bool {
	return (*j)[:cap(*j)][cap(*j)-1].Data == marker
}

/*
ClassDeclaration (<>, implements)
AssignmentExpression (!, as)
FieldDefinition (private, protected)
MethodDefinition (<>)
FormalParameters (:TYPE)
FunctionDeclaration (<>, :TYPE)
ArrowFunction (<>, :TYPE)
StatementListItem (enum, type, interface)
LexicalBinding (!:TYPE)
TryStatement (:TYPE)
ModuleItem (import type)
LeftHandSideExpression (<>)
*/

func (j *jsParser) SkipGeneric() {}

func (j *jsParser) SkipAsType() {}

func (j *jsParser) SkipColonType() {}

func (j *jsParser) SkipType() {}

func (j *jsParser) SkipInterface() {}

func (j *jsParser) SkipEnum() {}

func (j *jsParser) SkipImportType() {
}
