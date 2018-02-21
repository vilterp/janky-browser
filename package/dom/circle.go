package dom

import (
	"encoding/xml"
	"image/color"
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
	Fill   color.Color
}

var _ DOMNode = &CircleNode{}

func (cn *CircleNode) Name() string { return "circle" }
func (cn *CircleNode) Attrs() map[string]string {
	m := map[string]string{
		"radius": strconv.FormatFloat(cn.Radius, 'f', 2, 64),
		"x":      strconv.FormatFloat(cn.X, 'f', 2, 64),
		"y":      strconv.FormatFloat(cn.Y, 'f', 2, 64),
	}
	if cn.Fill != nil {
		m["fill"] = colorToString(cn.Fill)
	}
	return m
}
func (cn *CircleNode) Children() []DOMNode {
	return []DOMNode{}
}
func (cn *CircleNode) Draw(imd *imdraw.IMDraw) {
	if cn.Fill != nil {
		imd.Color = cn.Fill
	} else {
		imd.Color = colornames.Black
	}
	imd.Push(pixel.V(cn.X, cn.Y))
	imd.Circle(cn.Radius, 0)
	// TODO: support stroke as well
}
