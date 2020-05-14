package jankybrowser

import (
	"image"

	"github.com/llgcode/draw2d"
	"github.com/vilterp/janky-browser/package/dom"
)

type ContentRenderer struct {
	rootNode dom.Node // set when state = PageStateLoaded

	// the set of nodes the mouse was over when it was pressed.
	// empty if the mouse has not been pressed.
	mouseDownNodes map[dom.Node]bool

	mouseOverNodes map[dom.Node]bool

	highlightedNode dom.Node
}

func NewContentRenderer(rootNode dom.Node) *ContentRenderer {
	cr := &ContentRenderer{
		rootNode:       rootNode,
		mouseDownNodes: make(map[dom.Node]bool),
		mouseOverNodes: make(map[dom.Node]bool),
	}
	cr.rootNode.Init()
	return cr
}

// processClickState steps the click state machine, returning clicked nodes if there are any.
func (cr *ContentRenderer) processClickState(
	pt image.Point, mouseDown bool, mouseJustDown bool,
) []dom.Node {
	hoveredNodes := cr.GetHoveredNodes(pt)

	// Find nodes the mouse just went out of.
	// mouseOutNodes := cr.mouseOverNodes - hoveredNodes
	for wasOverNode := range cr.mouseOverNodes {
		if _, ok := hoveredNodes[wasOverNode]; !ok {
			if wasOverNode.Events().OnMouseOut != nil {
				wasOverNode.Events().OnMouseOut()
			}
			delete(cr.mouseOverNodes, wasOverNode)
		}
	}

	// Find nodes that we just went into.
	// mouseJustOverNodes = hoveredNodes - mouseOverNodes
	for hoveredNode := range hoveredNodes {
		if _, ok := cr.mouseOverNodes[hoveredNode]; !ok {
			cr.mouseOverNodes[hoveredNode] = true
			if hoveredNode.Events().OnMouseOver != nil {
				hoveredNode.Events().OnMouseOver()
			}
		}
	}

	var clickedNodes []dom.Node
	if mouseJustDown {
		// Record nodes the mouse was over when it was clicked.
		// copy(cr.mouseDownNodes, hoveredNodes)
		cr.mouseDownNodes = make(map[dom.Node]bool, len(hoveredNodes))
		for hoveredNode := range hoveredNodes {
			cr.mouseDownNodes[hoveredNode] = true
		}
	} else if !mouseDown && len(cr.mouseDownNodes) > 0 {
		// Mouse was just released. Find which nodes were clicked.
		// clickedNodes = intersect(cr.mouseDownNodes, mouseOverNodes)
		for hoveredNode := range hoveredNodes {
			if _, ok := cr.mouseDownNodes[hoveredNode]; ok {
				// This node was clicked.
				if hoveredNode.Events().OnClick != nil {
					hoveredNode.Events().OnClick()
				}
				clickedNodes = append(clickedNodes, hoveredNode)
			}
		}
		cr.mouseDownNodes = nil
	}
	return clickedNodes
}

func (cr *ContentRenderer) GetHoveredNodes(pt image.Point) map[dom.Node]bool {
	picked := dom.Pick(cr.rootNode, pt)
	// TODO: maybe have it returned as a map from pick
	asMap := make(map[dom.Node]bool, len(picked))
	for _, node := range picked {
		asMap[node] = true
	}
	return asMap
}

func (cr *ContentRenderer) Draw(gc draw2d.GraphicContext) {
	cr.rootNode.Draw(gc)

	// Draw highlight rect if we have a highlighted node.
	if cr.highlightedNode == nil {
		return
	}
	highlightRect := dom.RectFromBounds(cr.highlightedNode.GetBounds())
	highlightRect.Stroke = "red"
	highlightRect.Draw(gc)
}

func (cr *ContentRenderer) SetHighlightedNode(node dom.Node) {
	cr.highlightedNode = node
}
