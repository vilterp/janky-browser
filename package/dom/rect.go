package dom

import (
	"encoding/xml"
	"image/color"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type RectNode struct {
	XMLName xml.Name `xml:"rect"`

	x      float64
	y      float64
	width  float64
	height float64
	fill   color.Color
}

var _ DOMNode = &RectNode{}

func (rn *RectNode) Name() string { return "rect" }
func (rn *RectNode) Attrs() map[string]string {
	return map[string]string{
		"x":      strconv.FormatFloat(rn.x, 'f', 2, 64),
		"y":      strconv.FormatFloat(rn.y, 'f', 2, 64),
		"width":  strconv.FormatFloat(rn.width, 'f', 2, 64),
		"height": strconv.FormatFloat(rn.height, 'f', 2, 64),
		"fill":   colorToString(rn.fill),
	}
}
func (rn *RectNode) Children() []DOMNode {
	return []DOMNode{}
}
func (rn *RectNode) Draw(imd *imdraw.IMDraw) {
	imd.Color = rn.fill
	imd.Push(pixel.V(rn.x, rn.y))
	imd.Push(pixel.V(rn.x+rn.width, rn.y+rn.height))
	imd.Rectangle(0)
	// TODO: stroke
}
