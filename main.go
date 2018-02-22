package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/vilterp/jankybrowser/package"
	"golang.org/x/image/colornames"
)

const initPage = "http://localhost:8081/circleAndRect.svg"

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

	page := jankybrowser.NewBrowserPage(initPage)
	browser := jankybrowser.NewBrowser(win, page)

	fps := time.Tick(time.Second / 60)
	for !win.Closed() {
		win.Clear(colornames.White)

		browser.Draw(win)
		win.Update()

		<-fps
	}
}

func main() {
	pixelgl.Run(run)
}
