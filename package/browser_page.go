package jankybrowser

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/faiface/pixel"
	"github.com/vilterp/janky-browser/package/dom"
)

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
	loadError error // set when state = PageStateError

	renderer *ContentRenderer
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
	bp.renderer = NewContentRenderer(node)
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
		bp.renderer.Draw(t)
	case PageStateError:
		// TODO: render error state
		// make a <text> element and an error DOM and use it!
	}
}

func (bp *BrowserPage) numNodes() int {
	if bp.state != PageStateLoaded {
		return 0
	}
	return len(dom.GetAllNodes(bp.renderer.rootNode))
}

func (bp *BrowserPage) ProcessMouseEvents(pt pixel.Vec, mouseDown bool, mouseJustDown bool) string {
	if bp.state != PageStateLoaded {
		return ""
	}

	bp.mu.RLock()
	defer bp.mu.RUnlock()

	clickedNodes := bp.renderer.processClickState(pt, mouseDown, mouseJustDown)

	var navigateTo string
	if len(clickedNodes) > 0 {
		for _, hoveredNode := range clickedNodes {
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
