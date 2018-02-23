package jankybrowser

import (
	"reflect"

	"github.com/faiface/pixel"
	"github.com/vilterp/janky-browser/package/dom"
)

type ContentRenderer struct {
	rootNode dom.Node // set when state = PageStateLoaded

	// the set of nodes the mouse was over when it was pressed.
	// empty if the mouse has not been pressed.
	mouseDownNodes []dom.Node
}

func NewContentRenderer(rootNode dom.Node) *ContentRenderer {
	cr := &ContentRenderer{
		rootNode: rootNode,
	}
	cr.rootNode.Init()
	return cr
}

// processClickState steps the click state machine, returning clicked nodes if there are any.
func (cr *ContentRenderer) processClickState(
	pt pixel.Vec, mouseDown bool, mouseJustDown bool,
) []dom.Node {
	var res []dom.Node
	hoveredNodes := cr.GetHoveredNodes(pt)

	if mouseJustDown {
		cr.mouseDownNodes = hoveredNodes
	} else if !mouseDown && len(cr.mouseDownNodes) > 0 {
		if reflect.DeepEqual(hoveredNodes, cr.mouseDownNodes) {
			res = cr.mouseDownNodes
		}
		cr.mouseDownNodes = nil
	}
	return res
}

func (cr *ContentRenderer) GetHoveredNodes(pt pixel.Vec) []dom.Node {
	return dom.Pick(cr.rootNode, pt)
}

func (cr *ContentRenderer) Draw(t pixel.Target) {
	cr.rootNode.Draw(t)
}
