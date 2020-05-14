package dom

import (
	"image"

	"github.com/llgcode/draw2d"
)

type Node interface {
	Name() string
	Attrs() map[string]string
	Children() []Node

	Init()
	Draw(ctx draw2d.GraphicContext)
	Contains(image.Point) bool
	GetBounds() image.Rectangle

	Events() *EventHandlers
}

type baseNode struct {
	events EventHandlers
}

func (bn *baseNode) Events() *EventHandlers {
	return &bn.events
}

func GetAllNodes(tree Node) []Node {
	var output []Node
	SimpleVisit(tree, func(n Node, _ int) {
		output = append(output, n)
	})
	return output
}

// SimpleVisit does a pre-order traversal of the DOM tree.
func SimpleVisit(tree Node, visitor func(n Node, depth int)) {
	Visit(tree, visitor, nil)
}

func Visit(
	tree Node, beforeChildren func(n Node, depth int), afterChildren func(n Node, depth int),
) {
	doVisit(tree, beforeChildren, afterChildren, 0)
}

func doVisit(
	tree Node,
	beforeChildren func(n Node, depth int),
	afterChildren func(n Node, depth int),
	depth int,
) {
	if beforeChildren != nil {
		beforeChildren(tree, depth)
	}
	for _, child := range tree.Children() {
		doVisit(child, beforeChildren, afterChildren, depth+1)
	}
	if afterChildren != nil {
		afterChildren(tree, depth)
	}
}
