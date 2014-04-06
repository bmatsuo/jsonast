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

// NOTE assumes utf-8 input
func lexStart(lex *lexer.Lexer) lexer.StateFn {
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
	if !lex.Accept("n") {
		lex.Errorf("expected 'n'")
		return nil
	}
	if !lex.Accept("u") {
		lex.Errorf("expected 'u'")
		return nil
	}
	if !lex.Accept("l") {
		lex.Errorf("expected 'l'")
		return nil
	}
	if !lex.Accept("l") {
		lex.Errorf("expected 'l'")
		return nil
	}
	lex.Emit(lNull)
	return lexStart
}
func lexTrue(lex *lexer.Lexer) lexer.StateFn {
	if !lex.Accept("t") {
		lex.Errorf("expected 't'")
		return nil
	}
	if !lex.Accept("r") {
		lex.Errorf("expected 'r'")
		return nil
	}
	if !lex.Accept("u") {
		lex.Errorf("expected 'u'")
		return nil
	}
	if !lex.Accept("e") {
		lex.Errorf("expected 'e'")
		return nil
	}
	lex.Emit(lBoolean)
	return lexStart
}
func lexFalse(lex *lexer.Lexer) lexer.StateFn {
	if !lex.Accept("f") {
		lex.Errorf("expected 'f'")
		return nil
	}
	if !lex.Accept("a") {
		lex.Errorf("expected 'a'")
		return nil
	}
	if !lex.Accept("l") {
		lex.Errorf("expected 'l'")
		return nil
	}
	if !lex.Accept("s") {
		lex.Errorf("expected 's'")
		return nil
	}
	if !lex.Accept("e") {
		lex.Errorf("expected 'e'")
		return nil
	}
	lex.Emit(lBoolean)
	return lexStart
}

// can produce bad numbers. but there will not be any ambiguity of syntax error when that happens.
func lexNumber(lex *lexer.Lexer) lexer.StateFn {
	lex.Accept("-")
	if lex.Accept("0") {
		return lexNumberFraction
	}
	lex.AcceptRun("123456789")
	/*
	n := lex.AcceptRun("123456789")
	if n == 0 {
		return lex.Errorf("expected non-zero digit")
	}
	*/
	return lexNumberFraction
}
func lexNumberFraction(lex *lexer.Lexer) lexer.StateFn {
	if !lex.Accept(".") {
		return lexNumberExponent
	}
	lex.AcceptRun("0123456789")
	/*
	n := lex.AcceptRun("0123456789")
	if n == 0 {
		return lex.Errorf("expected digit")
	}
	*/
	return lexNumberExponent
}
func lexNumberExponent(lex *lexer.Lexer) lexer.StateFn {
	if !lex.Accept("eE") {
		lex.Emit(lNumber)
		return lexStart
	}
	lex.Accept("+-")
	lex.AcceptRun("0123456789")
	/*
	n := lex.AcceptRun("0123456789")
	if n == 0 {
		return lex.Errorf("expected digit")
	}
	*/
	lex.Emit(lNumber)
	return lexStart
}

func lexString(lex *lexer.Lexer) lexer.StateFn {
	if !lex.Accept(`"`) {
		return lex.Errorf("expected quote")
	}
	for {
		if lex.Accept(`"`) {
			lex.Emit(lString)
			return lexStart
		}
		if lex.Accept(`\`) {
			if lex.Accept(`"\/bfnrt`) {
				continue
			}
			if lex.Accept("u") {
				for i := 0; i < 4; i++ {
					if !lex.Accept("0123456789abcdefABCDEF") {
						return lex.Errorf("return expected hex digit")
					}
				}
			}
		}
		c, _ := lex.Advance()
		if unicode.IsControl(c) {
			return lex.Errorf("unexpected control character")
		}
	}
	panic("unreachable")
}
