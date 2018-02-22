package dom

import (
	"encoding/xml"
	"image/color"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type RectNode struct {
	XMLName xml.Name `xml:"rect"`

	X      float64 `xml:"x,attr"`
	Y      float64 `xml:"y,attr"`
	Width  float64 `xml:"width,attr"`
	Height float64 `xml:"height,attr"`
	fill   color.Color
}

var _ DOMNode = &RectNode{}

func (rn *RectNode) Name() string { return "rect" }
func (rn *RectNode) Attrs() map[string]string {
	m := map[string]string{
		"x":      strconv.FormatFloat(rn.X, 'f', 2, 64),
		"y":      strconv.FormatFloat(rn.Y, 'f', 2, 64),
		"width":  strconv.FormatFloat(rn.Width, 'f', 2, 64),
		"height": strconv.FormatFloat(rn.Height, 'f', 2, 64),
	}
	if rn.fill != nil {
		m["fill"] = colorToString(rn.fill)
	}
	return m
}
func (rn *RectNode) Children() []DOMNode {
	return []DOMNode{}
}
func (rn *RectNode) Draw(t pixel.Target) {
	imd := imdraw.New(nil)

	if rn.fill != nil {
		imd.Color = rn.fill
	} else {
		imd.Color = colornames.Black
	}
	imd.Push(pixel.V(rn.X, rn.Y))
	imd.Push(pixel.V(rn.X+rn.Width, rn.Y+rn.Height))
	imd.Rectangle(0)
	// TODO: stroke

	imd.Draw(t)
}
