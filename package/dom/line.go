package dom

import (
	"image"
	"strconv"

	"github.com/llgcode/draw2d"
	"golang.org/x/image/colornames"
)

type LineNode struct {
	baseNode

	X1     float64
	Y1     float64
	X2     float64
	Y2     float64
	Stroke string
}

var _ Node = &LineNode{}

func (ln *LineNode) Name() string     { return "line" }
func (ln *LineNode) Children() []Node { return []Node{} }
func (ln *LineNode) Init()            {}

func (ln *LineNode) Attrs() map[string]string {
	return map[string]string{
		"x1": strconv.FormatFloat(ln.X1, 'f', 2, 64),
		"y1": strconv.FormatFloat(ln.Y1, 'f', 2, 64),
		"x2": strconv.FormatFloat(ln.X2, 'f', 2, 64),
		"y2": strconv.FormatFloat(ln.Y2, 'f', 2, 64),
	}
}

func (ln *LineNode) Draw(gc draw2d.GraphicContext) {
	color, ok := colornames.Map[ln.Stroke]
	if ok {
		gc.SetStrokeColor(color)
		gc.MoveTo(ln.X1, ln.Y1)
		gc.LineTo(ln.X2, ln.Y2)
		gc.SetLineWidth(2)
		gc.Stroke() // stroke width
	}
}

func (ln *LineNode) Contains(image.Point) bool {
	return false
}

func (ln *LineNode) GetBounds() image.Rectangle {
	return image.Rect(int(ln.X1), int(ln.Y1), int(ln.X2), int(ln.X2))
}
