package jankybrowser

import (
	"fmt"
	"image"
	"log"
	"net/url"

	"github.com/llgcode/draw2d"
	"github.com/vilterp/janky-browser/package/dom"
	"github.com/vilterp/janky-browser/package/util"
)

type Browser struct {
	window      *util.Window
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

func NewBrowser(
	window *util.Window, initialURL string,
) *Browser {
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

func (b *Browser) Draw(gc draw2d.GraphicContext) {
	// Update & draw URL bar.
	b.drawChrome(gc)

	// Draw page.
	b.currentPage.Draw(gc)
}

// TODO: factor this out into its own DOMNode/Component which takes its own attributes
// and emits its own events... once we have those concepts...
func (b *Browser) drawChrome(gc draw2d.GraphicContext) {
	b.currentPage.mu.RLock()
	defer b.currentPage.mu.RUnlock()

	const urlBarStart = 170

	// Update URL input.
	if b.UrlInput.Value == b.currentPage.url {
		b.UrlInput.TextColor = "black"
	} else {
		b.UrlInput.TextColor = "blue"
	}
	b.UrlInput.Width = float64(b.window.Size.X - urlBarStart + 5)
	b.UrlInput.X = urlBarStart
	b.UrlInput.Y = 30

	// Update status text.
	b.stateText.Value = StateNames[b.currentPage.state]
	b.stateText.X = 80
	b.stateText.Y = 30

	// Update back button.
	if len(b.history) > 1 {
		b.backButton.Fill = "blue"
	} else {
		b.backButton.Fill = "grey"
	}
	b.backButton.X = 10
	b.backButton.Y = 30

	// Update error text.
	errorText := ""
	if b.currentPage.state == PageStateError {
		errorText = b.currentPage.loadError.Error()
	}
	b.errorText.Fill = "red"
	b.errorText.Value = errorText
	b.errorText.X = 10
	b.errorText.Y = 80

	b.chromeContentRenderer.Draw(gc)
}

func (b *Browser) ProcessMouseEvents(pt image.Point, mouseDown bool, mouseJustDown bool) {
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
