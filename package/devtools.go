package jankybrowser

import (
	"fmt"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/vilterp/janky-browser/package/dom"
)

type Devtools struct {
	win *pixelgl.Window

	renderer     *ContentRenderer
	domGroupNode *dom.GroupNode
}

func NewDevtools(win *pixelgl.Window) *Devtools {
	domGroup := &dom.GroupNode{}
	rootGroup := &dom.GroupNode{
		GroupNode: []*dom.GroupNode{
			domGroup,
		},
	}

	return &Devtools{
		win:          win,
		renderer:     NewContentRenderer(rootGroup),
		domGroupNode: domGroup,
	}
}

func (dt *Devtools) ProcessMouseEvents(pt pixel.Vec, mouseDown bool, mouseJustDown bool) {
	dt.renderer.processClickState(pt, mouseDown, mouseJustDown)
}

func (dt *Devtools) Draw(bp *BrowserPage) {
	dt.drawDOM(bp)
	dt.renderer.Draw(dt.win)
}

func (dt *Devtools) drawDOM(bp *BrowserPage) {
	dt.domGroupNode.TextNode = nil

	if bp.state != PageStateLoaded {
		return
	}

	// TODO: figure out a way to DRY this up...
	line := 0
	dom.Visit(
		bp.renderer.rootNode,
		func(n dom.Node, depth int) {
			indent := strings.Repeat("  ", depth)
			textNode := &dom.TextNode{
				Y:     dt.win.Bounds().H() - float64((line+1)*dom.TextHeight),
				Value: fmt.Sprintf("%s%s", indent, dom.FormatWithoutChildren(n)),
			}
			if bp.renderer.highlightedNode == n {
				textNode.Fill = "red"
			}
			textNode.Events().OnMouseOver = func() {
				bp.renderer.SetHighlightedNode(n)
			}
			textNode.Events().OnMouseOut = func() {
				bp.renderer.SetHighlightedNode(nil)
			}
			dt.domGroupNode.TextNode = append(dt.domGroupNode.TextNode, textNode)
			line++
		},
		func(n dom.Node, depth int) {
			if len(n.Children()) == 0 {
				return
			}
			indent := strings.Repeat("  ", depth)
			textNode := &dom.TextNode{
				Y:     dt.win.Bounds().H() - float64((line+1)*dom.TextHeight),
				Value: fmt.Sprintf("%s</%s>", indent, n.Name()),
			}
			if bp.renderer.highlightedNode == n {
				textNode.Fill = "red"
			}
			textNode.Events().OnMouseOver = func() {
				bp.renderer.SetHighlightedNode(n)
			}
			textNode.Events().OnMouseOut = func() {
				bp.renderer.SetHighlightedNode(nil)
			}
			dt.domGroupNode.TextNode = append(dt.domGroupNode.TextNode, textNode)
			line++
		},
	)
	dt.domGroupNode.Init()
}
