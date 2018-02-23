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

	UrlInput *dom.TextInputNode

	backButton *dom.CircleNode
	stateText  *dom.TextNode
	errorText  *dom.TextNode
}

func NewBrowser(window *pixelgl.Window, initialURL string) *Browser {
	// TODO: maybe group chrome stuff into a custom element.
	// Initialize nodes.
	backButton := &dom.CircleNode{
		Radius: 7,
	}
	stateText := &dom.TextNode{}
	errorText := &dom.TextNode{}
	urlInput := &dom.TextInputNode{}

	// Group them so we can draw in one go.
	chromeGroup := &dom.GroupNode{
		TextNode: []*dom.TextNode{
			stateText,
			errorText,
		},
		CircleNode:    []*dom.CircleNode{backButton},
		TextInputNode: []*dom.TextInputNode{urlInput},
	}

	b := &Browser{
		window: window,

		// Save nodes so we can reference them.
		backButton: backButton,
		stateText:  stateText,
		errorText:  errorText,
		UrlInput:   urlInput,

		chromeContentRenderer: NewContentRenderer(chromeGroup),
	}

	b.UrlInput.OnEnter = func(newUrl string) {
		b.NavigateTo(newUrl)
		b.UrlInput.UnFocus()
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

	const urlBarStart = 85

	// Update URL input.
	if b.UrlInput.Value == b.currentPage.url {
		b.UrlInput.TextColor = "black"
	} else {
		b.UrlInput.TextColor = "blue"
	}
	b.UrlInput.Width = b.window.Bounds().W() - urlBarStart + 5
	b.UrlInput.X = urlBarStart
	b.UrlInput.Y = b.window.Bounds().H() - 30

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

func (b *Browser) NavigateTo(url string) {
	log.Println("navigate to", url)
	b.newURL = url
	b.currentPage = NewBrowserPage(url)
	b.currentPage.Load()
	b.UrlInput.SetValue(url)

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
