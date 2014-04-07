package jsonast

type Type uint

const (
	TInvalid Type = iota
	TNull
	TString
	TNumber
	TBoolean
	TArray
	TObject
	tInvalid
)

func (t Type) IsValid() bool {
	return t > TInvalid && t < tInvalid
}

func IsNull(node ASTNode) bool {
	return node.Type() == TNull
}

type ASTNode interface {
	// the node's Type
	Type() Type
	// p with the node's json appended to it
	JSON(p []byte) []byte
	// the node's children. panics when the node's type is not TArray or
	// TObject. on a parsed node, if the node's type is TObject the returned
	// slice has an even number of children, with even-indexed elements of type
	// TString.
	Children() []ASTNode
	// add a child to the node. panics if the node's type is not TArray or
	// TObject.
	PushChild(ASTNode)
	// panics if Type() != TString
	String() string
	// panics if Type() != TNumber
	Float64() float64
	// panics if Type() != TNumber
	Int64() int64
	// panics if Type() != TBoolean
	Bool() bool
	// only types in this package can implement ASTNode
	sealASTNode()
}

type node struct {
	typ      Type
	raw      []byte    // only present for primatives
	children []ASTNode // only present for arrays and objects. key-value pairs are consecutive nodes.
}

func (nod *node) sealASTNode() {}

func (nod *node) Type() Type { return nod.typ }

// panics if nod.typ is note in unless
func (nod *node) panicType(unless ...Type) {
	for i := range unless {
		if nod.typ == unless[i] {
			return
		}
	}
	panic("invalid type")
}

func (nod *node) String() string {
	nod.panicType(TString)
	return "FIXME"
}

func (nod *node) Bool() bool {
	nod.panicType(TBoolean)
	return nod.raw[0] == 't'
}

func (nod *node) Int64() int64 {
	nod.panicType(TNumber)
	return -1
}

func (nod *node) Float64() float64 {
	nod.panicType(TNumber)
	return -1
}

func (nod *node) JSON(js []byte) []byte {
	if !nod.typ.IsValid() {
		panic("invalid node")
	}
	switch nod.typ {
	case TArray:
		return nod.arrayJSON(js)
	case TObject:
		return nod.objectJSON(js)
	default:
		return append(js, nod.raw...)
	}
}

func (nod *node) PushChild(c ASTNode) {
	if nod.typ == TObject || nod.typ == TArray {
		nod.children = append(nod.children, c)
		return
	}
	panic("invalid node type")
}
func (nod *node) Children() []ASTNode {
	if nod.typ == TObject || nod.typ == TArray {
		return nod.children
	}
	panic("invalid node type")
}

func (nod *node) arrayJSON(js []byte) []byte {
	n := len(nod.children)
	if n == 0 {
		return []byte{'[', ']'}
	}
	js = append(js, '[')
	for i := 0; i < n; i++ {
		if i > 0 {
			js = append(js, ',')
		}
		js = nod.children[i].JSON(js)
	}
	js = append(js, ']')
	return js
}

func (nod *node) objectJSON(js []byte) []byte {
	n := len(nod.children)
	if n == 0 {
		return []byte{'{', '}'}
	}
	if n%2 == 0 {
		panic("odd number of children")
	}
	js = append(js, '{')
	for i := 0; i < n; i += 2 {
		k, v := nod.children[i], nod.children[i+1]
		if k.Type() != TString {
			panic("non-string key")
		}
		if i > 0 {
			js = append(js, ',')
		}
		js = k.JSON(js)
		js = append(js, ':')
		js = v.JSON(js)
	}
	js = append(js, '}')
	return js
}
