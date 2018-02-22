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
	chromeContentRenderer *ContentRenderer
	backgroundRect        *dom.RectNode
	backButton            *dom.CircleNode
	urlBar                *dom.TextNode
	stateText             *dom.TextNode
	errorText             *dom.TextNode

	urlBarFocused bool
}

func NewBrowser(window *pixelgl.Window, initialURL string) *Browser {
	backButton := &dom.CircleNode{
		Radius: 7,
		Fill:   "blue",
	}
	backgroundRect := &dom.RectNode{
		Fill: "lightgrey",
	}
	urlBar := &dom.TextNode{}
	stateText := &dom.TextNode{}
	errorText := &dom.TextNode{}
	chromeGroup := &dom.GroupNode{
		TextNode: []*dom.TextNode{
			urlBar,
			stateText,
			errorText,
		},
		CircleNode: []*dom.CircleNode{backButton},
		RectNode:   []*dom.RectNode{
		//backgroundRect,
		},
	}

	b := &Browser{
		window: window,

		urlBar:         urlBar,
		backButton:     backButton,
		backgroundRect: backgroundRect,
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

	// Update background rect.
	b.backgroundRect.Width = b.window.Bounds().W()
	b.backgroundRect.X = 0
	b.backgroundRect.Y = b.window.Bounds().H() - 30
	b.backgroundRect.Height = 30

	// Update URL bar.
	urlToShow := b.newURL
	if urlToShow != b.currentPage.url {
		b.urlBar.Fill = "blue"
	} else {
		b.urlBar.Fill = "black"
	}
	b.urlBar.Value = urlToShow
	b.urlBar.X = 90
	b.urlBar.Y = b.window.Bounds().H() - 20

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

func (b *Browser) AddText(t string) {
	b.newURL += t
}

func (b *Browser) Backspace() {
	if b.newURL == "" {
		return
	}
	b.newURL = b.newURL[:len(b.newURL)-1]
}

func (b *Browser) NavigateToNewURL() {
	b.NavigateTo(b.newURL)
}

func (b *Browser) ToggleURLFocus() {
	b.urlBarFocused = !b.urlBarFocused
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
