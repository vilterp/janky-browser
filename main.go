package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
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

	fps := time.Tick(time.Second / 120)
	for !win.Closed() {
		win.Clear(colornames.White)
		fmt.Println("draw")
		<-fps
	}
}

func main() {
	pixelgl.Run(run)
}
