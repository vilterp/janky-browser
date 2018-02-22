package dom

import (
	"encoding/xml"

	"github.com/faiface/pixel"
)

type GroupNode struct {
	XMLName xml.Name `xml:"g"`

	Href string `xml:"href,attr"`

	// TODO: this destroys ordering...
	// not sure how to get it to understand an interface...
	RectNode   []*RectNode   `xml:"rect"`
	CircleNode []*CircleNode `xml:"circle"`
	TextNode   []*TextNode   `xml:"text"`
	GroupNode  []*GroupNode  `xml:"g"`
	LineNode   []*LineNode   `xml:"line"`
	//
	//children []Node
}

var _ Node = &GroupNode{}

func (gn *GroupNode) Init() {
	for _, child := range gn.Children() {
		child.Init()
	}
}

func (gn *GroupNode) Name() string             { return "g" }
func (gn *GroupNode) Attrs() map[string]string { return make(map[string]string) }

func (gn *GroupNode) Children() []Node {
	var ret []Node
	for _, rect := range gn.RectNode {
		ret = append(ret, rect)
	}
	for _, circle := range gn.CircleNode {
		ret = append(ret, circle)
	}
	for _, text := range gn.TextNode {
		ret = append(ret, text)
	}
	for _, group := range gn.GroupNode {
		ret = append(ret, group)
	}
	for _, line := range gn.LineNode {
		ret = append(ret, line)
	}
	return ret
}
func (gn *GroupNode) Draw(t pixel.Target) {
	for _, child := range gn.Children() {
		// TODO: draw witn transform
		child.Draw(t)
	}
}

func (gn *GroupNode) Contains(pt pixel.Vec) bool {
	for _, child := range gn.Children() {
		if child.Contains(pt) {
			return true
		}
	}
	return false
}
