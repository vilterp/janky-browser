package jankybrowser

import "golang.org/x/image/colornames"

func ExampleDOMTree() DOMNode {
	return &GroupNode{
		children: []DOMNode{
			&CircleNode{
				fill:   colornames.Red,
				radius: 50,
				x:      200,
				y:      300,
			},
			&RectNode{
				fill:   colornames.Blue,
				x:      20,
				y:      30,
				width:  500,
				height: 100,
			},
		},
	}
}
