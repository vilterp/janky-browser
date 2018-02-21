package jankybrowser

import (
	"fmt"
	"image/color"
	"sort"
	"strconv"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
	"golang.org/x/net/html"
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

// domNodeFromParserNode discards anything it doesn't understand.
// TODO: break this up, move it away...
func domNodeFromParserNode(node *html.Node) DOMNode {
	switch node.Type {
	case html.ElementNode:
		switch node.Data {
		case "g":
			g := &GroupNode{}
			for child := node.FirstChild; child.NextSibling != nil; child = child.NextSibling {
				childDOMNode := domNodeFromParserNode(child)
				if childDOMNode != nil {
					g.children = append(g.children, childDOMNode)
				}
			}
		case "circle":
			circle := &CircleNode{}
			for _, attr := range node.Attr {
				switch attr.Key {
				case "x":
					f, _ := strconv.ParseFloat(attr.Val, 2)
					circle.x = f
				case "y":
					f, _ := strconv.ParseFloat(attr.Val, 2)
					circle.y = f
				case "radius":
					f, _ := strconv.ParseFloat(attr.Val, 2)
					circle.radius = f
				case "color":
					color, ok := colornames.Map[attr.Val]
					if ok {
						circle.fill = color
					}
				}
			}
			return circle
		case "rect":
			rect := &RectNode{}
			for _, attr := range node.Attr {
				switch attr.Key {
				case "x":
					f, _ := strconv.ParseFloat(attr.Val, 2)
					rect.x = f
				case "y":
					f, _ := strconv.ParseFloat(attr.Val, 2)
					rect.y = f
				case "width":
					f, _ := strconv.ParseFloat(attr.Val, 2)
					rect.width = f
				case "height":
					f, _ := strconv.ParseFloat(attr.Val, 2)
					rect.height = f
				case "color":
					color, ok := colornames.Map[attr.Val]
					if ok {
						rect.fill = color
					}
				}
			}
			return rect
		}
	default:
		return nil
	}
	return nil
}

type GroupNode struct {
	children []DOMNode
}

var _ DOMNode = &GroupNode{}

func (gn *GroupNode) Name() string             { return "g" }
func (gn *GroupNode) Attrs() map[string]string { return make(map[string]string) }
func (gn *GroupNode) Children() []DOMNode {
	return gn.children
}
func (gn *GroupNode) Draw(imd *imdraw.IMDraw) {
	for _, child := range gn.children {
		// TODO: draw witn transform
		child.Draw(imd)
	}
}

type CircleNode struct {
	radius float64
	x      float64
	y      float64
	fill   color.Color
}

var _ DOMNode = &CircleNode{}

func (cn *CircleNode) Name() string { return "circle" }
func (cn *CircleNode) Attrs() map[string]string {
	return map[string]string{
		"radius": strconv.FormatFloat(cn.radius, 'f', 2, 64),
		"x":      strconv.FormatFloat(cn.x, 'f', 2, 64),
		"y":      strconv.FormatFloat(cn.y, 'f', 2, 64),
		"fill":   colorToString(cn.fill),
	}
}
func (cn *CircleNode) Children() []DOMNode {
	return []DOMNode{}
}
func (cn *CircleNode) Draw(imd *imdraw.IMDraw) {
	imd.Color = cn.fill
	imd.Push(pixel.V(cn.x, cn.y))
	imd.Circle(cn.radius, 0)
	// TODO: support stroke as well
}

type RectNode struct {
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
