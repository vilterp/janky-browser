package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	jankybrowser "github.com/vilterp/jankybrowser/package"
	"golang.org/x/image/colornames"
)

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

	page := jankybrowser.NewBrowserPage("http://example.com/", jankybrowser.ExampleDOMTree())
	browser := jankybrowser.NewBrowser(win, page)

	fps := time.Tick(time.Second / 120)
	for !win.Closed() {
		win.Clear(colornames.White)

		imd := imdraw.New(nil)
		browser.Draw(imd)
		imd.Draw(win)
		win.Update()

		<-fps
	}
}

func main() {
	pixelgl.Run(run)
}
