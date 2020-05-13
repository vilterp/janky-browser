package main

import (
	"fmt"
	"image"
	"log"

	"github.com/llgcode/draw2d/draw2dimg"
	jankybrowser "github.com/vilterp/janky-browser/package"
	"github.com/vilterp/janky-browser/package/util"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

const initPage = "http://localhost:8084/circleRectText.svg"

func main() {
	driver.Main(func(curScreen screen.Screen) {
		initialSize := image.Pt(800, 800)

		window, err := curScreen.NewWindow(&screen.NewWindowOptions{
			Title:  "Janky Browser",
			Width:  initialSize.X,
			Height: initialSize.Y,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer window.Release()

		winWrap := &util.Window{
			Win:  window,
			Size: initialSize,
		}

		browser := jankybrowser.NewBrowser(winWrap, initPage)

		for {
			evt := window.NextEvent()

			fmt.Printf("event: %#v\n", evt)

			switch tEvt := evt.(type) {
			case paint.Event:
				buf, err := curScreen.NewBuffer(winWrap.Size)
				if err != nil {
					log.Fatal(err)
				}
				img := buf.RGBA()
				gc := draw2dimg.NewGraphicContext(img)

				browser.Draw(gc)

				window.Upload(image.Point{}, buf, buf.Bounds())
				window.Publish()
				buf.Release()
			case mouse.Event:
				// Handle mouse events.
				browser.ProcessMouseEvents(
					image.Pt(int(tEvt.X), int(tEvt.Y)),
					tEvt.Direction == mouse.DirPress,
					true, // ???
				)
			case key.Event:
				fmt.Println("rune:", string(tEvt.Rune))
				typed := tEvt.Rune
				if len(string(typed)) > 0 { // ???
					browser.UrlInput.ProcessTyping(string(typed))
				}
				if tEvt.Code == key.CodeDeleteBackspace {
					browser.UrlInput.ProcessBackspace()
				}
				if tEvt.Code == key.CodeReturnEnter {
					browser.UrlInput.ProcessEnter()
				}
				superDown := tEvt.Modifiers&key.ModMeta > 0 // ??
				if tEvt.Code == key.CodeL && superDown {
					browser.UrlInput.Focus()
				}
				if tEvt.Code == key.CodeA && superDown {
					browser.UrlInput.SelectAll()
				}
				if tEvt.Code == key.CodeTab || tEvt.Code == key.CodeEscape {
					browser.UrlInput.UnFocus()
				}
				shiftDown := tEvt.Modifiers&key.ModShift > 0
				if tEvt.Code == key.CodeLeftArrow {
					browser.UrlInput.ProcessLeftKey(shiftDown, superDown)
				}
				if tEvt.Code == key.CodeRightArrow {
					browser.UrlInput.ProcessRightKey(shiftDown, superDown)
				}
			case size.Event:
				winWrap.Size = tEvt.Size()
			case lifecycle.Event:
				if tEvt.To == lifecycle.StageDead {
					return
				}
			}
		}
	})
}
