# javascript
--
    import "vimagination.zapto.org/javascript"

Package javascript provides tools to tokenise and parse javascript source files

## Usage

```go
const (
	TokenWhitespace parser.TokenType = iota
	TokenLineTerminator
	TokenSingleLineComment
	TokenMultiLineComment
	TokenIdentifier
	TokenBooleanLiteral
	TokenKeyword
	TokenPunctuator
	TokenNumericLiteral
	TokenStringLiteral
	TokenNoSubstitutionTemplate
	TokenTemplateHead
	TokenTemplateMiddle
	TokenTemplateTail
	TokenDivPunctuator
	TokenRightBracePunctuator
	TokenRegularExpressionLiteral
)
```
Javascript Token values

#### func  SetTokeniser

```go
func SetTokeniser(t *parser.Tokeniser) *parser.Tokeniser
```
SetTokeniser provides javascript parsing functions to a Tokeniser

#### type Boolean

```go
type Boolean bool
```

Boolean represents a literal boolean value (true, false)

#### func (Boolean) Data

```go
func (tk Boolean) Data() interface{}
```
Data returns typed interpreted data

#### func (Boolean) WriteTo

```go
func (tk Boolean) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type Identifier

```go
type Identifier string
```

Identifier represents a non-keyword identified, usually a variable name or
global

#### func (Identifier) Data

```go
func (tk Identifier) Data() interface{}
```
Data returns typed interpreted data

#### func (Identifier) String

```go
func (tk Identifier) String() string
```
String returns the unescaped string

#### func (Identifier) WriteTo

```go
func (tk Identifier) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type Keyword

```go
type Keyword string
```

Keyword represents one of the following reserved keywords:

await, break, case, catch, class, const, continue, debugger, default, delete,
do, else, export, extends, finally, for, function, if, import, in, instanceof,
new, return, super, switch, this, throw, try, typeof, var, void, while, with,
yield

#### func (Keyword) Data

```go
func (tk Keyword) Data() interface{}
```
Data returns typed interpreted data

#### func (Keyword) WriteTo

```go
func (tk Keyword) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type LineTerminators

```go
type LineTerminators string
```

LineTerminators is sequential line feeds, carriage returns, line separators and
paragraph separators

#### func (LineTerminators) Data

```go
func (tk LineTerminators) Data() interface{}
```
Data returns typed interpreted data

#### func (LineTerminators) WriteTo

```go
func (tk LineTerminators) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type MultiLineComment

```go
type MultiLineComment string
```

MultiLineComment is a comment started by a slash and an asterix

#### func (MultiLineComment) Comment

```go
func (tk MultiLineComment) Comment() string
```
Comment returns the contents of the comment

#### func (MultiLineComment) Data

```go
func (tk MultiLineComment) Data() interface{}
```
Data returns typed interpreted data

#### func (MultiLineComment) WriteTo

```go
func (tk MultiLineComment) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type NoSubstitutionTemplate

```go
type NoSubstitutionTemplate string
```

NoSubstitutionTemplate represents a template literal

#### func (NoSubstitutionTemplate) Data

```go
func (tk NoSubstitutionTemplate) Data() interface{}
```
Data returns typed interpreted data

#### func (NoSubstitutionTemplate) String

```go
func (tk NoSubstitutionTemplate) String() string
```
String returns the unescaped string

#### func (NoSubstitutionTemplate) WriteTo

```go
func (tk NoSubstitutionTemplate) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type Number

```go
type Number string
```

Number represents a number literal

#### func (Number) Data

```go
func (tk Number) Data() interface{}
```
Data returns typed interpreted data

#### func (Number) Number

```go
func (tk Number) Number() float64
```
Number returns a a numeric interpretation of the data

#### func (Number) WriteTo

```go
func (tk Number) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type NumberBinary

```go
type NumberBinary string
```

NumberBinary represents a binary literal

#### func (NumberBinary) Data

```go
func (tk NumberBinary) Data() interface{}
```
Data returns typed interpreted data

#### func (NumberBinary) Number

```go
func (tk NumberBinary) Number() float64
```
Number returns a a numeric interpretation of the data

#### func (NumberBinary) WriteTo

```go
func (tk NumberBinary) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type NumberHexadecimal

```go
type NumberHexadecimal string
```

NumberHexadecimal represents a hexadecimal literal

#### func (NumberHexadecimal) Data

```go
func (tk NumberHexadecimal) Data() interface{}
```
Data returns typed interpreted data

#### func (NumberHexadecimal) Number

```go
func (tk NumberHexadecimal) Number() float64
```
Number returns a a numeric interpretation of the data

#### func (NumberHexadecimal) WriteTo

```go
func (tk NumberHexadecimal) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type NumberOctal

```go
type NumberOctal string
```

NumberOctal represents an octal literal

#### func (NumberOctal) Data

```go
func (tk NumberOctal) Data() interface{}
```
Data returns typed interpreted data

#### func (NumberOctal) Number

```go
func (tk NumberOctal) Number() float64
```
Number returns a a numeric interpretation of the data

#### func (NumberOctal) WriteTo

```go
func (tk NumberOctal) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type Punctuator

```go
type Punctuator string
```

Punctuator represents one of the following:

{, }, [, ], (, ), ;, ,, ?, :, ~, ., <, >, >=, <=, >>=, <<=, >>>, =, ==, ===,
==>, !, !=, !==, +, +=, -, -=, *, **, *=, /, /= &, &&, |, ||, &=, |=, %, %=, ^,
^=

#### func (Punctuator) Data

```go
func (tk Punctuator) Data() interface{}
```
Data returns typed interpreted data

#### func (Punctuator) WriteTo

```go
func (tk Punctuator) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type Regex

```go
type Regex string
```

Regex represents a regular expression literal

#### func (Regex) Data

```go
func (tk Regex) Data() interface{}
```
Data returns typed interpreted data

#### func (Regex) WriteTo

```go
func (tk Regex) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type SingleLineComment

```go
type SingleLineComment string
```

SingleLineComment is a comment started by two slashes

#### func (SingleLineComment) Comment

```go
func (tk SingleLineComment) Comment() string
```
Comment returns the contents of the comment

#### func (SingleLineComment) Data

```go
func (tk SingleLineComment) Data() interface{}
```
Data returns typed interpreted data

#### func (SingleLineComment) WriteTo

```go
func (tk SingleLineComment) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type String

```go
type String string
```

String represents a string literal

#### func (String) Data

```go
func (tk String) Data() interface{}
```
Data returns typed interpreted data

#### func (String) String

```go
func (tk String) String() string
```
String returns the unescaped string

#### func (String) WriteTo

```go
func (tk String) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type Template

```go
type Template struct {
	TemplateStart
	TemplateMiddle Tokens
	TemplateEnd
}
```

Template represents a template with substitutions

#### func (Template) Data

```go
func (tk Template) Data() interface{}
```
Data returns typed interpreted data

#### func (Template) WriteTo

```go
func (tk Template) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type TemplateEnd

```go
type TemplateEnd string
```

TemplateEnd represents the end of a template

#### func (TemplateEnd) Data

```go
func (tk TemplateEnd) Data() interface{}
```
Data returns typed interpreted data

#### func (TemplateEnd) String

```go
func (tk TemplateEnd) String() string
```
String returns the unescaped string

#### func (TemplateEnd) WriteTo

```go
func (tk TemplateEnd) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type TemplateMiddle

```go
type TemplateMiddle string
```

TemplateMiddle represents a middle chunk of a template

#### func (TemplateMiddle) Data

```go
func (tk TemplateMiddle) Data() interface{}
```
Data returns typed interpreted data

#### func (TemplateMiddle) String

```go
func (tk TemplateMiddle) String() string
```
String returns the unescaped string

#### func (TemplateMiddle) WriteTo

```go
func (tk TemplateMiddle) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type TemplateStart

```go
type TemplateStart string
```

TemplateStart represents the opening of a template

#### func (TemplateStart) Data

```go
func (tk TemplateStart) Data() interface{}
```
Data returns typed interpreted data

#### func (TemplateStart) String

```go
func (tk TemplateStart) String() string
```
String returns the unescaped string

#### func (TemplateStart) WriteTo

```go
func (tk TemplateStart) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token

#### type Token

```go
type Token interface {
	io.WriterTo
	Data() interface{}
}
```

Token represents a parsed Token in the tree

#### type Tokens

```go
type Tokens []Token
```

Tokens represents a list of Tokens

#### func  Tree

```go
func Tree(t parser.Tokeniser) (Tokens, error)
```
Tree uses the given Tokeniser to produce a tree of tokens

#### func (Tokens) Data

```go
func (tk Tokens) Data() interface{}
```
Data returns typed interpreted data

#### func (Tokens) WriteTo

```go
func (tk Tokens) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the tokens

#### type Whitespace

```go
type Whitespace string
```

Whitespace is sequential tabs, vertical tabs, form feeds, spaces, no-break
spaces and zero width no-break spaces

#### func (Whitespace) Data

```go
func (tk Whitespace) Data() interface{}
```
Data returns typed interpreted data

#### func (Whitespace) WriteTo

```go
func (tk Whitespace) WriteTo(w io.Writer) (int64, error)
```
WriteTo implements the io.WriterTo interface and writes to the Writer the
original contents of the token
