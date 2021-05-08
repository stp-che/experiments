package ui

import (
	"experiments/pkg/sim"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	HEIGHT = 900
	WIDTH  = 1600
)

type Ui struct {
	Sim   *sim.Simulation
	Win   *pixelgl.Window
	scene *Scene
}

func New(sim *sim.Simulation) (*Ui, error) {
	cfg := pixelgl.WindowConfig{
		Title:  "Simulation",
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return nil, err
	}

	win.Clear(sceneBgColor)

	return &Ui{
		Sim: sim,
		Win: win,
	}, nil
}

func (ui *Ui) Update() {
	ui.Win.Clear(sceneBgColor)
	ui.drawScene()
	ui.Win.Update()
}

func (ui *Ui) Closed() bool {
	return ui.Win.Closed()
}

func (ui *Ui) drawScene() {
	if ui.scene == nil {
		ui.scene = newScene(ui.Sim, ui.Win.Bounds())
	}
	ui.scene.Draw(ui.Win)
}
