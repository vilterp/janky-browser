package dom

import (
	"encoding/xml"
	"image"
	"strconv"

	"github.com/llgcode/draw2d"
	"golang.org/x/image/colornames"
)

type TextNode struct {
	baseNode

	XMLName xml.Name `xml:"text"`

	Value string  `xml:"value,attr"`
	X     float64 `xml:"x,attr"`
	Y     float64 `xml:"y,attr"`
	Fill  string  `xml:"fill,attr"`

	bounds image.Rectangle
}

var _ Node = &TextNode{}

func (tn *TextNode) Name() string     { return "text" }
func (tn *TextNode) Children() []Node { return []Node{} }

func (tn *TextNode) Attrs() map[string]string {
	return map[string]string{
		"value": tn.Value,
		"x":     strconv.FormatFloat(tn.X, 'f', 2, 64),
		"y":     strconv.FormatFloat(tn.Y, 'f', 2, 64),
		"fill":  tn.Fill,
	}
}

func (tn *TextNode) Init() {}

func (tn *TextNode) Draw(gc draw2d.GraphicContext) {
	gc.SetFontData(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: 0})
	// TODO: configurable font size and color
	gc.SetFontSize(TextHeight)

	color, ok := colornames.Map[tn.Fill]
	if ok {
		gc.SetFillColor(color)
	} else {
		gc.SetFillColor(image.Black)
	}

	gc.FillStringAt(tn.Value, tn.X, tn.Y)
	// TODO: uh... what do these numbers mean?
	left, top, right, bottom := gc.GetStringBounds(tn.Value)
	tn.bounds = image.Rect(int(left), int(top), int(right), int(bottom))
}

func (tn *TextNode) Contains(pt image.Point) bool {
	return pt.In(tn.bounds)
}

// TODO: support multiple font sizes
const TextHeight = 13

func (tn *TextNode) GetBounds() image.Rectangle {
	return tn.bounds
}
