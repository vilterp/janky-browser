package dom

import (
	"strconv"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type LineNode struct {
	X1 float64
	Y1 float64
	X2 float64
	Y2 float64

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

func (ln *LineNode) Draw(t pixel.Target) {
	imd := imdraw.New(nil)

	color, ok := colornames.Map[ln.Stroke]
	if ok {
		imd.Color = color
		imd.Push(pixel.V(ln.X1, ln.Y1))
		imd.Push(pixel.V(ln.X2, ln.Y2))
		imd.Line(2) // TODO: strokeWidth
		imd.Draw(t)
	}
}

func (ln *LineNode) Contains(pixel.Vec) bool {
	return false
}
