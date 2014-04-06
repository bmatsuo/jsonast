package jsonast

import (
	"github.com/bmatsuo/go-lexer"

	"fmt"
)

var ErrUnexpectedEOF = fmt.Errorf("unexpected end of input")

type parseState struct {
	in    []byte
	lex   *lexer.Lexer
	stack []ASTNode
	root  ASTNode
}

func (state *parseState) pop() {
	n := len(state.stack)
	if n == 0 {
		panic("pop with an empty stack")
	}
	state.stack[n] = nil
	state.stack = state.stack[:n-1]
}

/*
XXX
XXX this is not a good start
XXX
*/

type parseStateFn func(*parseState) (parseStateFn, error)

func parseStart(state *parseState) (parseStateFn, error) {
	for {
		item := state.lex.Next()
		if item == nil {
			return nil, ErrUnexpectedEOF
		}
		if item.Type == lLeftCurly {
			nod := new(node)
			nod.typ = TObject
			if state.root == nil {
				state.root = nod
				state.stack = append(state.stack, state.root)
			} else {
				panic("non-empty stack")
			}
			return parseObject, nil
		}
		if item.Type == lLeftSquare {
			if state.root == nil {
				nod := new(node)
				nod.typ = TArray
				state.root = nod
				state.stack = append(state.stack, state.root)
			} else {
				panic("non-empty stack")
			}
			return parseArray, nil
		}
	}
	panic("unreachable")
}

func parseObject(state *parseState) (parseStateFn, error) {
	for {
		item := state.lex.Next()
		if item.Type != lString {
			return nil, fmt.Errorf("expected string")
		}
		nod := new(node)
		nod.typ = TString
		nod.raw = []byte(item.Value)
		item = state.lex.Next()
		if item.Type != lColon {
			return nil, fmt.Errorf("expected ':'")
		}
		item = state.lex.Next()
		if item == nil {
			return nil, ErrUnexpectedEOF
		}
		switch item.Type {
		}
		state.stack[len(state.stack)-1].PushChild(nod)
		if item.Type == lRightCurly {
			state.pop()
			return parseStart, nil
		}
	}
	panic("unreachable")
}

func parseArray(state *parseState) (parseStateFn, error) {
	for {
		item := state.lex.Next()
		if item == nil {
			return nil, ErrUnexpectedEOF
		}
		if item.Type == lRightSquare {
			state.pop()
			return parseStart, nil
		}
	}
	panic("unreachable")
}

func (state *parseState) parseValue() error {
	for {
		item := state.lex.Next()
		if item == nil {

		}
	}
	panic("unreachable")
}

func Parse(p []byte) (ASTNode, error) {
	state := parseState{
		in:  p,
		lex: lexer.New(lexStart, string(p)),
	}
	_ = state
	return nil, fmt.Errorf("not yet")
}
