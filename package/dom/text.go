package dom

import (
	"encoding/xml"
	"log"
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
}

var _ DOMNode = &TextNode{}

var Atlas *text.Atlas

func init() {
	fontName := "Roboto-Regular.ttf"
	face, err := LoadTTF(fontName, 15)
	if err != nil {
		panic(err)
	}
	Atlas = text.NewAtlas(face, text.ASCII)
	log.Println("loaded font ")
}

func (tn *TextNode) Name() string        { return "text" }
func (tn *TextNode) Children() []DOMNode { return []DOMNode{} }
func (tn *TextNode) Attrs() map[string]string {
	return map[string]string{
		"value": tn.Value,
		"x":     strconv.FormatFloat(tn.X, 'f', 2, 64),
		"y":     strconv.FormatFloat(tn.Y, 'f', 2, 64),
	}
}
func (tn *TextNode) Draw(t pixel.Target) {
	//txt := text.New(pixel.V(tn.X, tn.Y), atlas)
	txt := text.New(pixel.V(0, 0), Atlas)
	txt.Color = colornames.Black // TODO: set color as attribute
	txt.Clear()
	txt.WriteString(tn.Value)
	txt.Draw(t, pixel.IM.Moved(pixel.V(tn.X, tn.Y)))
}
