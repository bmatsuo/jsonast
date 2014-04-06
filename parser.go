package jsonast

type Parser interface {
	Parse([]byte) (ASTNode, error)
}
