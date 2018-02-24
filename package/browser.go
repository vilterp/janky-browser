package jankybrowser

import (
	"fmt"
	"log"
	"net/url"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/vilterp/janky-browser/package/dom"
)

type Browser struct {
	window      *pixelgl.Window
	currentPage *BrowserPage

	history []string

	// Stuff for drawing the chrome.
	// TODO: wrap this up in its own struct somehow.
	chromeContentRenderer *ContentRenderer

	UrlInput *dom.TextInputNode

	backButton *dom.TextNode
	stateText  *dom.TextNode
	errorText  *dom.TextNode
}

func NewBrowser(window *pixelgl.Window, initialURL string) *Browser {
	// TODO: maybe group chrome stuff into a custom element.
	// Initialize nodes.
	backButton := &dom.TextNode{
		Value: "BACK",
	}
	stateText := &dom.TextNode{}
	errorText := &dom.TextNode{}
	urlInput := &dom.TextInputNode{}

	// Group them so we can draw in one go.
	chromeGroup := &dom.GroupNode{
		TextNode: []*dom.TextNode{
			stateText,
			errorText,
			backButton,
		},
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

	const urlBarStart = 90

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
	b.stateText.X = 45
	b.stateText.Y = b.window.Bounds().H() - 20

	// Update back button.
	if len(b.history) > 1 {
		b.backButton.Fill = "blue"
	} else {
		b.backButton.Fill = "grey"
	}
	b.backButton.X = 10
	b.backButton.Y = b.window.Bounds().H() - 20

	// Update error text.
	errorText := ""
	if b.currentPage.state == PageStateError {
		errorText = b.currentPage.loadError.Error()
	}
	b.errorText.Value = errorText
	b.errorText.X = 20
	b.errorText.Y = b.window.Bounds().H() - 50

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
		b.NavigateTo(b.resolveURL(navigateTo))
	}
}

func (b *Browser) resolveURL(unresolvedURL string) string {
	parsed, err := url.Parse(unresolvedURL)
	if err != nil {
		return unresolvedURL
	}
	if parsed.Host == "" {
		currentURL, _ := url.Parse(b.currentPage.url)
		parsed.Scheme = currentURL.Scheme
		parsed.Host = currentURL.Host
	}
	return parsed.String()
}

func (b *Browser) NavigateTo(newURL string) {
	log.Println("navigate to", newURL)
	b.currentPage = NewBrowserPage(newURL)
	b.currentPage.Load()
	b.UrlInput.Value = newURL

	b.history = append(b.history, newURL)
}

func (b *Browser) NavigateBack() error {
	toURL, err := b.popHistory()
	if err != nil {
		return err
	}
	b.NavigateTo(toURL)
	return nil
}

func (b *Browser) popHistory() (string, error) {
	if len(b.history) == 1 {
		return "", fmt.Errorf("can't go back; already on last page")
	}

	popped := b.history[len(b.history)-2]
	b.history = b.history[:len(b.history)-2]
	return popped, nil
}

func (b *Browser) UnHighlightNode() {
	b.currentPage.UnHighlightNode()
}

func (b *Browser) HighlightNextNode() {
	b.currentPage.HighlightNextNode()
}

func (b *Browser) HighlightPrevNode() {
	b.currentPage.HighlightPrevNode()
}
