package jankybrowser

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type Browser struct {
	window      *pixelgl.Window
	currentPage *BrowserPage

	// Text for drawing URL
	txt *text.Text
}

func NewBrowser(window *pixelgl.Window, currentPage *BrowserPage) *Browser {
	face, err := loadTTF("Roboto-Regular.ttf", 15)
	if err != nil {
		panic(err)
	}
	atlas := text.NewAtlas(face, text.ASCII)
	txt := text.New(pixel.V(0, 0), atlas)
	txt.Color = colornames.Black

	return &Browser{
		currentPage: currentPage,
		window:      window,
		txt:         txt,
	}
}

func (b *Browser) Draw(t pixel.Target) {
	b.txt.Clear()
	b.txt.WriteString(b.currentPage.url)

	b.txt.Draw(t, pixel.IM.Moved(pixel.V(10, b.window.Bounds().H()-20.0)))

	imd := imdraw.New(nil)
	b.currentPage.Draw(imd)
	imd.Draw(t)
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
