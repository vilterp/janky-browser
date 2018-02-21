package jankybrowser

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/vilterp/jankybrowser/package/dom"
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

	b := &Browser{
		currentPage: currentPage,
		window:      window,
		txt:         txt,
	}
	b.currentPage.Load()
	return b
}

func (b *Browser) Draw(t pixel.Target) {
	b.txt.Clear()

	b.currentPage.mu.RLock()
	str := fmt.Sprintf("%s | %s", StateNames[b.currentPage.state], b.currentPage.url)
	if b.currentPage.state == PageStateError {
		str = fmt.Sprintf("%s | %s", str, b.currentPage.loadError.Error())
	}
	b.txt.WriteString(str)
	b.currentPage.mu.RUnlock()

	b.txt.Draw(t, pixel.IM.Moved(pixel.V(10, b.window.Bounds().H()-20.0)))

	imd := imdraw.New(nil)
	b.currentPage.Draw(imd)
	imd.Draw(t)
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

func (bp *BrowserPage) Draw(imd *imdraw.IMDraw) {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	switch bp.state {
	case PageStateInit:
		break
	case PageStateLoading:
		break
	case PageStateLoaded:
		bp.rootNode.Draw(imd)
	case PageStateError:
		// TODO: render error state
		// make a <text> element and an error DOM and use it!
	}
}
