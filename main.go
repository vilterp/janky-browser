package main

import (
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/vilterp/jankybrowser/package"
	"golang.org/x/image/colornames"
)

const initPage = "http://localhost:8081/circleRectText.svg"

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

	browser := jankybrowser.NewBrowser(win, initPage)

	fps := time.Tick(time.Second / 60)
	for !win.Closed() {
		win.Clear(colornames.White)

		browser.ProcessMouseEvents(
			win.MousePosition(),
			win.Pressed(pixelgl.MouseButton1),
			win.JustPressed(pixelgl.MouseButton1),
		)

		typed := win.Typed()
		if len(typed) > 0 {
			log.Println("typed", typed)
			browser.ProcessTyping(typed)
		}
		if win.JustReleased(pixelgl.KeyBackspace) {
			browser.ProcessBackspace()
		}
		if win.JustReleased(pixelgl.KeyEnter) {
			browser.ProcessEnter()
		}
		if win.JustReleased(pixelgl.KeyL) && win.Pressed(pixelgl.KeyLeftSuper) {
			browser.FocusURLBar()
		}
		if win.JustReleased(pixelgl.KeyTab) || win.JustReleased(pixelgl.KeyEscape) {
			browser.UnFocusURLBar()
		}

		// TODO: handle clicks, not just position
		// TODO: handle keyboard events

		browser.Draw(win)
		win.Update()

		<-fps
	}
}

func main() {
	pixelgl.Run(run)
}
