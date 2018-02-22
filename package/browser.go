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
	b.DrawChrome(t)

	// Draw page.
	b.currentPage.Draw(t)
}

// TODO: factor this out into its own DOMNode/Component which takes its own attributes
// and emits its own events... once we have those concepts...
func (b *Browser) DrawChrome(t pixel.Target) {
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
}

func (b *Browser) ProcessMouseEvents(pt pixel.Vec, mouseDown bool, mouseJustDown bool) {
	b.currentPage.mu.RLock()
	defer b.currentPage.mu.RUnlock()

	if len(b.history) > 1 && b.currentPage.state != PageStateLoading && b.backButton.Contains(pt) {
		b.NavigateBack()
		return
	}

	navigateTo := b.currentPage.ProcessMouseEvents(pt, mouseDown, mouseJustDown)
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
