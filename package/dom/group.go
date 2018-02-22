package dom

import (
	"encoding/xml"

	"github.com/faiface/pixel"
)

type GroupNode struct {
	XMLName xml.Name `xml:"g"`

	RectNode   []RectNode   `xml:"rect"`
	CircleNode []CircleNode `xml:"circle"`
	//
	//children []DOMNode
}

var _ DOMNode = &GroupNode{}

func (gn *GroupNode) Name() string             { return "g" }
func (gn *GroupNode) Attrs() map[string]string { return make(map[string]string) }
func (gn *GroupNode) Children() []DOMNode {
	var ret []DOMNode
	for _, rect := range gn.RectNode {
		ret = append(ret, &rect)
	}
	for _, circle := range gn.CircleNode {
		ret = append(ret, &circle)
	}
	return ret
}
func (gn *GroupNode) Draw(t pixel.Target) {
	for _, child := range gn.Children() {
		// TODO: draw witn transform
		child.Draw(t)
	}
}
