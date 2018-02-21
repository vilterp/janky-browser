package jankybrowser

import (
	"testing"

	"golang.org/x/image/colornames"
)

func TestDOMFormatting(t *testing.T) {
	node := &GroupNode{
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
	formatted := Format(node)
	expected := `<g>
  <circle fill="rgba(65535, 0, 0, 65535)" radius="5.00" x="2.00" y="3.00" />
  <rect fill="rgba(0, 0, 65535, 65535)" height="10.00" width="5.00" x="2.00" y="3.00" />
</g>`
	if formatted != expected {
		t.Fatalf("expected:\n%s\n got:\n%s", expected, formatted)
	}
}
