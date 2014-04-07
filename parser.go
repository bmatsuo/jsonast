package jsonast

import (
	"github.com/bmatsuo/go-lexer"

	"fmt"
)

func Parse(p []byte) (roots []ASTNode, err error) {
	state := new(parseState)
	state.in = p
	state.lex = lexer.New(lexStart, p)
	err = state.loop(false)
	if err != nil {
		return nil, err
	}
	return state.roots, nil
}

func parse(p []byte, debug bool) (roots []ASTNode, err error) {
	state := new(parseState)
	state.in = p
	state.lex = lexer.New(lexStart, p)
	err = state.loop(debug)
	if err != nil {
		return nil, err
	}
	return state.roots, nil
}

type parseState struct {
	in    []byte
	lex   *lexer.Lexer
	stack []ASTNode
	top   ASTNode
	root  ASTNode
	roots []ASTNode
}

func (state *parseState) push(nod ASTNode) (isroot bool) {
	if nod == nil {
		panic("nil node")
	}
	n := len(state.stack)
	if n == 0 {
		state.root = nod
		state.roots = append(state.roots, nod)
		isroot = true
	} else {
		state.stack[n-1].PushChild(nod)
	}
	state.stack = append(state.stack, nod)
	state.top = nod
	return isroot
}

func (state *parseState) pop() (nod ASTNode, isempty bool) {
	n := len(state.stack)
	if n == 0 {
		return nil, true
	}
	nod = state.top
	state.stack[n-1] = nil
	state.stack = state.stack[:n-1 : cap(state.stack)]
	n-- // new len
	if n > 0 {
		state.top = state.stack[n-1]
		return nod, false
	}
	state.top = nil
	return nod, true
}

type ParseError struct {
	Reason string
	*lexer.Item
	in []byte
}

func (err *ParseError) Error() string {
	// reported 'offset' is dump. should be (line, col)
	return fmt.Sprintf("offset %d: %s", err.Item.Pos, err.Reason)
}

func (state *parseState) unexpected(item *lexer.Item) error {
	return &ParseError{fmt.Sprintf("unexpected '%s'", item.String()), item, state.in}
}

func (state *parseState) empty(item *lexer.Item) error {
	return &ParseError{"empty", item, state.in}
}

func (state *parseState) eof(item *lexer.Item) error {
	if len(state.stack) == 0 {
		return nil
	}
	return state.unexpected(item)
}

// loop sets the parse root to the first value seen and iterates until
// there are no values on the stack.
func (state *parseState) loop(debug bool) error {
	for {
		item := state.lex.Next()
		if item == nil {
			panic("missing EOF")
		}
		err := item.Error()
		if err != nil {
			return err
		}
		if debug {
			if state.top != nil {
				fmt.Printf("top: %d children: %d\n", state.top.Type(), len(state.top.Children()))
			} else {
				fmt.Println("empty stack")
			}
		}
		if item.Type == lexer.ItemEOF {
			// do checks
			return state.eof(item)
		}
		if item.Type == lString {
			// optimization because of high string frequency
			err = parseJumpString(state, item)
		} else {
			jmp, ok := parseJump[item.Type]
			if ok {
				err = jmp(state, item)
				if err != nil {
					return err
				}
			} else {
				return state.unexpected(item)
			}
		}

		state.lex.Free(item)
	}
}

func parseJumpString(state *parseState, item *lexer.Item) error {
	state.push(&node{typ: TString, raw: item.Value})
	state.pop()
	return nil
}

var parseJump = map[lexer.ItemType]func(state *parseState, item *lexer.Item) error{
	lLeftCurly: func(state *parseState, item *lexer.Item) error {
		state.push(&node{typ: TObject})
		return nil
	},
	lRightCurly: func(state *parseState, item *lexer.Item) error {
		nod, _ := state.pop()
		if nod == nil {
			return state.unexpected(item)
		}
		if nod.Type() != TObject {
			return state.unexpected(item)
		}
		return nil
	},
	lLeftSquare: func(state *parseState, item *lexer.Item) error {
		state.push(&node{typ: TArray})
		return nil
	},
	lRightSquare: func(state *parseState, item *lexer.Item) error {
		nod, _ := state.pop()
		if nod == nil {
			return state.unexpected(item)
		}
		if nod.Type() != TArray {
			return state.unexpected(item)
		}
		return nil
	},
	lColon: func(state *parseState, item *lexer.Item) error {
		if state.top == nil {
			return state.unexpected(item)
		}
		if state.top.Type() != TObject {
			return state.unexpected(item)
		}
		if len(state.top.Children())%2 == 0 {
			return state.unexpected(item)
		}
		return nil
	},
	lComma: func(state *parseState, item *lexer.Item) error {
		if state.top == nil {
			return state.unexpected(item)
		}
		if state.top.Type() == TObject {
			if len(state.top.Children())%2 == 1 {
				return state.unexpected(item)
			}
		} else if state.top.Type() != TArray {
			return state.unexpected(item)
		}
		return nil
	},
	lString: parseJumpString,
	lNumber: func(state *parseState, item *lexer.Item) error {
		state.push(&node{typ: TNumber, raw: item.Value})
		state.pop()
		return nil
	},
	lBoolean: func(state *parseState, item *lexer.Item) error {
		state.push(&node{typ: TBoolean, raw: item.Value})
		state.pop()
		return nil
	},
	lNull: func(state *parseState, item *lexer.Item) error {
		state.push(&node{typ: TNull, raw: item.Value})
		state.pop()
		return nil
	},
}
