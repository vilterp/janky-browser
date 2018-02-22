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

	history []string

	// Text for drawing URL
	backButton *dom.CircleNode
	urlBar     *dom.TextNode
}

func NewBrowser(window *pixelgl.Window, initialURL string) *Browser {
	b := &Browser{
		window: window,
		urlBar: &dom.TextNode{},
		backButton: &dom.CircleNode{
			Radius: 7,
			Fill:   "blue",
		},
	}
	b.urlBar.Init()
	b.NavigateTo(initialURL)
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
	b.urlBar.X = 35
	b.urlBar.Y = b.window.Bounds().H() - 20
	b.urlBar.Draw(t)

	// Draw back button.
	if len(b.history) > 1 {
		b.backButton.Fill = "lightblue"
	} else {
		b.backButton.Fill = "grey"
	}
	b.backButton.X = 20
	b.backButton.Y = b.window.Bounds().H() - 15
	b.backButton.Draw(t)

	// Draw page.
	b.currentPage.Draw(t)
}

func (b *Browser) ProcessMouseEvents(pt pixel.Vec) {
	b.currentPage.mu.RLock()
	defer b.currentPage.mu.RUnlock()

	if len(b.history) > 1 && b.currentPage.state != PageStateLoading && b.backButton.Contains(pt) {
		b.NavigateBack()
		return
	}

	navigateTo := b.currentPage.ProcessMouseEvents(pt)
	if navigateTo != "" {
		b.NavigateTo(navigateTo)
	}
}

func (b *Browser) NavigateTo(url string) {
	log.Println("navigate to", url)
	b.currentPage = NewBrowserPage(url)
	b.currentPage.Load()

	b.history = append(b.history, url)
}

func (b *Browser) NavigateBack() error {
	url, err := b.PopHistory()
	if err != nil {
		return err
	}
	b.NavigateTo(url)
	return nil
}

func (b *Browser) PopHistory() (string, error) {
	if len(b.history) == 1 {
		return "", fmt.Errorf("can't go back; already on last page")
	}

	url := b.history[len(b.history)-2]
	b.history = b.history[:len(b.history)-2]
	return url, nil
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
	loadError error    // set when state = PageStateError
	rootNode  dom.Node // set when state = PageStateLoaded
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
	bp.rootNode.Init()
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

func (bp *BrowserPage) ProcessMouseEvents(pt pixel.Vec) string {
	hoveredNodes := bp.GetHoveredNodes(pt)
	var navigateTo string
	if len(hoveredNodes) > 0 {
		//var hoveredNodeStrs []string
		//for _, hoveredNode := range hoveredNodes {
		//	hoveredNodeStrs = append(hoveredNodeStrs, dom.Format(hoveredNode))
		//}
		for _, hoveredNode := range hoveredNodes {
			switch n := hoveredNode.(type) {
			case *dom.GroupNode:
				if n.Href != "" {
					navigateTo = n.Href
				}
			}
		}
	}
	return navigateTo
}

func (bp *BrowserPage) GetHoveredNodes(pt pixel.Vec) []dom.Node {
	switch bp.state {
	case PageStateLoaded:
		bp.mu.RLock()
		defer bp.mu.RUnlock()
		return dom.Pick(bp.rootNode, pt)
	default:
		return nil
	}
}
