package main

import (
	"log"
	"math/rand"
	"time"

	"experiments/pkg/sim"
	"experiments/pkg/sim/behaviour"
	"experiments/pkg/ui"

	"github.com/faiface/pixel/pixelgl"
	"go.uber.org/zap"
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

func setupLogger() (*zap.Logger, *zap.SugaredLogger) {
	cfg := zap.NewDevelopmentConfig()
	cfg.DisableCaller = true
	cfg.OutputPaths = []string{"log/sim.log"}
	logger, _ := cfg.Build()
	sugar := logger.Sugar()
	behaviour.SetLogger(sugar)
	return logger, sugar
}

func run() {
	// Uncomment this to enable logging
	// logger, sugar := setupLogger()
	// defer logger.Sync()
	// sugar.Infof("Start")

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
