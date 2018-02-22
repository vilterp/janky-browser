package main

import (
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
