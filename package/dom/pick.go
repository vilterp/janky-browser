package dom

import (
	"image"
)

// TODO: really, Pick should return a tree, because
// you can be over multiple things at once.
func Pick(node Node, pt image.Point) []Node {
	switch node.(type) {
	case *RectNode, *CircleNode, *TextNode:
		if node.Contains(pt) {
			return []Node{node}
		}
		return []Node{}
	case *GroupNode:
		// TODO: support transforms on groups
		var res []Node
		for _, child := range node.Children() {
			childRes := Pick(child, pt)
			res = append(res, childRes...)
		}
		if len(res) > 0 {
			res = append(res, node)
		}
		return res
	}
	return []Node{}
}
