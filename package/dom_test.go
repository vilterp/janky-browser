package jankybrowser

import (
	"testing"

	"golang.org/x/image/colornames"
)

const circleAndRectSource = `<g>
  <circle Fill="rgba(65535, 0, 0, 65535)" Radius="5.00" X="2.00" Y="3.00" />
  <rect Fill="rgba(0, 0, 65535, 65535)" height="10.00" width="5.00" X="2.00" Y="3.00" />
</g>`

var circleAndRect = &GroupNode{
	Circle: []CircleNode{
		{
			Fill:   colornames.Red,
			Radius: 5,
			X:      2,
			Y:      3,
		},
	},
	Rect: []RectNode{
		{
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

func TestDOMParse(t *testing.T) {
	//source := circleAndRectSource
	source := `<circle radius="10" x="11" y="12" />`
	parsed, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(parsed)
	expected := &CircleNode{
		Radius: 10,
		X:      11,
		Y:      12,
	}
	if parsed == nil {
		t.Fatalf("expected:\n%s\ngot nil", Format(expected))
	}
	if parsed != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", Format(expected), Format(parsed))
	}
}
