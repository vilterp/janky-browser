package util

import (
	"image"

	"golang.org/x/exp/shiny/screen"
)

type Window struct {
	Win  screen.Window
	Size image.Point
}
