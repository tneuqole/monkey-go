package evaluator

import (
	"ast"
	"object"
)

func quote(node ast.Node) object.Object {
	return &object.Quote{Node: node}
}
