package jankybrowser

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"sync"

	"github.com/faiface/pixel"
	"github.com/vilterp/jankybrowser/package/dom"
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
	loadError error    // set when state = PageStateError
	rootNode  dom.Node // set when state = PageStateLoaded

	// the set of nodes the mouse was over when it was pressed.
	// empty if the mouse has not been pressed.
	mouseDownNodes []dom.Node
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
	log.Println("parsed DOM tree:", dom.Format(bp.rootNode))
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

func (bp *BrowserPage) ProcessMouseEvents(pt pixel.Vec, mouseDown bool, mouseJustDown bool) string {
	clickedNodes := bp.processClickState(pt, mouseDown, mouseJustDown)

	var navigateTo string
	if len(clickedNodes) > 0 {
		//var hoveredNodeStrs []string
		//for _, hoveredNode := range hoveredNodes {
		//	hoveredNodeStrs = append(hoveredNodeStrs, dom.Format(hoveredNode))
		//}
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

// processClickState steps the click state machine, returning clicked nodes if there are any.
func (bp *BrowserPage) processClickState(
	pt pixel.Vec, mouseDown bool, mouseJustDown bool,
) []dom.Node {
	var res []dom.Node
	hoveredNodes := bp.GetHoveredNodes(pt)

	if mouseJustDown {
		log.Println("begin click on", hoveredNodes)
		bp.mouseDownNodes = hoveredNodes
	} else if !mouseDown && len(bp.mouseDownNodes) > 0 {
		if reflect.DeepEqual(hoveredNodes, bp.mouseDownNodes) {
			log.Println("clicked on", bp.mouseDownNodes)
			res = bp.mouseDownNodes
		}
		bp.mouseDownNodes = nil
	}
	return res
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
