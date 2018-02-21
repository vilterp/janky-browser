package jankybrowser

import (
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Browser struct {
	window      *pixelgl.Window
	currentPage *BrowserPage
}

func NewBrowser(window *pixelgl.Window, currentPage *BrowserPage) *Browser {
	return &Browser{
		currentPage: currentPage,
		window:      window,
	}
}

func (b *Browser) Draw(imd *imdraw.IMDraw) {
	b.currentPage.Draw(imd)
}

type BrowserPage struct {
	url      string
	rootNode DOMNode
}

func NewBrowserPage(url string, rootNode DOMNode) *BrowserPage {
	return &BrowserPage{
		rootNode: rootNode,
		url:      url,
	}
}

func (bp *BrowserPage) Draw(imd *imdraw.IMDraw) {
	bp.rootNode.Draw(imd)
}

func ExampleDOMTree() DOMNode {
	return &GroupNode{
		children: []DOMNode{
			&CircleNode{
				fill:   colornames.Red,
				radius: 50,
				x:      200,
				y:      300,
			},
			&RectNode{
				fill:   colornames.Blue,
				x:      20,
				y:      30,
				width:  500,
				height: 100,
			},
		},
	}
}
