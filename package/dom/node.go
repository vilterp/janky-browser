package dom

import (
	"github.com/faiface/pixel"
)

type Node interface {
	Name() string
	Attrs() map[string]string
	Children() []Node

	Init()
	Draw(target pixel.Target)
	Contains(pixel.Vec) bool
	GetBounds() pixel.Rect
}

func GetAllNodes(tree Node) []Node {
	var output []Node
	Visit(tree, func(n Node) {
		output = append(output, n)
	})
	return output
}

// Visit does a pre-order traversal of the DOM tree.
func Visit(tree Node, visitor func(n Node)) {
	visitor(tree)
	for _, child := range tree.Children() {
		Visit(child, visitor)
	}
}
