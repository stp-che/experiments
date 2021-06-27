package sim

import (
	"experiments/pkg/sim/behaviour"
	"sync/atomic"
)

type BotBrain struct {
	ID         uint64
	Generation uint64
	Brain      behaviour.IBrain
}

var botBrainIDCounter uint64 = 1

func nextBotBrainID() uint64 {
	defer atomic.AddUint64(&botBrainIDCounter, 1)
	return botBrainIDCounter
}

func randomBotBrain() *BotBrain {
	return &BotBrain{
		ID:         nextBotBrainID(),
		Generation: 1,
		Brain:      behaviour.RandomBrain(),
	}
}

func (b *BotBrain) Mutate(n int) *BotBrain {
	return &BotBrain{
		ID:         nextBotBrainID(),
		Generation: b.Generation + 1,
		Brain:      b.Brain.Mutate(n),
	}
}

func (b *BotBrain) Process(o []behaviour.OuterInput, i behaviour.InnerInput) *behaviour.ProcessingResult {
	return b.Brain.Process(o, i)
}

func (b *BotBrain) VisionRange() int {
	return b.Brain.VisionRange()
}
