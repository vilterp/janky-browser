package jankybrowser

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/faiface/pixel"
	"github.com/vilterp/janky-browser/package/dom"
	"github.com/vilterp/janky-browser/package/util"
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

	renderer           *ContentRenderer
	highlightedNodeIdx *int
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
		if bp.highlightedNodeIdx == nil {
			bp.renderer.highlightedNode = nil
		} else {
			allNodes := dom.GetAllNodes(bp.renderer.rootNode)
			bp.renderer.highlightedNode = allNodes[*bp.highlightedNodeIdx]
		}
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

func (bp *BrowserPage) HighlightNextNode() {
	bp.changeHighlightedNodeIdx(1)
}

func (bp *BrowserPage) HighlightPrevNode() {
	bp.changeHighlightedNodeIdx(-1)
}

func (bp *BrowserPage) UnHighlightNode() {
	bp.highlightedNodeIdx = nil
}

func (bp *BrowserPage) changeHighlightedNodeIdx(by int) {
	if bp.highlightedNodeIdx == nil {
		zero := 0
		bp.highlightedNodeIdx = &zero
		return
	}
	numNodes := bp.numNodes()
	newIdx := util.Clamp(0, numNodes-1, *bp.highlightedNodeIdx+by)
	bp.highlightedNodeIdx = &newIdx
	log.Println("highlight node:", *bp.highlightedNodeIdx)
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
