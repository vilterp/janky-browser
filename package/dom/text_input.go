package dom

import (
	"fmt"
	"strconv"

	"github.com/faiface/pixel"
)

// TODO: additional keyboard movement
// - word boundaries on option+{left, right}

type TextInputNode struct {
	// props
	X         float64
	Y         float64
	Value     string
	Width     float64
	TextColor string
	Focused   bool
	OnEnter   func(string)

	// state
	cursorPos      int
	selectionStart *int

	// children
	group          *GroupNode
	valueText      *TextNode
	cursorLine     *LineNode
	backgroundRect *RectNode
	selectionRect  *RectNode
}

var _ Node = &TextInputNode{}

func (tin *TextInputNode) Init() {
	tin.backgroundRect = &RectNode{}
	tin.selectionRect = &RectNode{}
	tin.valueText = &TextNode{}
	tin.cursorLine = &LineNode{}
	tin.group = &GroupNode{
		RectNode: []*RectNode{
			tin.backgroundRect,
			tin.selectionRect,
		},
		LineNode: []*LineNode{
			tin.cursorLine,
		},
		TextNode: []*TextNode{
			tin.valueText,
		},
	}
	tin.group.Init()
}

func (tin *TextInputNode) Name() string { return "textInput" }
func (tin *TextInputNode) Children() []Node {
	// not reporting any children... the "shadow dom" I guess...
	return []Node{}
}

func (tin *TextInputNode) Attrs() map[string]string {
	return map[string]string{
		"x":         strconv.FormatFloat(tin.X, 'f', 2, 64),
		"y":         strconv.FormatFloat(tin.Y, 'f', 2, 64),
		"value":     tin.Value,
		"width":     strconv.FormatFloat(tin.Width, 'f', 2, 64),
		"textColor": tin.TextColor,
		"focused":   fmt.Sprintf("%v", tin.Focused),
	}
}

func (tin *TextInputNode) Draw(t pixel.Target) {
	// Update background rect.
	tin.backgroundRect.Width = tin.Width
	tin.backgroundRect.X = tin.X
	tin.backgroundRect.Y = tin.Y
	tin.backgroundRect.Height = 30
	if tin.Focused {
		tin.backgroundRect.Stroke = "black"
	} else {
		tin.backgroundRect.Stroke = ""
	}

	textStartX := tin.X + 5

	// Update text.
	tin.valueText.Fill = tin.TextColor
	tin.valueText.Value = tin.Value
	tin.valueText.X = textStartX
	tin.valueText.Y = tin.Y + 10

	// Update cursor.
	const charWidth = float64(7)
	cursorX := textStartX + float64(tin.cursorPos)*charWidth
	tin.cursorLine.X1 = cursorX
	tin.cursorLine.X2 = cursorX
	tin.cursorLine.Y1 = tin.Y + 21
	tin.cursorLine.Y2 = tin.Y + 7
	if tin.Focused {
		tin.cursorLine.Stroke = "black"
	} else {
		tin.cursorLine.Stroke = ""
	}

	// Update selection.
	if tin.selectionStart == nil {
		tin.selectionRect.Fill = ""
	} else {
		tin.selectionRect.Fill = "pink"
		startIdx, endIdx := tin.GetSelection()
		tin.selectionRect.X = textStartX + float64(startIdx)*charWidth
		tin.selectionRect.Y = tin.Y + 7
		tin.selectionRect.Width = float64(endIdx-startIdx) * charWidth
		tin.selectionRect.Height = 13
	}

	tin.group.Draw(t)
}

func (tin *TextInputNode) Contains(pt pixel.Vec) bool {
	return tin.backgroundRect.Contains(pt)
}

// Event handling stuff.

func (tin *TextInputNode) ProcessTyping(t string) {
	if !tin.Focused {
		return
	}
	if tin.selectionStart != nil {
		tin.DeleteSelection()
	}
	tin.Value = tin.Value[:tin.cursorPos] + t + tin.Value[tin.cursorPos:]
	tin.cursorPos += 1
}

func (tin *TextInputNode) ProcessBackspace() {
	if !tin.Focused {
		return
	}
	if tin.Value == "" {
		return
	}
	if tin.selectionStart == nil {
		tin.Value = tin.Value[:tin.cursorPos-1] + tin.Value[tin.cursorPos:]
		tin.cursorPos -= 1
	} else {
		tin.DeleteSelection()
	}
}

func (tin *TextInputNode) DeleteSelection() {
	startIdx, endIdx := tin.GetSelection()
	tin.Value = tin.Value[:startIdx] + tin.Value[endIdx:]
	tin.CancelSelection()
	tin.cursorPos = startIdx
}

func (tin *TextInputNode) GetSelection() (int, int) {
	if tin.selectionStart == nil {
		return tin.cursorPos, tin.cursorPos
	}
	selectionStart := *tin.selectionStart
	var startIdx int
	var endIdx int
	if selectionStart < tin.cursorPos {
		startIdx = selectionStart
		endIdx = tin.cursorPos
	} else {
		endIdx = selectionStart
		startIdx = tin.cursorPos
	}
	return startIdx, endIdx
}

func (tin *TextInputNode) ProcessEnter() {
	if !tin.Focused {
		return
	}
	tin.OnEnter(tin.Value)
}

func (tin *TextInputNode) Focus() {
	tin.Focused = true
	if len(tin.Value) > 0 {
		tin.cursorPos = len(tin.Value)
		selectionStart := 0
		tin.selectionStart = &selectionStart
		return
	}
	tin.cursorPos = 0
	tin.selectionStart = nil
}

func (tin *TextInputNode) UnFocus() {
	tin.Focused = false
	tin.CancelSelection()
}

func (tin *TextInputNode) ProcessLeftKey(shiftDown bool, superDown bool) {
	if !tin.Focused {
		return
	}
	tin.MaybeStartSelection(shiftDown)
	if superDown {
		tin.cursorPos = 0
		return
	}
	tin.cursorPos = tin.cursorPos - 1
	if tin.cursorPos < 0 {
		tin.cursorPos = 0
	}
}

func (tin *TextInputNode) ProcessRightKey(shiftDown bool, superDown bool) {
	if !tin.Focused {
		return
	}
	tin.MaybeStartSelection(shiftDown)
	if superDown {
		tin.cursorPos = len(tin.Value)
		return
	}
	tin.cursorPos = tin.cursorPos + 1
	if tin.cursorPos > len(tin.Value) {
		tin.cursorPos = len(tin.Value)
	}
}

func (tin *TextInputNode) MaybeStartSelection(shiftDown bool) {
	if !shiftDown {
		tin.CancelSelection()
		return
	}
	if tin.selectionStart != nil {
		return
	}
	selectionStart := tin.cursorPos
	tin.selectionStart = &selectionStart
}

func (tin *TextInputNode) CancelSelection() {
	tin.selectionStart = nil
}
