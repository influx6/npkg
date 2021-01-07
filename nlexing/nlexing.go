package nlexing

import (
	"unicode/utf8"
)

const eof = rune(0)

type TokenType int

type ResultFunc func(b string, t TokenType) error

type TokenFunc func(l *Lexer, res ResultFunc) (TokenFunc, error)

type Tokenizer struct {
	l        *Lexer
	start    TokenFunc
	resultFn ResultFunc
}

func NewTokenizer(l *Lexer, starter TokenFunc, results ResultFunc) *Tokenizer {
	return &Tokenizer{
		l:        l,
		start:    starter,
		resultFn: results,
	}
}

func (t *Tokenizer) Run() error {
	var stateFn = t.start
	var err error
	for {
		if stateFn == nil {
			break
		}
		stateFn, err = stateFn(t.l, t.resultFn)
	}
	return err
}

type Lexer struct {
	pos   int
	width int
	start int
	input string
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		pos:   0,
		width: 0,
		start: 0,
		input: input,
	}
}

func (b *Lexer) isAtEnd() bool {
	return b.pos >= len(b.input)
}

func (b *Lexer) backup() {
	b.pos -= b.width
}

func (b *Lexer) next() rune {
	if b.pos >= len(b.input) {
		b.width = 0
		return eof
	}
	var rn, width = utf8.DecodeRuneInString(b.input[b.pos:])
	b.width = width
	b.pos += width
	return rn
}

func (b *Lexer) skipNext() rune {
	var nr = b.next()
	b.ignore()
	return nr
}

func (b *Lexer) peek() rune {
	var nr = b.next()
	b.backup()
	return nr
}

func (b *Lexer) skipRune() {
	b.start += b.width
}

func (b *Lexer) ignore() {
	b.start = b.pos
}

func (b *Lexer) peekSlice() string {
	// if b.pos >= len(b.input) {
	// 	return b.input
	// }
	var sl = b.input[b.start:b.pos]
	return sl
}

func (b *Lexer) slice() string {
	if b.pos >= len(b.input) {
		return b.input[b.start:]
	}
	var sl = b.input[b.start:b.pos]
	b.start = b.pos
	return sl
}
