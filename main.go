package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/vilterp/janky-browser/package"
	"golang.org/x/image/colornames"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

const initPage = "http://localhost:8084/circleRectText.svg"

func run() {
	// Conditionally initialize profiler.
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// Initialize windows.
	cfg := pixelgl.WindowConfig{
		Title:     "JankyBrowser",
		Bounds:    pixel.R(0, 0, 1024, 768),
		Resizable: true,
		VSync:     true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	devtoolsCfg := pixelgl.WindowConfig{
		Title:     "Devtools | JankyBrowser",
		Bounds:    pixel.R(0, 0, 500, 500),
		Resizable: true,
		VSync:     true,
	}
	devtoolsWin, err := pixelgl.NewWindow(devtoolsCfg)
	if err != nil {
		panic(err)
	}

	// Initialize browser and devtools.
	devtools := jankybrowser.NewDevtools(devtoolsWin)
	browser := jankybrowser.NewBrowser(win, initPage, devtools)

	// Main loop.
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
		if win.JustPressed(pixelgl.KeyBackspace) || win.Repeated(pixelgl.KeyBackspace) {
			browser.UrlInput.ProcessBackspace()
		}
		if win.JustPressed(pixelgl.KeyEnter) {
			browser.UrlInput.ProcessEnter()
		}
		superDown := win.Pressed(pixelgl.KeyLeftSuper) || win.Pressed(pixelgl.KeyRightSuper)
		if win.JustPressed(pixelgl.KeyL) && superDown {
			browser.UrlInput.Focus()
		}
		if win.JustPressed(pixelgl.KeyA) && superDown {
			browser.UrlInput.SelectAll()
		}
		if win.JustPressed(pixelgl.KeyTab) || win.JustPressed(pixelgl.KeyEscape) {
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
		browser.Draw()
		win.Update()
		devtoolsWin.Update()

		<-fps
	}
}

func main() {
	pixelgl.Run(run)
}
