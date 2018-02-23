package jankybrowser

import (
	"fmt"
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/vilterp/jankybrowser/package/dom"
)

type Browser struct {
	window      *pixelgl.Window
	currentPage *BrowserPage

	history []string

	newURL string

	// Stuff for drawing the chrome.
	// TODO: wrap this up in its own struct somehow.
	chromeContentRenderer *ContentRenderer

	// TODO: factor these out into text field
	urlBar         *dom.TextNode
	urlCursor      *dom.LineNode
	backgroundRect *dom.RectNode
	selectionRect  *dom.RectNode
	urlBarFocused  bool
	cursorPos      int
	selectionStart *int

	backButton *dom.CircleNode
	stateText  *dom.TextNode
	errorText  *dom.TextNode
}

func NewBrowser(window *pixelgl.Window, initialURL string) *Browser {
	backButton := &dom.CircleNode{
		Radius: 7,
	}
	backgroundRect := &dom.RectNode{}
	selectionRect := &dom.RectNode{}
	urlBar := &dom.TextNode{}
	stateText := &dom.TextNode{}
	errorText := &dom.TextNode{}
	urlCursor := &dom.LineNode{}
	chromeGroup := &dom.GroupNode{
		TextNode: []*dom.TextNode{
			urlBar,
			stateText,
			errorText,
		},
		CircleNode: []*dom.CircleNode{backButton},
		RectNode: []*dom.RectNode{
			backgroundRect,
			selectionRect,
		},
		LineNode: []*dom.LineNode{
			urlCursor,
		},
	}

	b := &Browser{
		window: window,

		urlBar:         urlBar,
		urlCursor:      urlCursor,
		backButton:     backButton,
		backgroundRect: backgroundRect,
		selectionRect:  selectionRect,
		stateText:      stateText,
		errorText:      errorText,

		chromeContentRenderer: NewContentRenderer(chromeGroup),
	}
	b.NavigateTo(initialURL)
	return b
}

func (b *Browser) Draw(t pixel.Target) {
	// Update & draw URL bar.
	b.DrawChrome(t)

	// Draw page.
	b.currentPage.Draw(t)
}

// TODO: factor this out into its own DOMNode/Component which takes its own attributes
// and emits its own events... once we have those concepts...
func (b *Browser) DrawChrome(t pixel.Target) {
	b.currentPage.mu.RLock()
	defer b.currentPage.mu.RUnlock()

	const urlBarStart = 90

	// Update background rect.
	b.backgroundRect.Width = b.window.Bounds().W() - urlBarStart + 5
	b.backgroundRect.X = urlBarStart - 5
	b.backgroundRect.Y = b.window.Bounds().H() - 30
	b.backgroundRect.Height = 30
	if b.urlBarFocused {
		b.backgroundRect.Fill = "lightgrey"
	} else {
		b.backgroundRect.Fill = "white"
	}

	// Update URL bar.
	urlToShow := b.newURL
	if urlToShow != b.currentPage.url {
		b.urlBar.Fill = "blue"
	} else {
		b.urlBar.Fill = "black"
	}
	b.urlBar.Value = urlToShow
	b.urlBar.X = urlBarStart
	b.urlBar.Y = b.window.Bounds().H() - 20

	// Update cursor.
	const charWidth = float64(7)
	cursorX := float64(b.cursorPos)*charWidth + urlBarStart
	b.urlCursor.X1 = cursorX
	b.urlCursor.X2 = cursorX
	b.urlCursor.Y1 = b.window.Bounds().H() - 8
	b.urlCursor.Y2 = b.window.Bounds().H() - 22
	if b.urlBarFocused {
		b.urlCursor.Stroke = "red"
	} else {
		b.urlCursor.Stroke = ""
	}

	// Update selection.

	if b.selectionStart == nil {
		b.selectionRect.Fill = ""
	} else {
		b.selectionRect.Fill = "pink"
		startIdx, endIdx := b.GetSelection()
		b.selectionRect.X = float64(startIdx)*charWidth + urlBarStart
		b.selectionRect.Y = b.window.Bounds().H() - 23
		b.selectionRect.Width = float64(endIdx-startIdx) * charWidth
		b.selectionRect.Height = 13
	}

	// Update status text.
	b.stateText.Value = StateNames[b.currentPage.state]
	b.stateText.X = 35
	b.stateText.Y = b.window.Bounds().H() - 20

	// Update error text.
	errorText := ""
	if b.currentPage.state == PageStateError {
		errorText = b.currentPage.loadError.Error()
	}
	b.errorText.Value = errorText
	b.errorText.X = 20
	b.errorText.Y = b.window.Bounds().H() - 50

	// Update back button.
	if len(b.history) > 1 {
		b.backButton.Fill = "lightblue"
	} else {
		b.backButton.Fill = "grey"
	}
	b.backButton.X = 20
	b.backButton.Y = b.window.Bounds().H() - 15

	b.chromeContentRenderer.Draw(t)
}

func (b *Browser) ProcessMouseEvents(pt pixel.Vec, mouseDown bool, mouseJustDown bool) {
	b.currentPage.mu.RLock()
	defer b.currentPage.mu.RUnlock()

	clickedNodes := b.chromeContentRenderer.processClickState(pt, mouseDown, mouseJustDown)
	if len(clickedNodes) > 0 && clickedNodes[0] == b.backButton {
		if len(b.history) > 1 && b.currentPage.state != PageStateLoading {
			b.NavigateBack()
			return
		}
	}

	navigateTo := b.currentPage.ProcessMouseEvents(pt, mouseDown, mouseJustDown)
	if navigateTo != "" {
		b.NavigateTo(navigateTo)
	}
}

func (b *Browser) ProcessTyping(t string) {
	if !b.urlBarFocused {
		return
	}
	if b.selectionStart != nil {
		b.DeleteSelection()
	}
	b.newURL = b.newURL[:b.cursorPos] + t + b.newURL[b.cursorPos:]
	b.cursorPos += 1
}

func (b *Browser) ProcessBackspace() {
	if !b.urlBarFocused {
		return
	}
	if b.newURL == "" {
		return
	}
	if b.selectionStart == nil {
		b.newURL = b.newURL[:b.cursorPos-1] + b.newURL[b.cursorPos:]
		b.cursorPos -= 1
	} else {
		b.DeleteSelection()
	}
}

func (b *Browser) DeleteSelection() {
	startIdx, endIdx := b.GetSelection()
	b.newURL = b.newURL[:startIdx] + b.newURL[endIdx:]
	b.CancelSelection()
}

func (b *Browser) GetSelection() (int, int) {
	if b.selectionStart == nil {
		log.Println("nil selection start")
		return b.cursorPos, b.cursorPos
	}
	selectionStart := *b.selectionStart
	var startIdx int
	var endIdx int
	if selectionStart < b.cursorPos {
		startIdx = selectionStart
		endIdx = b.cursorPos
	} else {
		endIdx = selectionStart
		startIdx = b.cursorPos
	}
	return startIdx, endIdx
}

func (b *Browser) ProcessEnter() {
	if !b.urlBarFocused {
		return
	}
	b.NavigateTo(b.newURL)
	b.urlBarFocused = false
}

func (b *Browser) FocusURLBar() {
	b.urlBarFocused = true
	b.cursorPos = len(b.newURL)
}

func (b *Browser) UnFocusURLBar() {
	b.urlBarFocused = false
}

func (b *Browser) ProcessLeftKey(shiftDown bool, superDown bool) {
	if !b.urlBarFocused {
		return
	}
	b.MaybeStartSelection(shiftDown)
	if superDown {
		b.cursorPos = 0
		return
	}
	b.cursorPos = b.cursorPos - 1
	if b.cursorPos < 0 {
		b.cursorPos = 0
	}
}

func (b *Browser) ProcessRightKey(shiftDown bool, superDown bool) {
	if !b.urlBarFocused {
		return
	}
	b.MaybeStartSelection(shiftDown)
	if superDown {
		b.cursorPos = len(b.newURL)
		return
	}
	b.cursorPos = b.cursorPos + 1
	if b.cursorPos > len(b.newURL) {
		b.cursorPos = len(b.newURL)
	}
}

func (b *Browser) MaybeStartSelection(shiftDown bool) {
	if !shiftDown {
		b.CancelSelection()
		return
	}
	if b.selectionStart != nil {
		return
	}
	selectionStart := b.cursorPos
	b.selectionStart = &selectionStart
}

func (b *Browser) CancelSelection() {
	b.selectionStart = nil
}

func (b *Browser) NavigateTo(url string) {
	log.Println("navigate to", url)
	b.newURL = url
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
