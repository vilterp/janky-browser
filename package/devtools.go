package jankybrowser

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/vilterp/janky-browser/package/dom"
)

type Devtools struct {
	win *pixelgl.Window

	renderer   *ContentRenderer
	domStrNode *dom.TextNode
}

func NewDevtools(win *pixelgl.Window) *Devtools {
	domStrNode := &dom.TextNode{}
	rootGroup := &dom.GroupNode{
		TextNode: []*dom.TextNode{
			domStrNode,
		},
	}

	return &Devtools{
		win:        win,
		renderer:   NewContentRenderer(rootGroup),
		domStrNode: domStrNode,
	}
}

func (dt *Devtools) Draw(bp *BrowserPage) {
	if bp.state != PageStateLoaded {
		dt.domStrNode.Value = ""
	} else {
		dt.domStrNode.Value = dom.Format(bp.renderer.rootNode)
	}
	dt.domStrNode.Y = dt.win.Bounds().H() - 10

	dt.renderer.Draw(dt.win)
}
