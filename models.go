package main


type TokenType int 

const (
	EOF TokenType = iota
	TEXT
	ITALIC
	INLINE_CODE
	QUOTE
)

type Token struct{
	Type TokenType
	value string
}


type Lex struct{
	input string 
	pos   int //current pos
	char  byte //current character
}

func NewLexer(input string) *Lex{
	l := &Lex{input: input, pos: -1}
	l.ReadChar()
	return  l
}

func (l *Lex) ReadChar(){
	l.pos++
	if l.pos>=len(l.input){
		l.char = 0
	}else{
		l.char = l.input[l.pos]
	}
}

func (l *Lex) NextToken() Token{
	start := l.pos
	if l.char==0{
		return Token{Type: EOF}
	}

	if l.char=='`'{
		l.ReadChar()
		val := l.ReadInlineCode()
		return Token{
			Type: INLINE_CODE,
			value: val,
		}
	}

	if l.char=='*'{
		l.ReadChar()
		val := l.ReadItalic()
		return Token{
			Type: ITALIC,
			value: val,
		}
	}
	for l.char!='*' && l.char!='`' && l.char != 0{
		l.ReadChar()
	}
	return Token{
		Type: TEXT,
		value: l.input[start:l.pos],
	}
}

func (l *Lex) ReadInlineCode() string{
	start := l.pos
	for l.char!=0 && l.char!='`'{
		l.ReadChar()
	}
	value := l.input[start:l.pos]
	if l.char=='`'{
	l.ReadChar()
	}
	// if l.char=='`'{
	// 	l.ReadChar()
	// }
	return value
}


func (l *Lex) ReadItalic() string{
	start := l.pos
	for l.char!=0 && l.char!='*'{
		l.ReadChar()
	}
	value := l.input[start:l.pos]
	if l.char=='*'{
		l.ReadChar()
	}
	// if l.char=='`'{
	// 	l.ReadChar()
	// }
	return value
}
