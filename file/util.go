package main

import (
	"errors"

	"github.com/goccy/go-yaml/ast"
)

// getStringNode gets the string value of a string node, or an error if not a string node.
func getStringNode(node ast.Node) (string, error) {
	switch n := node.(type) {
	case *ast.StringNode:
		return n.Value, nil
	default:
		return "", errors.New("node is not string")
	}
}
