package dom

import (
	"encoding/xml"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
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

func (rn *RectNode) Draw(t pixel.Target) {
	imd := imdraw.New(nil)

	// Draw fill.
	fillColor, ok := colornames.Map[rn.Fill]
	if ok {
		imd.Color = util.WithTransparency(fillColor, rn.Transparency)
		imd.Push(pixel.V(rn.X, rn.Y))
		imd.Push(pixel.V(rn.X+rn.Width, rn.Y+rn.Height))
		imd.Rectangle(0)
	}

	// Draw stroke.
	strokeColor, ok := colornames.Map[rn.Stroke]
	if ok {
		imd.Color = strokeColor
		imd.Push(pixel.V(rn.X, rn.Y))
		imd.Push(pixel.V(rn.X+rn.Width, rn.Y+rn.Height))
		imd.Rectangle(2)
	}

	imd.Draw(t)
}

func (rn *RectNode) Contains(pt pixel.Vec) bool {
	return rn.GetBounds().Contains(pt)
}

func (rn *RectNode) GetBounds() pixel.Rect {
	return pixel.R(rn.X, rn.Y, rn.X+rn.Width, rn.Y+rn.Height)
}

func RectFromBounds(bounds pixel.Rect) *RectNode {
	return &RectNode{
		X:      bounds.Min.X,
		Y:      bounds.Min.Y,
		Width:  bounds.W(),
		Height: bounds.H(),
	}
}
