package jankybrowser

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/vilterp/jankybrowser/package/dom"
)

type Browser struct {
	window      *pixelgl.Window
	currentPage *BrowserPage

	// Text for drawing URL
	urlBar *dom.TextNode
}

func NewBrowser(window *pixelgl.Window, currentPage *BrowserPage) *Browser {
	b := &Browser{
		currentPage: currentPage,
		window:      window,
		urlBar:      &dom.TextNode{},
	}
	b.currentPage.Load()
	return b
}

func (b *Browser) Draw(t pixel.Target) {
	// Update & draw URL bar.
	b.currentPage.mu.RLock()
	str := fmt.Sprintf("%s | %s", StateNames[b.currentPage.state], b.currentPage.url)
	if b.currentPage.state == PageStateError {
		str = fmt.Sprintf("%s | %s", str, b.currentPage.loadError.Error())
	}
	b.currentPage.mu.RUnlock()

	b.urlBar.Value = str
	b.urlBar.X = 10
	b.urlBar.Y = b.window.Bounds().H() - 20
	b.urlBar.Draw(t)

	// Draw page.
	b.currentPage.Draw(t)
}

func (b *Browser) ProcessMouseEvents(pt pixel.Vec) {
	b.currentPage.ProcessMouseEvents(pt)
}

type PageState = int

const (
	PageStateInit PageState = iota
	PageStateLoading
	PageStateLoaded
	PageStateError
)

var StateNames = map[PageState]string{
	PageStateInit:    "INIT",
	PageStateLoading: "LOADING",
	PageStateLoaded:  "LOADED",
	PageStateError:   "ERROR",
}

type BrowserPage struct {
	mu sync.RWMutex

	url       string
	state     PageState
	loadError error       // set when state = PageStateError
	rootNode  dom.DOMNode // set when state = PageStateLoaded
}

func NewBrowserPage(url string) *BrowserPage {
	return &BrowserPage{
		state: PageStateInit,
		url:   url,
	}
}

func (bp *BrowserPage) Load() {
	go func() {
		bp.doLoad()
	}()
}

// doLoad is blocking. Don't call in the main UI thread.
func (bp *BrowserPage) doLoad() {
	bp.mu.Lock()
	bp.state = PageStateLoading
	bp.mu.Unlock()

	response, err := http.Get(bp.url)

	bp.mu.Lock()
	defer bp.mu.Unlock()

	if err != nil {
		bp.state = PageStateError
		bp.loadError = err
		return
	}
	if response.StatusCode != 200 {
		bp.state = PageStateError
		// TODO: structured error
		bp.loadError = fmt.Errorf("non-200 status code: %d", response.StatusCode)
		return
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		bp.state = PageStateError
		bp.loadError = fmt.Errorf("error reading input stream: %s", err.Error())
		return
	}

	node, err := dom.Parse(bytes)
	if err != nil {
		bp.state = PageStateError
		// TODO: structured error
		bp.loadError = fmt.Errorf("parse error: %s", err.Error())
		return
	}

	bp.state = PageStateLoaded
	if node == nil {
		node = &dom.GroupNode{}
	}
	bp.rootNode = node
	log.Println("rootNode:", dom.Format(bp.rootNode))
}

func (bp *BrowserPage) Draw(t pixel.Target) {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	switch bp.state {
	case PageStateInit:
		break
	case PageStateLoading:
		break
	case PageStateLoaded:
		bp.rootNode.Draw(t)
	case PageStateError:
		// TODO: render error state
		// make a <text> element and an error DOM and use it!
	}
}

func (bp *BrowserPage) ProcessMouseEvents(pt pixel.Vec) {
	node := bp.GetHoveredNode(pt)
	if node != nil {
		log.Println("hovering over", node)
	}
}

func (bp *BrowserPage) GetHoveredNode(pt pixel.Vec) dom.DOMNode {
	switch bp.state {
	case PageStateLoaded:
		bp.mu.RLock()
		defer bp.mu.RUnlock()
		return dom.Pick(bp.rootNode, pt)
	default:
		return nil
	}
}
