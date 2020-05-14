package dom

import (
	"encoding/xml"
	"image"
	"math"
	"strconv"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dkit"
	"golang.org/x/image/colornames"
)

type CircleNode struct {
	baseNode

	XMLName xml.Name `xml:"circle"`

	Radius float64 `xml:"radius,attr"`
	X      float64 `xml:"x,attr"`
	Y      float64 `xml:"y,attr"`
	Fill   string  `xml:"fill,attr"`
}

var _ Node = &CircleNode{}

func (cn *CircleNode) Init()            {}
func (cn *CircleNode) Name() string     { return "circle" }
func (cn *CircleNode) Children() []Node { return []Node{} }

func (cn *CircleNode) Attrs() map[string]string {
	return map[string]string{
		"radius": strconv.FormatFloat(cn.Radius, 'f', 2, 64),
		"x":      strconv.FormatFloat(cn.X, 'f', 2, 64),
		"y":      strconv.FormatFloat(cn.Y, 'f', 2, 64),
		"fill":   cn.Fill,
	}
}

func (cn *CircleNode) Draw(gc draw2d.GraphicContext) {
	color, ok := colornames.Map[cn.Fill]
	if !ok {
		gc.SetFillColor(colornames.Black)
	} else {
		gc.SetFillColor(color)
	}
	draw2dkit.Circle(gc, cn.X, cn.Y, cn.Radius)
	gc.Fill()
	// TODO: support stroke as well
}

func (cn *CircleNode) Contains(pt image.Point) bool {
	center := image.Pt(int(cn.X), int(cn.Y))
	diff := pt.Sub(center)
	return pointLength(diff) <= cn.Radius
}

func (cn *CircleNode) GetBounds() image.Rectangle {
	return image.Rect(
		int(cn.X-cn.Radius), int(cn.Y-cn.Radius),
		int(cn.X+cn.Radius), int(cn.Y+cn.Radius),
	)
}

func pointLength(pt image.Point) float64 {
	return math.Sqrt(math.Pow(float64(pt.X), 2) + math.Pow(float64(pt.Y), 2))
}
