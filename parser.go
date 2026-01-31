package main

type parser struct {
	tokens []Token
	pos    int

	isBoldOpen        bool
	isItalicOpen      bool
	isInlineQuoteOpen bool
	isMathOpen        bool

	isParagraph bool
	isList      bool
}
