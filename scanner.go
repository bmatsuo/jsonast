package jsonast

import (
	"github.com/bmatsuo/go-lexer"

	"unicode"
)

const (
	lLeftCurly lexer.ItemType = iota
	lRightCurly
	lLeftSquare
	lRightSquare
	lColon
	lComma
	lString
	lNumber
	lBoolean
	lNull
)

var spaceRunes = []rune(" \t\n")

func lexSlurpSpace(lex *lexer.Lexer) int {
	return lex.AcceptRunRunes(spaceRunes)
}

// NOTE assumes utf-8 input
func lexStart(lex *lexer.Lexer) lexer.StateFn {
	if lexSlurpSpace(lex) > 0 {
		lex.Ignore()
	}
	c, _ := lex.Peek()
	if c == lexer.EOF {
		lex.Emit(lexer.ItemEOF)
		return nil
	}
	typ, found := lexNom[c]
	if found {
		lex.Advance()
		lex.Emit(typ)
		return lexStart
	}
	switch c {
	case '"':
		return lexString
	case 't':
		return lexTrue
	case 'f':
		return lexFalse
	case 'n':
		return lexNull
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		return lexNumber
	default:
		return lex.Errorf("unexpected rune '%c'", c)
	}
}

var lexNom = map[rune]lexer.ItemType{
	'{': lLeftCurly,
	'}': lRightCurly,
	'[': lLeftSquare,
	']': lRightSquare,
	',': lComma,
	':': lColon,
}

func lexNull(lex *lexer.Lexer) lexer.StateFn {
	if !lex.AcceptRune('n') {
		lex.Errorf("expected 'n'")
		return nil
	}
	if !lex.AcceptRune('u') {
		lex.Errorf("expected 'u'")
		return nil
	}
	if !lex.AcceptRune('l') {
		lex.Errorf("expected 'l'")
		return nil
	}
	if !lex.AcceptRune('l') {
		lex.Errorf("expected 'l'")
		return nil
	}
	lex.Emit(lNull)
	return lexStart
}
func lexTrue(lex *lexer.Lexer) lexer.StateFn {
	if !lex.AcceptRune('t') {
		lex.Errorf("expected 't'")
		return nil
	}
	if !lex.AcceptRune('r') {
		lex.Errorf("expected 'r'")
		return nil
	}
	if !lex.AcceptRune('u') {
		lex.Errorf("expected 'u'")
		return nil
	}
	if !lex.AcceptRune('e') {
		lex.Errorf("expected 'e'")
		return nil
	}
	lex.Emit(lBoolean)
	return lexStart
}
func lexFalse(lex *lexer.Lexer) lexer.StateFn {
	if !lex.AcceptRune('f') {
		lex.Errorf("expected 'f'")
		return nil
	}
	if !lex.AcceptRune('a') {
		lex.Errorf("expected 'a'")
		return nil
	}
	if !lex.AcceptRune('l') {
		lex.Errorf("expected 'l'")
		return nil
	}
	if !lex.AcceptRune('s') {
		lex.Errorf("expected 's'")
		return nil
	}
	if !lex.AcceptRune('e') {
		lex.Errorf("expected 'e'")
		return nil
	}
	lex.Emit(lBoolean)
	return lexStart
}

func lexNumber(lex *lexer.Lexer) lexer.StateFn {
	lex.AcceptRune('-')
	if lex.AcceptRune('0') {
		return lexNumberFraction
	}
	n := lex.AcceptRunMinMax('1', '9')
	if n == 0 {
		return lex.Errorf("expected non-zero digit")
	}
	return lexNumberFraction
}
func lexNumberFraction(lex *lexer.Lexer) lexer.StateFn {
	if !lex.AcceptRune('.') {
		return lexNumberExponent
	}
	n := lex.AcceptRunMinMax('0', '9')
	if n == 0 {
		return lex.Errorf("expected digit")
	}
	return lexNumberExponent
}
func lexNumberExponent(lex *lexer.Lexer) lexer.StateFn {
	if !lex.Accept("eE") {
		lex.Emit(lNumber)
		return lexStart
	}
	lex.Accept("+-")
	n := lex.AcceptRunMinMax('0', '9')
	if n == 0 {
		return lex.Errorf("expected digit")
	}
	lex.Emit(lNumber)
	return lexStart
}

var escapeRunes = []rune(`"\/bfnrt`)

func lexString(lex *lexer.Lexer) lexer.StateFn {
	if !lex.AcceptRune('"') {
		return lex.Errorf("expected quote")
	}
	for {
		c, _ := lex.Advance()
		if c == '"' {
			lex.Emit(lString)
			return lexStart
		}
		if c == '\\' {
			if lex.AcceptRunes(escapeRunes) {
				continue
			}
			if lex.AcceptRune('u') {
				for i := 0; i < 4; i++ {
					accept := lex.AcceptMinMax('0', '9') ||
						lex.AcceptMinMax('a', 'f') ||
						lex.AcceptMinMax('A', 'F')
					if !accept {
						return lex.Errorf("expected hex digit")
					}
				}
			}
			continue
		}
		if unicode.IsControl(c) {
			return lex.Errorf("unexpected control character")
		}
	}
	panic("unreachable")
}
