package nlexing

import (
	"strings"
	"unicode"

	"github.com/influx6/npkg/nerror"
)

const (
	VariantToken TokenType = iota + 10
	PrefixToken
	GroupStartToken
	GroupEndToken
	TargetFinished
	TargetToken
	NegationToken
)

const (
	dash              = '-'
	prefixer          = '~'
	underscore        = '_'
	colon             = ':'
	comma             = ','
	apostrophe        = '!'
	leftBracketDelim  = '('
	rightBracketDelim = ')'
)

type intStack []int

func (th *intStack) Clear() {
	*th = (*th)[:0]
}

func (th *intStack) ClearCount(n int) {
	var count = len(*th)
	if count == 0 {
		return
	}
	*th = (*th)[0 : count-n]
}

func (th *intStack) Len() int {
	return len(*th)
}

func (th *intStack) Pop() int {
	var count = len(*th)
	if count == 0 {
		return -1
	}
	var last = (*th)[count-1]
	*th = (*th)[0 : count-1]
	return last
}

func (th *intStack) Push(t int) {
	*th = append(*th, t)
}

type stack []string

func (th *stack) Join(with string) string {
	return strings.Join(*th, with)
}

func (th *stack) Len() int {
	return len(*th)
}

func (th *stack) Clear() {
	*th = (*th)[:0]
}

func (th *stack) ClearCount(n int) {
	var count = len(*th)
	if count == 0 {
		return
	}
	*th = (*th)[0 : count-n]
}

func (th *stack) Pop() string {
	var count = len(*th)
	if count == 0 {
		return ""
	}
	var last = (*th)[count-1]
	*th = (*th)[0 : count-1]
	return last
}

func (th *stack) Push(t string) {
	*th = append(*th, t)
}

// ParseVariantDirectives implements a parser which handles parsing
// tokens in the following formats:
// 1. Variant:Prefix-text-text*
// 2. Variant:(Prefix-text, Prefix2-text2) which expanded is Variant:Prefix-text and Variant:Prefix2-text2
// 3. Variant:(Variant2:Prefix-text, Prefix-text) which expanded is Variant:Variant2:Prefix-text and Variant:Prefix-text
// 4. Variant:(Variant2:Prefix-text, Variant3:Prefix-text) which expanded is Variant:Variant2:Prefix-text and Variant:Variant3:Prefix-text
// 5. Variant:(Variant2:(Prefix-text, Prefix2-text2), Variant3:Prefix-text3) which expanded is Variant:Variant2:Prefix-text, Variant:Variant2:Prefix2-text2 and Variant:Variant3:Prefix-text
// 6. Pf~(Prefix-text, Prefix2-text2) which expanded is Variant:Pf-Prefix-text and Variant:Pf-Prefix2-text2
//
func ParseVariantDirectives(v string) ([]string, error) {
	var parsed []string

	var cls stack
	var gps intStack

	var pls stack
	var pgps intStack

	var variantCount int
	var prefixCount int
	var groups int

	var lexer = NewLexer(v)
	var tokenizer = NewTokenizer(lexer, LexCompactDirective, func(b string, t TokenType) error {
		switch t {
		case TargetToken:
			if len(b) == 0 {
				break
			}

			gps.Push(variantCount)
			variantCount = 0

			pgps.Push(prefixCount)
			prefixCount = 0

			var prefixed = b
			if pls.Len() > 0 {
				prefixed = pls.Join("-") + "-" + prefixed
			}

			if cls.Len() > 0 {
				prefixed = cls.Join(":") + ":" + prefixed
			}

			parsed = append(parsed, prefixed)
		case TargetFinished:
			if gps.Len() > 0 {
				var grpVariantCount = gps.Pop()
				cls.ClearCount(grpVariantCount)
			}
			if pgps.Len() > 0 {
				var pgrpVariantCount = pgps.Pop()
				pls.ClearCount(pgrpVariantCount)
			}
		case PrefixToken:
			prefixCount++
			pls.Push(b)
		case VariantToken:
			variantCount++
			cls.Push(b)
		case GroupStartToken:
			groups++
			gps.Push(variantCount)
			variantCount = 0

			pgps.Push(prefixCount)
			prefixCount = 0
		case GroupEndToken:
			groups--
			var grpVariantCount = gps.Pop()
			cls.ClearCount(grpVariantCount)

			var pgrpVariantCount = pgps.Pop()
			pls.ClearCount(pgrpVariantCount)
		}
		return nil
	})

	if err := tokenizer.Run(); err != nil {
		return nil, nerror.WrapOnly(err)
	}

	if groups < 0 {
		return parsed, nerror.New("seems grouping closer ')' is more than opener '('")
	}

	if groups > 0 {
		return parsed, nerror.New("seems grouping opener '(' is more than closer ')'")
	}

	return parsed, nil
}

// LexCompactDirective implements a parser which handles parsing
// tokens in the following formats:
// 1. Variant:Prefix-text-text*
// 2. Variant:(Prefix-text, Prefix2-text2) which expanded is Variant:Prefix-text and Variant:Prefix2-text2
// 3. Variant:(Variant2:Prefix-text, Prefix-text) which expanded is Variant:Variant2:Prefix-text and Variant:Prefix-text
// 4. Variant:(Variant2:Prefix-text, Variant3:Prefix-text) which expanded is Variant:Variant2:Prefix-text and Variant:Variant3:Prefix-text
// 5. Variant:(Variant2:(Prefix-text, Prefix2-text2), Variant3:Prefix-text3) which expanded is Variant:Variant2:Prefix-text, Variant:Variant2:Prefix2-text2 and Variant:Variant3:Prefix-text
// 6. Pf~(Prefix-text, Prefix2-text2) which expanded is Variant:Pf-Prefix-text and Variant:Pf-Prefix2-text2
//
func LexCompactDirective(l *Lexer, result ResultFunc) (TokenFunc, error) {
	if l.isAtEnd() {
		return nil, nil
	}

	return lexVariant, nil
}

// This specifically searches for ([\w\d]+): matching tokens
// which then are sent as the variant token type
func lexVariant(l *Lexer, result ResultFunc) (TokenFunc, error) {
	var pr rune
	for {
		if l.isAtEnd() {
			return LexCompactDirective, nil
		}

		if lexSpaceUntil(l) {
			l.ignore()
			continue
		}

		var lexedText = lexTextUntil(l)

		// if we lex text values, then check what is the next
		// non text token
		pr = l.peek()
		switch pr {
		case eof:
			l.next()
			if lexedText {
				if err := result(l.slice(), TargetToken); err != nil {
					return nil, nerror.WrapOnly(err)
				}
			}
			return nil, nil
		case comma:
			if lexedText {
				if err := result(l.slice(), TargetToken); err != nil {
					return nil, nerror.WrapOnly(err)
				}
			}
			if err := result("", TargetFinished); err != nil {
				return nil, nerror.WrapOnly(err)
			}
			l.skipNext()
			return lexVariant, nil
		case prefixer:
			if err := result(l.slice(), PrefixToken); err != nil {
				return nil, nerror.WrapOnly(err)
			}
			l.skipNext()
			return lexVariant, nil
		case colon:
			if err := result(l.slice(), VariantToken); err != nil {
				return nil, nerror.WrapOnly(err)
			}
			l.skipNext()
			return lexVariant, nil
		case apostrophe:
			if lexedText {
				if err := result(l.slice(), TargetToken); err != nil {
					return nil, nerror.WrapOnly(err)
				}
			}

			if err := result("", NegationToken); err != nil {
				return nil, nerror.WrapOnly(err)
			}
			l.skipNext()
			return lexVariant, nil
		case rightBracketDelim:
			if lexedText {
				if err := result(l.slice(), TargetToken); err != nil {
					return nil, nerror.WrapOnly(err)
				}
			}

			if err := result("", GroupEndToken); err != nil {
				return nil, nerror.WrapOnly(err)
			}
			l.skipNext()
			return lexVariant, nil
		case leftBracketDelim:
			if lexedText {
				if err := result(l.slice(), TargetToken); err != nil {
					return nil, nerror.WrapOnly(err)
				}
			}

			if err := result("", GroupStartToken); err != nil {
				return nil, nerror.WrapOnly(err)
			}
			l.skipNext()
			return lexVariant, nil
		default:
			if !isAlphaNumeric(pr) {
				return nil, nerror.New("undefined token %q found while lexing %q in %q", pr, l.slice(), l.input)
			}
		}

		if err := result(l.slice(), TargetToken); err != nil {
			return nil, nerror.WrapOnly(err)
		}
	}
}

// func lexGroupStart(l *Lexer, resultFunc ResultFunc) TokenFunc {
//
// 	return nil
// }

func isDash(r rune) bool {
	return r == dash
}

func isUnderscore(r rune) bool {
	return r == underscore
}

func isColon(r rune) bool {
	return r == colon
}

func isLeftBracket(r rune) bool {
	return r == leftBracketDelim
}

func isRightBracket(r rune) bool {
	return r == rightBracketDelim
}

// lexTextUntil scans a run of alphaneumeric characters.
func lexTextUntil(l *Lexer) bool {
	var found = false
	var numSpaces int
	for {
		if l.isAtEnd() {
			break
		}

		var r = l.peek()
		if !isAlphaNumeric(r) {
			break
		}
		l.next()
		found = true
		numSpaces++
	}
	return found
}

// lexSpaceUntil scans a run of alphaneumeric characters.
func lexSpaceUntil(l *Lexer) bool {
	var foundSpace = false
	var r rune
	for {
		r = l.peek()
		if !isSpace(r) {
			break
		}
		foundSpace = true
		l.next()
	}
	return foundSpace
}

// lexTextWith scans a run of alphanumeric characters.
func lexTextWith(fn TokenFunc) TokenFunc {
	return func(l *Lexer, rs ResultFunc) (TokenFunc, error) {
		var r rune
		for {
			r = l.peek()
			if !isAlphaNumeric(r) {
				break
			}
			l.next()
		}
		return fn, nil
	}
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == underscore || r == apostrophe || r == dash || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumericAndDot(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.'
}
