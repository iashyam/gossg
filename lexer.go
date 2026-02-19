package main

type TokenType int

const (
	EOF TokenType = iota
	TEXT
	ITALIC
	BOLD
	INLINE_CODE
	QUOTE
	NEWLINE
	BLANKLINE
	HEADING
	LIST
	MATH
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
	case BLANKLINE:
		return "BLANKLINE"
	case NEWLINE:
		return "NEWLINE"
	case HEADING:
		return "HEADING"
	case LIST:
		return "LIST"
	case MATH:
		return "MATH"
	default:
		return "NONE"
	}
}

type State int

const (
	StateText State = iota
	StateBold
	StateItalic
	StateInLineCode
	StateQuote
)

type Token struct {
	Type  TokenType
	value string
}

type Lex struct {
	input         string
	pos           int  //current pos
	char          byte //current character
	isAtLineStart bool
	state         State
}

func NewLexer(input string) *Lex {
	l := &Lex{input: input, pos: -1, state: StateText, isAtLineStart: true}
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
	if l.pos >= (len(l.input) - 1) {
		return 0
	}
	return l.input[l.pos+1]
}

func (l *Lex) QuoteHandler() Token {
	l.ReadChar() // '>'
	l.ReadChar() // ' '
	return Token{Type: QUOTE, value: "> "}
}

func (l *Lex) HeadingHandler() Token {
	tokenVal := ""
	for l.char == '#' {
		tokenVal += string(l.char)
		l.ReadChar()
	}
	tokenVal += " "
	l.ReadChar() // ' '
	l.isAtLineStart = false
	return Token{Type: HEADING, value: tokenVal}
}

func (l *Lex) BoldHandler() Token {
	l.ReadChar() // '*'
	l.ReadChar() // '*'
	return Token{Type: BOLD, value: "**"}
}

func (l *Lex) ItalicHandler() Token {
	l.ReadChar() // '*'
	return Token{Type: ITALIC, value: "*"}
}

func (l *Lex) InlineCodeHandler() Token {
	l.ReadChar() // '`'
	return Token{Type: INLINE_CODE, value: "`"}
}

func (l *Lex) ReadText() string {
	start := l.pos
	if l.char == 0 {
		return ""
	}
	for l.char != 0 {
		if (l.state != StateInLineCode && (l.char == '*' || l.char == '$')) || l.char == '\n' || l.char == '`' || l.char == '\\' {
			break
		}
		l.ReadChar()
	}
	return l.input[start:l.pos]
}

func (l *Lex) ListHandler() Token {
	l.ReadChar() // '-'
	l.ReadChar() // ' '
	return Token{Type: LIST, value: "- "}
}

func (l *Lex) MathHandler() Token {
	l.ReadChar() // '$'
	start := l.pos
	for l.char != 0 && l.char != '$' {
		l.ReadChar()
	}
	if l.char == '$' {
		l.ReadChar() // consume closing '$'
		return Token{Type: MATH, value: l.input[start : l.pos-1]}
	}
	return Token{Type: MATH, value: l.input[start:l.pos]}
}

func (l *Lex) ReadNextToken() Token {

	switch l.state {
	case StateText, StateItalic, StateBold, StateQuote:
		if l.char == 0 {
			return Token{Type: EOF, value: ""}
		}

		if l.char == '\\' {
			l.ReadChar()
			char := l.char
			l.ReadChar()
			return Token{Type: TEXT, value: string(char)}
		}

		if l.char == '\n' {
			l.isAtLineStart = true
			if l.PeekAhead() == '\n' {
				l.ReadChar()
				l.ReadChar()
				return Token{Type: BLANKLINE, value: "BLANKLINE"}
			}
			l.ReadChar()

			return Token{Type: NEWLINE, value: "NEWLINE"}
		}

		if l.isAtLineStart && l.char == '#' {
			count := 0
			for i := l.pos; i < len(l.input); i++ {
				if l.input[i] == '#' {
					count++
				} else {
					break
				}
			}

			if count > 0 && count <= 6 && l.pos+count < len(l.input) && l.input[l.pos+count] == ' ' {
				return l.HeadingHandler()
			}
		}

		if l.isAtLineStart && (l.char == '-' || l.char == '*' || l.char == '+') && l.PeekAhead() == ' ' {
			return l.ListHandler()
		}
		if l.char == '>' && l.PeekAhead() == ' ' {
			l.state = StateQuote
			return l.QuoteHandler()
		} else if l.char == '`' {
			l.state = StateInLineCode
			return l.InlineCodeHandler()
		} else if l.char == '*' && l.PeekAhead() == '*' {
			l.state = StateBold
			return l.BoldHandler()
		} else if l.char == '*' && l.PeekAhead() != '*' {
			l.state = StateItalic
			return l.ItalicHandler()
		} else if l.char == '$' {
			return l.MathHandler()
		} else {
			l.isAtLineStart = false
			value := l.ReadText()
			return Token{Type: TEXT, value: value}
		}
	case StateInLineCode:
		if l.char == 0 {
			return Token{Type: EOF, value: ""}
		} else if l.char == '`' {
			l.state = StateText
			return l.InlineCodeHandler()
		} else {
			value := l.ReadText()
			return Token{Type: TEXT, value: value}
		}

	default:
		return Token{Type: EOF, value: ""}
	}
}
