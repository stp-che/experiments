package ui

import "github.com/faiface/pixel"

func topLeft(r pixel.Rect) pixel.Vec {
	return pixel.V(r.Min.X, r.Max.Y)
}
