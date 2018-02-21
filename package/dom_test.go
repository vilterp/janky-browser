package jankybrowser

import (
	"strings"
	"testing"

	"golang.org/x/image/colornames"
	"golang.org/x/net/html"
)

const circleAndRectSource = `<g>
  <circle fill="rgba(65535, 0, 0, 65535)" radius="5.00" x="2.00" y="3.00" />
  <rect fill="rgba(0, 0, 65535, 65535)" height="10.00" width="5.00" x="2.00" y="3.00" />
</g>`

var circleAndRect = &GroupNode{
	children: []DOMNode{
		&CircleNode{
			fill:   colornames.Red,
			radius: 5,
			x:      2,
			y:      3,
		},
		&RectNode{
			fill:   colornames.Blue,
			x:      2,
			y:      3,
			width:  5,
			height: 10,
		},
	},
}

func TestDOMFormat(t *testing.T) {
	node := circleAndRect
	formatted := Format(node)
	expected := circleAndRectSource
	if formatted != expected {
		t.Fatalf("expected:\n%s\n got:\n%s", expected, formatted)
	}
}

func TestDOMFromParserNode(t *testing.T) {
	source := circleAndRectSource
	parsed, err := html.Parse(strings.NewReader(source))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(parsed)
	domNode := domNodeFromParserNode(parsed)
	expected := circleAndRect
	if domNode == nil {
		t.Fatalf("expected:\n%s\ngot nil", Format(expected))
	}
	if domNode != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", Format(expected), Format(domNode))
	}
}
