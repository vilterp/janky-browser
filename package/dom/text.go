package dom

import (
	"encoding/xml"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type TextNode struct {
	XMLName xml.Name `xml:"text"`

	Value string  `xml:"value,attr"`
	X     float64 `xml:"x,attr"`
	Y     float64 `xml:"y,attr"`
	Fill  string  `xml:"fill,attr"`

	txt *text.Text
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

func (tn *TextNode) Init() {
	tn.txt = text.New(pixel.V(0, 0), Atlas)
}

func (tn *TextNode) Draw(t pixel.Target) {
	color, ok := colornames.Map[tn.Fill]
	if ok {
		tn.txt.Color = color
	} else {
		tn.txt.Color = colornames.Black
	}
	tn.txt.Clear()
	tn.txt.WriteString(tn.Value)
	tn.txt.Draw(t, pixel.IM.Moved(pixel.V(tn.X, tn.Y)))
}

func (tn *TextNode) Contains(pt pixel.Vec) bool {
	txtBounds := tn.txt.Bounds()
	movedBounds := txtBounds.Moved(pixel.V(tn.X, tn.Y))
	return movedBounds.Contains(pt)
}

var Atlas *text.Atlas

func init() {
	//fontName := "Roboto-Regular.ttf"
	//face, err := LoadTTF(fontName, 15)
	//if err != nil {
	//	panic(err)
	//}
	//Atlas = text.NewAtlas(face, text.ASCII)
	//log.Println("loaded font ")
	Atlas = text.Atlas7x13
}
