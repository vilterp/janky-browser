package dom

import (
	"github.com/faiface/pixel"
)

type Node interface {
	Name() string
	Attrs() map[string]string
	Children() []Node

	Init()
	Draw(target pixel.Target)
	Contains(pixel.Vec) bool
}
