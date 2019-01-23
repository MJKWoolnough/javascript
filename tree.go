// Package javascript provides tools to tokenise and parse javascript source files
package javascript

import (
	"math"
	"strings"

	"vimagination.zapto.org/parser"
)

var (
	pInf = math.Inf(1)
)

func Tree(t parser.Tokeniser) (Tokens, error) {
	var (
		j       jsParser
		tree    = make([]Tokens, 1, 32)
		treeLen = 0
	)
	t.TokeniserState(j.inputElement)
	for {
		tk, err := t.GetToken()
		if err != nil {
			return nil, err
		}
		switch tk.Type {
		case parser.TokenDone:
			return tree[0], nil
		case TokenWhitespace:
			tree[treeLen] = append(tree[treeLen], Whitespace(tk.Data))
		case TokenLineTerminator:
			tree[treeLen] = append(tree[treeLen], LineTerminators(tk.Data))
		case TokenSingleLineComment:
			tree[treeLen] = append(tree[treeLen], SingleLineComment(tk.Data))
		case TokenMultiLineComment:
			tree[treeLen] = append(tree[treeLen], MultiLineComment(tk.Data))
		case TokenIdentifier:
			tree[treeLen] = append(tree[treeLen], Identifier(tk.Data))
		case TokenBooleanLiteral:
			if tk.Data == "true" {
				tree[treeLen] = append(tree[treeLen], Boolean(true))
			} else {
				tree[treeLen] = append(tree[treeLen], Boolean(false))
			}
		case TokenKeyword:
			tree[treeLen] = append(tree[treeLen], Keyword(tk.Data))
		case TokenPunctuator:
			tree[treeLen] = append(tree[treeLen], Punctuator(tk.Data))
			tree = append(tree, nil)
			treeLen++
		case TokenNumericLiteral:
			if strings.HasPrefix(tk.Data, "0b") || strings.HasPrefix(tk.Data, "0B") {
				tree[treeLen] = append(tree[treeLen], NumberBinary(tk.Data))
			} else if strings.HasPrefix(tk.Data, "0o") || strings.HasPrefix(tk.Data, "0O") {
				tree[treeLen] = append(tree[treeLen], NumberOctal(tk.Data))
			} else if strings.HasPrefix(tk.Data, "0x") || strings.HasPrefix(tk.Data, "0X") {
				tree[treeLen] = append(tree[treeLen], NumberHexadecimal(tk.Data))
			} else {
				tree[treeLen] = append(tree[treeLen], Number(tk.Data))
			}
		case TokenStringLiteral:
			tree[treeLen] = append(tree[treeLen], String(tk.Data))
		case TokenNoSubstitutionTemplate:
			tree[treeLen] = append(tree[treeLen], NoSubstitutionTemplate(tk.Data))
		case TokenTemplateHead:
			tree = append(tree, nil)
			treeLen++
			tree[treeLen] = append(tree[treeLen], TemplateStart(tk.Data))
		case TokenTemplateMiddle:
			tree[treeLen] = append(tree[treeLen], TemplateMiddle(tk.Data))
		case TokenTemplateTail:
			tree[treeLen] = append(tree[treeLen], TemplateEnd(tk.Data))
			treeLen--
			tree[treeLen] = append(tree[treeLen], Template(tree[treeLen+1]))
			tree = tree[:treeLen+1]
		case TokenDivPunctuator, TokenRightBracePunctuator:
			tree[treeLen-1] = append(tree[treeLen-1], tree[treeLen])
			tree = tree[:treeLen]
			treeLen--
			tree[treeLen] = append(tree[treeLen], Punctuator(tk.Data))
		case TokenRegularExpressionLiteral:
			tree[treeLen] = append(tree[treeLen], Regex(tk.Data))
		}
	}
}
