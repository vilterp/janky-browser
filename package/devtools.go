package jankybrowser

import (
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

func (dt *Devtools) drawDOM(bp *BrowserPage) {
	domStrNode := &dom.TextNode{}
	domStrNode.Init()
	if bp.state != PageStateLoaded {
		domStrNode.Value = ""
	} else {
		domStrNode.Value = dom.Format(bp.renderer.rootNode)
	}
	domStrNode.Y = dt.win.Bounds().H() - 10
	dt.domGroupNode.TextNode = []*dom.TextNode{domStrNode}
}
