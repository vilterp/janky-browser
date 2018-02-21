package dom

import (
	"encoding/xml"

	"github.com/faiface/pixel/imdraw"
)

type GroupNode struct {
	XMLName xml.Name `xml:"g"`

	Rect   []RectNode
	Circle []CircleNode
	//
	//children []DOMNode
}

var _ DOMNode = &GroupNode{}

func (gn *GroupNode) Name() string             { return "g" }
func (gn *GroupNode) Attrs() map[string]string { return make(map[string]string) }
func (gn *GroupNode) Children() []DOMNode {
	var ret []DOMNode
	for _, rect := range gn.Rect {
		ret = append(ret, &rect)
	}
	for _, circle := range gn.Circle {
		ret = append(ret, &circle)
	}
	return ret
}
func (gn *GroupNode) Draw(imd *imdraw.IMDraw) {
	for _, child := range gn.Children() {
		// TODO: draw witn transform
		child.Draw(imd)
	}
}
