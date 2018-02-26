package util

import "image/color"

func Clamp(min int, max int, val int) int {
	if val > max {
		return max
	}
	if val < min {
		return min
	}
	return val
}

func WithTransparency(c color.Color, transparency float64) color.Color {
	r, g, b, _ := c.RGBA()
	t := transparency * 255
	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255 - uint8(t),
	}
}
