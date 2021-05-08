package ui

import (
	"experiments/pkg/sim"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Scene struct {
	WorldMap *WorldMap
	Sim      *sim.Simulation
	imd      *imdraw.IMDraw
}

func (s *Scene) Draw(win *pixelgl.Window) {
	if s.imd == nil {
		s.imd = imdraw.New(nil)
	}
	s.imd.Clear()
	s.imd.Reset()
	s.WorldMap.Render(s.imd)
	s.imd.Draw(win)
}

func newScene(sim *sim.Simulation, bounds pixel.Rect) *Scene {
	wMap := newWorldMap(sim.World, bounds.Resized(topLeft(bounds), bounds.Size().Scaled(0.8)).Moved(pixel.V(scenePadding, -scenePadding)))
	s := &Scene{
		WorldMap: wMap,
		Sim:      sim,
	}
	return s
}
