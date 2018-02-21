package dom

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel/imdraw"
)

type DOMNode interface {
	Name() string
	Attrs() map[string]string
	Children() []DOMNode

	Draw(*imdraw.IMDraw)
}

func colorToString(c color.Color) string {
	r, g, b, a := c.RGBA()
	return fmt.Sprintf("rgba(%d, %d, %d, %d)", r, g, b, a)
}
