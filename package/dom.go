package jankybrowser

import (
	"encoding/xml"
	"fmt"
	"image/color"
	"sort"
	"strconv"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type DOMNode interface {
	Name() string
	Attrs() map[string]string
	Children() []DOMNode

	Draw(*imdraw.IMDraw)
}

func Format(node DOMNode) string {
	return doFormat(node, 1)
}

func doFormat(node DOMNode, indent int) string {
	var attrs []string
	for key, val := range node.Attrs() {
		attrs = append(attrs, fmt.Sprintf("%s=%#v", key, val))
	}
	sort.Strings(attrs)
	attrsStr := strings.Join(attrs, " ")
	if len(attrs) > 0 {
		attrsStr = " " + attrsStr
	}

	children := node.Children()
	if len(children) > 0 {
		indentStr := strings.Repeat("  ", indent)
		var childrenLines []string
		for _, child := range children {
			childrenLines = append(childrenLines, indentStr+doFormat(child, indent+1))
		}
		childrenStr := strings.Join(childrenLines, "\n")
		return fmt.Sprintf("<%s%s>\n%s\n</%s>", node.Name(), attrsStr, childrenStr, node.Name())
	}
	return fmt.Sprintf("<%s%s />", node.Name(), attrsStr)
}

func Parse(data []byte) (DOMNode, error) {
	g := CircleNode{}
	err := xml.Unmarshal(data, &g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// domNodeFromParserNode discards anything it doesn't understand.
// TODO: break this up, move it away...
//func domNodeFromParserNode(node *html.Node) DOMNode {
//	switch node.Type {
//	case html.ElementNode:
//		switch node.Data {
//		case "g":
//			g := &GroupNode{}
//			for child := node.FirstChild; child.NextSibling != nil; child = child.NextSibling {
//				childDOMNode := domNodeFromParserNode(child)
//				if childDOMNode != nil {
//					g.children = append(g.children, childDOMNode)
//				}
//			}
//		case "circle":
//			circle := &CircleNode{}
//			for _, attr := range node.Attr {
//				switch attr.Key {
//				case "X":
//					f, _ := strconv.ParseFloat(attr.Val, 2)
//					circle.X = f
//				case "Y":
//					f, _ := strconv.ParseFloat(attr.Val, 2)
//					circle.Y = f
//				case "Radius":
//					f, _ := strconv.ParseFloat(attr.Val, 2)
//					circle.Radius = f
//				case "color":
//					color, ok := colornames.Map[attr.Val]
//					if ok {
//						circle.Fill = color
//					}
//				}
//			}
//			return circle
//		case "rect":
//			rect := &RectNode{}
//			for _, attr := range node.Attr {
//				switch attr.Key {
//				case "X":
//					f, _ := strconv.ParseFloat(attr.Val, 2)
//					rect.X = f
//				case "Y":
//					f, _ := strconv.ParseFloat(attr.Val, 2)
//					rect.Y = f
//				case "width":
//					f, _ := strconv.ParseFloat(attr.Val, 2)
//					rect.width = f
//				case "height":
//					f, _ := strconv.ParseFloat(attr.Val, 2)
//					rect.height = f
//				case "color":
//					color, ok := colornames.Map[attr.Val]
//					if ok {
//						rect.Fill = color
//					}
//				}
//			}
//			return rect
//		}
//	default:
//		return nil
//	}
//	return nil
//}

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

type RectNode struct {
	XMLName xml.Name `xml:"rect"`

	x      float64
	y      float64
	width  float64
	height float64
	fill   color.Color
}

var _ DOMNode = &RectNode{}

func (rn *RectNode) Name() string { return "rect" }
func (rn *RectNode) Attrs() map[string]string {
	return map[string]string{
		"x":      strconv.FormatFloat(rn.x, 'f', 2, 64),
		"y":      strconv.FormatFloat(rn.y, 'f', 2, 64),
		"width":  strconv.FormatFloat(rn.width, 'f', 2, 64),
		"height": strconv.FormatFloat(rn.height, 'f', 2, 64),
		"fill":   colorToString(rn.fill),
	}
}
func (rn *RectNode) Children() []DOMNode {
	return []DOMNode{}
}
func (rn *RectNode) Draw(imd *imdraw.IMDraw) {
	imd.Color = rn.fill
	imd.Push(pixel.V(rn.x, rn.y))
	imd.Push(pixel.V(rn.x+rn.width, rn.y+rn.height))
	imd.Rectangle(0)
	// TODO: stroke
}

func colorToString(c color.Color) string {
	r, g, b, a := c.RGBA()
	return fmt.Sprintf("rgba(%d, %d, %d, %d)", r, g, b, a)
}
