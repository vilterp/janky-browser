package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/vilterp/janky-browser/package"
	"golang.org/x/image/colornames"
)

const initPage = "http://localhost:8084/circleRectText.svg"

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "JankyBrowser",
		Bounds:    pixel.R(0, 0, 1024, 768),
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	devtoolsCfg := pixelgl.WindowConfig{
		Title:     "Devtools | JankyBrowser",
		Bounds:    pixel.R(0, 0, 500, 500),
		Resizable: true,
	}
	devtoolsWin, err := pixelgl.NewWindow(devtoolsCfg)
	if err != nil {
		panic(err)
	}

	devtools := jankybrowser.NewDevtools(devtoolsWin)
	browser := jankybrowser.NewBrowser(win, initPage, devtools)

	fps := time.Tick(time.Second / 60)
	for !win.Closed() {
		win.Clear(colornames.White)
		devtoolsWin.Clear(colornames.White)

		// Handle mouse events.
		browser.ProcessMouseEvents(
			win.MousePosition(),
			win.Pressed(pixelgl.MouseButton1),
			win.JustPressed(pixelgl.MouseButton1),
		)

		// Handle keyboard events.
		typed := win.Typed()
		if len(typed) > 0 {
			browser.UrlInput.ProcessTyping(typed)
		}
		if win.JustReleased(pixelgl.KeyBackspace) || win.Repeated(pixelgl.KeyBackspace) {
			browser.UrlInput.ProcessBackspace()
		}
		if win.JustReleased(pixelgl.KeyEnter) {
			browser.UrlInput.ProcessEnter()
		}
		superDown := win.Pressed(pixelgl.KeyLeftSuper) || win.Pressed(pixelgl.KeyRightSuper)
		if win.JustReleased(pixelgl.KeyL) && superDown {
			browser.UrlInput.Focus()
		}
		if win.JustReleased(pixelgl.KeyTab) || win.JustReleased(pixelgl.KeyEscape) {
			if browser.UrlInput.Focused {
				browser.UrlInput.UnFocus()
			} else {
				browser.UnHighlightNode()
			}
		}
		shiftDown := win.Pressed(pixelgl.KeyLeftShift) || win.Pressed(pixelgl.KeyRightShift)
		if win.JustPressed(pixelgl.KeyLeft) || win.Repeated(pixelgl.KeyLeft) {
			if browser.UrlInput.Focused {
				browser.UrlInput.ProcessLeftKey(shiftDown, superDown)
			} else {
				browser.HighlightPrevNode()
			}
		}
		if win.JustPressed(pixelgl.KeyRight) || win.Repeated(pixelgl.KeyRight) {
			if browser.UrlInput.Focused {
				browser.UrlInput.ProcessRightKey(shiftDown, superDown)
			} else {
				browser.HighlightNextNode()
			}
		}

		// Draw.
		browser.Draw(win)
		win.Update()
		devtoolsWin.Update()

		<-fps
	}
}

func main() {
	pixelgl.Run(run)
}
