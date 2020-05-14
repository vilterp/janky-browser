package dom

import (
	"encoding/xml"
	"image"
	"strconv"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/vilterp/janky-browser/package/util"
	"golang.org/x/image/colornames"
)

type RectNode struct {
	baseNode

	XMLName xml.Name `xml:"rect"`

	X            float64 `xml:"x,attr"`
	Y            float64 `xml:"y,attr"`
	Width        float64 `xml:"width,attr"`
	Height       float64 `xml:"height,attr"`
	Fill         string  `xml:"fill,attr"`
	Transparency float64 `xml:"transparency,attr"` // [0, 1]. TODO: this should really be in the fill itself.
	Stroke       string  `xml:"stroke,attr"`
}

var _ Node = &RectNode{}

func (rn *RectNode) Init()            {}
func (rn *RectNode) Name() string     { return "rect" }
func (rn *RectNode) Children() []Node { return []Node{} }

func (rn *RectNode) Attrs() map[string]string {
	return map[string]string{
		"x":      strconv.FormatFloat(rn.X, 'f', 2, 64),
		"y":      strconv.FormatFloat(rn.Y, 'f', 2, 64),
		"width":  strconv.FormatFloat(rn.Width, 'f', 2, 64),
		"height": strconv.FormatFloat(rn.Height, 'f', 2, 64),
		"fill":   rn.Fill,
		"stroke": rn.Stroke,
	}
}

func (rn *RectNode) Draw(gc draw2d.GraphicContext) {
	// Draw fill.
	fillColor, ok := colornames.Map[rn.Fill]
	if ok {
		gc.SetFillColor(util.WithTransparency(fillColor, rn.Transparency))
		draw2dkit.Rectangle(gc, rn.X, rn.Y, rn.X+rn.Width, rn.Y+rn.Height)
		gc.Fill()
	}

	// Draw stroke.
	strokeColor, ok := colornames.Map[rn.Stroke]
	if ok {
		gc.SetStrokeColor(strokeColor)
		draw2dkit.Rectangle(gc, rn.X, rn.Y, rn.X+rn.Width, rn.Y+rn.Height)
		gc.Stroke()
	}
}

func (rn *RectNode) Contains(pt image.Point) bool {
	return pt.In(rn.GetBounds())
}

func (rn *RectNode) GetBounds() image.Rectangle {
	return image.Rect(int(rn.X), int(rn.Y), int(rn.X+rn.Width), int(rn.Y+rn.Height))
}

func RectFromBounds(bounds image.Rectangle) *RectNode {
	return &RectNode{
		X:      float64(bounds.Min.X),
		Y:      float64(bounds.Min.Y),
		Width:  float64(bounds.Dx()),
		Height: float64(bounds.Dy()),
	}
}
