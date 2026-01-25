package main

type TokenType int

const (
	EOF TokenType = iota
	TEXT
	ITALIC
	BOLD
	INLINE_CODE
	QUOTE
)

func (t TokenType) String() string {
	switch t {
	case TEXT:
		return "TEXT"
	case ITALIC:
		return "ITALIC"
	case BOLD:
		return "BOLD"
	case INLINE_CODE:
		return "INLINE_CODE"
	case QUOTE:
		return "QUOTE"
	default:
		return "None"
	}
}

type State int

const (
	StateText State = iota
	StateBold
	SateItalic
	StateInLineCode
	StateQuote
)

type Token struct {
	Type  TokenType
	value string
}

type Lex struct {
	input string
	pos   int  //current pos
	char  byte //current character
	state State
}

func NewLexer(input string) *Lex {
	l := &Lex{input: input, pos: -1, state: StateText}
	l.ReadChar()
	return l
}

func (l *Lex) ReadChar() {
	l.pos++
	if l.pos >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.pos]
	}
}

func (l *Lex) PeekAhead() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}

func (l *Lex) ReadNextToken() Token {

	switch l.state {
	case StateText:
		if l.char == '>' {
			l.state = StateQuote
		} else if l.char == '`' {
			l.state = StateInLineCode
		} else if l.char == '*' && l.PeekAhead() == '*' {
			l.state = StateBold
		} else if l.char == '*' && l.PeekAhead() != '*' {
			l.state = SateItalic
		} else {
			value := l.ReadText()
			return Token{Type: TEXT, value: value}
		}
	default:
		return Token{Type: EOF, value: ""}
	}

	return Token{Type: EOF, value: ""}
}

func (l *Lex) ReadText() string {
	start := l.pos
	if l.char == 0 {
		return ""
	}
	for l.char != 0 {
		l.ReadChar()
	}
	return l.input[start:l.pos]
}

func (l *Lex) NextToken() Token {
	start := l.pos
	if l.char == 0 {
		return Token{Type: EOF}
	}

	if l.char == '`' {
		l.ReadChar()
		val := l.ReadInlineCode()
		return Token{
			Type:  INLINE_CODE,
			value: val,
		}
	}

	if l.char == '*' && l.PeekAhead() != '*' {
		l.ReadChar()
		val := l.ReadItalic()
		return Token{
			Type:  ITALIC,
			value: val,
		}
	}

	if l.char == '*' && l.PeekAhead() == '*' {
		l.state = StateBold
		l.ReadChar()
		l.ReadChar()
		val := l.ReadItalic()
		l.ReadChar()
		return Token{
			Type:  BOLD,
			value: val,
		}
	}

	for l.char != '*' && l.char != '`' && l.char != 0 {
		l.ReadChar()
	}
	return Token{
		Type:  TEXT,
		value: l.input[start:l.pos],
	}
}

func (l *Lex) ReadInlineCode() string {
	start := l.pos
	for l.char != 0 && l.char != '`' {
		l.ReadChar()
	}
	value := l.input[start:l.pos]
	if l.char == '`' {
		l.ReadChar()
	}
	// if l.char=='`'{
	// 	l.ReadChar()
	// }
	return value
}

func (l *Lex) ReadItalic() string {
	start := l.pos
	for l.char != 0 && l.char != '*' {
		l.ReadChar()
	}
	value := l.input[start:l.pos]
	if l.char == '*' {
		l.ReadChar()
	}
	// if l.char=='`'{
	// 	l.ReadChar()
	// }
	return value
}
