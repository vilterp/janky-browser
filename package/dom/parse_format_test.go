package dom

import (
	"testing"
)

const circleAndRectSource = `
<g>
  <rect height="10.00" width="5.00" x="2.00" y="3.00" />
  <circle radius="5.00" x="2.00" y="3.00" />
</g>`

var circleAndRect = &GroupNode{
	CircleNode: []*CircleNode{
		{
			Radius: 5,
			X:      2,
			Y:      3,
		},
	},
	RectNode: []*RectNode{
		{
			X:      2,
			Y:      3,
			Width:  5,
			Height: 10,
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
	source := circleAndRectSource
	parsed, err := Parse([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(parsed)
	expected := circleAndRect
	if parsed == nil {
		t.Fatalf("expected:\n%s\ngot nil", Format(expected))
	}
	if Format(parsed) != Format(expected) {
		t.Fatalf("expected:\n%s\ngot:\n%s", Format(expected), Format(parsed))
	}
}
