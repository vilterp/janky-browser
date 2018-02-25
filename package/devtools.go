package jankybrowser

import (
	"fmt"

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

func (dt *Devtools) Draw(bp *BrowserPage) {
	dt.drawDOM(bp)
	dt.renderer.Draw(dt.win)
}

const indent = dom.CharWidth * 2

func (dt *Devtools) drawDOM(bp *BrowserPage) {
	dt.domGroupNode.TextNode = nil

	if bp.state != PageStateLoaded {
		return
	}

	line := 0
	dom.Visit(
		bp.renderer.rootNode,
		func(n dom.Node, depth int) {
			dt.domGroupNode.TextNode = append(dt.domGroupNode.TextNode, &dom.TextNode{
				Y:     dt.win.Bounds().H() - float64((line+1)*dom.TextHeight),
				X:     float64(depth * indent),
				Value: dom.FormatWithoutChildren(n),
			})
			line++
		},
		func(n dom.Node, depth int) {
			if len(n.Children()) == 0 {
				return
			}
			dt.domGroupNode.TextNode = append(dt.domGroupNode.TextNode, &dom.TextNode{
				Y:     dt.win.Bounds().H() - float64((line+1)*dom.TextHeight),
				X:     float64(depth * indent),
				Value: fmt.Sprintf("</%s>", n.Name()),
			})
			line++
		},
	)
	dt.domGroupNode.Init()
}
