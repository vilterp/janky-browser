package dom

import (
	"encoding/xml"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type CircleNode struct {
	XMLName xml.Name `xml:"circle"`

	Radius float64 `xml:"radius,attr"`
	X      float64 `xml:"x,attr"`
	Y      float64 `xml:"y,attr"`
	Fill   string  `xml:"fill,attr"`
}

var _ DOMNode = &CircleNode{}

func (cn *CircleNode) Name() string { return "circle" }
func (cn *CircleNode) Attrs() map[string]string {
	return map[string]string{
		"radius": strconv.FormatFloat(cn.Radius, 'f', 2, 64),
		"x":      strconv.FormatFloat(cn.X, 'f', 2, 64),
		"y":      strconv.FormatFloat(cn.Y, 'f', 2, 64),
		"fill":   cn.Fill,
	}
}
func (cn *CircleNode) Children() []DOMNode {
	return []DOMNode{}
}
func (cn *CircleNode) Draw(t pixel.Target) {
	imd := imdraw.New(nil)

	color, ok := colornames.Map[cn.Fill]
	if !ok {
		imd.Color = colornames.Black
	} else {
		imd.Color = color
	}
	imd.Push(pixel.V(cn.X, cn.Y))
	imd.Circle(cn.Radius, 0)
	// TODO: support stroke as well

	imd.Draw(t)
}
