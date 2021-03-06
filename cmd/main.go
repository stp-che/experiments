package main

import (
	"log"
	"math/rand"
	"time"

	"experiments/pkg/sim"
	"experiments/pkg/ui"

	"github.com/faiface/pixel/pixelgl"
)

var simConfig = sim.Config{
	ExperimentsNumber: 500,
	WorldWidth:        100,
	WorldHeight:       60,
	BrainsNumber:      8,
	GroupSize:         8,
	MutantsPerGroup:   2,
	FoodAmount:        100,
}

func run() {
	rand.Seed(time.Now().UnixNano())
	sim, err := sim.NewSimulation(simConfig)
	if err != nil {
		log.Fatalln(err)
	}
	simUi, err := ui.New(sim)
	if err != nil {
		log.Fatalln(err)
	}

	simUi.Update()

	simUpdates, err := sim.Run()
	if err != nil {
		log.Fatalln(err)
	}

	tick := time.Tick(time.Millisecond * 100)

	for !simUi.Closed() {
		select {
		case <-simUpdates:
			simUi.Update()
		case <-tick:
			simUi.Win.Update()
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func main() {
	pixelgl.Run(run)
}
