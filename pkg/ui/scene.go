package ui

import (
	"experiments/pkg/sim"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Scene struct {
	WorldMap *WorldMap
	BotsInfo *BotsComponent
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
	s.BotsInfo.Render(s.imd)
	s.imd.Draw(win)
	s.BotsInfo.RenderText(win)

	// basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	// basicTxt := text.New(pixel.V(100, 500), basicAtlas)
	// basicTxt.Color = colornames.Black

	// fmt.Fprintln(basicTxt, "Hello, text!")
	// fmt.Fprintln(basicTxt, "I support multiple lines!")
	// fmt.Fprintf(basicTxt, "And I'm an %s, yay!", "io.Writer")
	// basicTxt.Draw(win, pixel.IM)
}

func newScene(sim *sim.Simulation, bounds pixel.Rect) *Scene {
	wMap := newWorldMap(sim.World, bounds.Resized(topLeft(bounds), bounds.Size().Scaled(0.8)).Moved(pixel.V(scenePadding, -scenePadding)))
	botsInfo := &BotsComponent{
		TopLeft: topLeft(bounds).Add(pixel.V(wMap.Bounds().Max.X+50, -scenePadding)),
	}
	s := &Scene{
		WorldMap: wMap,
		BotsInfo: botsInfo.Init(sim.BotsGroups()),
		Sim:      sim,
	}
	return s
}
