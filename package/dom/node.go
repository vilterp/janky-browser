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
	output = append(output, tree)
	for _, child := range tree.Children() {
		output = append(output, GetAllNodes(child)...)
	}
	return output
}
