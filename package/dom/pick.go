package dom

import "github.com/faiface/pixel"

// TODO: really, Pick should return a tree, because
// you can be over multiple things at once.
func Pick(node Node, pt pixel.Vec) Node {
	switch node.(type) {
	case *RectNode, *CircleNode, *TextNode:
		if node.Contains(pt) {
			return node
		}
		return nil
	case *GroupNode:
		// TODO: support transforms on groups
		for _, child := range node.Children() {
			res := Pick(child, pt)
			if res != nil {
				return res
			}
		}
		return nil
	}
	return nil
}
