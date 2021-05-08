package ui

import (
	"experiments/pkg/sim"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type filledCell struct {
	Cell  sim.Pos
	Color color.RGBA
}

type WorldMap struct {
	Board         *sim.World
	cellSize      int
	cellSizeFloat float64
	bounds        pixel.Rect
	filledCells   []filledCell
}

func (b *WorldMap) FitInto(bounds pixel.Rect) {
	b.cellSize = int(bounds.W()) / b.Board.Cols
	if b.cellSize*b.Board.Rows > int(bounds.H()) {
		b.cellSize = int(bounds.H()) / b.Board.Rows
	}
	b.bounds = bounds.Resized(topLeft(bounds), pixel.V(float64(b.cellSize*b.Board.Cols), float64(b.cellSize*b.Board.Rows)))
	b.cellSizeFloat = float64(b.cellSize)
	// fmt.Printf("%v\n", bounds)
	// fmt.Printf("%v\n", b.bounds)
}

func (b *WorldMap) CellSize() float64 {
	return b.cellSizeFloat
}

func (b *WorldMap) Bounds() pixel.Rect {
	return b.bounds
}

func (b *WorldMap) FillCell(cell sim.Pos, clr color.RGBA) {
	if b.filledCells == nil {
		b.filledCells = make([]filledCell, 10)
	}
	b.filledCells = append(b.filledCells, filledCell{
		Cell:  cell,
		Color: clr,
	})
}

func (b *WorldMap) Clear() {
	b.filledCells = b.filledCells[:0]
}

func (b *WorldMap) Render(imd *imdraw.IMDraw) {
	b.renderTable(imd)
	for _, cell := range b.filledCells {
		b.renderFilledCell(cell, imd)
	}
}

func (b *WorldMap) renderTable(imd *imdraw.IMDraw) {
	imd.Color = worlMapRegionBorderColor
	imd.Push(b.bounds.Min)
	imd.Push(b.bounds.Max)
	imd.Rectangle(1)
	for i := 1; i < b.Board.Cols; i++ {
		x := b.bounds.Min.X + float64(b.cellSize*i)
		imd.Push(pixel.V(x, b.bounds.Min.Y))
		imd.Push(pixel.V(x, b.bounds.Max.Y))
		imd.Line(1)
	}
	for i := 1; i < b.Board.Rows; i++ {
		y := b.bounds.Max.Y - float64(b.cellSize*i)
		imd.Push(pixel.V(b.bounds.Min.X, y))
		imd.Push(pixel.V(b.bounds.Max.X, y))
		imd.Line(1)
	}
}

func (b *WorldMap) renderFilledCell(cell filledCell, imd *imdraw.IMDraw) {
	imd.Color = cell.Color
	topLeft := pixel.V(
		b.bounds.Min.X+float64(b.cellSize*cell.Cell[0]),
		b.bounds.Max.Y-float64(b.cellSize*cell.Cell[1])-1,
	)
	imd.Push(topLeft)
	imd.Push(topLeft.Add(pixel.V(b.cellSizeFloat-1, -b.cellSizeFloat+1)))
	imd.Rectangle(0)
}

func newWorldMap(board *sim.World, bounds pixel.Rect) *WorldMap {
	b := &WorldMap{
		Board: board,
	}
	b.FitInto(bounds)
	return b
}
