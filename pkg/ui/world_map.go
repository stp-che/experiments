package ui

import (
	"experiments/pkg/sim"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type WorldMap struct {
	World           *sim.World
	regionSize      int
	regionSizeFloat float64
	bounds          pixel.Rect
}

func (b *WorldMap) fitInto(bounds pixel.Rect) {
	b.regionSize = int(bounds.W()) / b.World.Cols
	if b.regionSize*b.World.Rows > int(bounds.H()) {
		b.regionSize = int(bounds.H()) / b.World.Rows
	}
	b.bounds = bounds.Resized(topLeft(bounds), pixel.V(float64(b.regionSize*b.World.Cols), float64(b.regionSize*b.World.Rows)))
	b.regionSizeFloat = float64(b.regionSize)
}

func (b *WorldMap) CellSize() float64 {
	return b.regionSizeFloat
}

func (b *WorldMap) Bounds() pixel.Rect {
	return b.bounds
}

func (b *WorldMap) Render(imd *imdraw.IMDraw) {
	b.renderTable(imd)
	for i, r := range b.World.Regions {
		b.renderRegion(i, r, imd)
	}
}

func (b *WorldMap) renderTable(imd *imdraw.IMDraw) {
	imd.Color = worlMapRegionBorderColor
	imd.Push(b.bounds.Min)
	imd.Push(b.bounds.Max)
	imd.Rectangle(1)
	for i := 1; i < b.World.Cols; i++ {
		x := b.bounds.Min.X + float64(b.regionSize*i)
		imd.Push(pixel.V(x, b.bounds.Min.Y))
		imd.Push(pixel.V(x, b.bounds.Max.Y))
		imd.Line(1)
	}
	for i := 1; i < b.World.Rows; i++ {
		y := b.bounds.Max.Y - float64(b.regionSize*i)
		imd.Push(pixel.V(b.bounds.Min.X, y))
		imd.Push(pixel.V(b.bounds.Max.X, y))
		imd.Line(1)
	}
}

func (b *WorldMap) renderRegion(pos int, reg *sim.Region, imd *imdraw.IMDraw) {
	if reg.Content == sim.RCNone {
		return
	}
	switch reg.Content {
	case sim.RCWall:
		imd.Color = wallColor
	case sim.RCBot:
		imd.Color = botColor
	case sim.RCFood:
		imd.Color = foodColor
	}
	x, y := b.World.XYPos(pos)
	topLeft := pixel.V(
		b.bounds.Min.X+float64(b.regionSize*x),
		b.bounds.Max.Y-float64(b.regionSize*y)-1,
	)
	imd.Push(topLeft)
	imd.Push(topLeft.Add(pixel.V(b.regionSizeFloat-1, -b.regionSizeFloat+1)))
	imd.Rectangle(0)
}

func newWorldMap(world *sim.World, bounds pixel.Rect) *WorldMap {
	m := &WorldMap{
		World: world,
	}
	m.fitInto(bounds)
	return m
}
