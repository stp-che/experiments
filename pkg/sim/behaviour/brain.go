package behaviour

import (
	"experiments/pkg/test_helpers"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

type IBrain interface {
	Process([]OuterInput, InnerInput) *ProcessingResult
	VisionRange() int
	Mutate(int) IBrain
}

type Brain struct {
	uuid                uuid.UUID
	origUuid            uuid.UUID
	OuterReceptor       OuterReceptor
	HealthAnalyzerNet   HealthAnalyzerNet
	OuterAnalyzersCount int
	OuterAnalyzerNet    OuterAnalyzerNet
	ManipulationSystem  ManipulationSystem
}

func (b *Brain) Process(o []OuterInput, i InnerInput) *ProcessingResult {
	defer func() {
		if r := recover(); r != nil {
			msg := b.strf("PANIC %v\n\nBrain content:\n\n%s\n", r, test_helpers.Inspect(b))
			fmt.Print(msg)
			fatalf(msg)
		}
	}()
	activation := ManipulationSystemActivation{}
	correction := b.HealthAnalyzerNet.Correction(i[0])
	for _, inp := range o {
		collectedSignal := b.OuterReceptor.CollectSignal(inp.Signal)
		activation[inp.Direction] = b.OuterAnalyzerNet.Activation(collectedSignal, correction)
	}
	return &ProcessingResult{
		Decision:   b.ManipulationSystem.ComputeIntention(activation),
		EnergyCost: b.energyCost(activation),
	}
}

func (b *Brain) VisionRange() int {
	return int(b.OuterReceptor.visionRange)
}

func (b *Brain) energyCost(activations ManipulationSystemActivation) int {
	baseCost := b.OuterReceptor.Size() +
		len(b.HealthAnalyzerNet) +
		len(b.OuterAnalyzerNet) +
		len(b.ManipulationSystem) +
		b.OuterAnalyzersCount + 1 // number of outer analyzers and one inner analyzer
	activityCost := 0
	for _, activation := range activations {
		for _, pow := range activation {
			if pow < 0 {
				pow = -pow
			}
			activityCost += int(pow)
		}
	}
	return baseCost + (activityCost-1)/32 + 1
}

func (b *Brain) Mutate(n int) IBrain {
	debugf(b.strf("Mutate(%d) START", n))
	defer debugf(b.strf("Mutate END"))

	newBrain := b
	for i := 0; i < n; i++ {
		prevBrain := newBrain
		mutation := newBrain.randomMutation()
		debugf(prevBrain.strf("%#v", mutation))
		newBrain = mutation.apply(newBrain)
		debugf(prevBrain.strf(" -> [%s] %s", newBrain.uuid, newBrain))
	}
	return newBrain
}

func (b *Brain) strf(pattern string, args ...interface{}) string {
	args = append([]interface{}{b.origUuid, b.uuid}, args...)
	return fmt.Sprintf("%s [%s] "+pattern, args...)
}

func (b *Brain) copy() *Brain {
	return &Brain{
		uuid:                uuid.New(),
		origUuid:            b.origUuid,
		OuterReceptor:       b.OuterReceptor,
		HealthAnalyzerNet:   b.HealthAnalyzerNet,
		OuterAnalyzersCount: b.OuterAnalyzersCount,
		OuterAnalyzerNet:    b.OuterAnalyzerNet,
		ManipulationSystem:  b.ManipulationSystem,
	}
}

func (b *Brain) String() string {
	pattern := "{uuid: %s, origUuid: %s, OuterReceptor: %s, HealthAnalyzerNet: %s, OuterAnalyzersCount: %d, OuterAnalyzerNet: %s, ManipulationSystem: %s}"
	return fmt.Sprintf(
		pattern,
		b.uuid.String(),
		b.origUuid.String(),
		b.OuterReceptor.String(),
		b.HealthAnalyzerNet.String(),
		b.OuterAnalyzersCount,
		b.OuterAnalyzerNet.String(),
		b.ManipulationSystem.String(),
	)
}

func RandomBrain() *Brain {
	b := &Brain{
		uuid:                uuid.New(),
		OuterAnalyzersCount: rand.Intn(10) + 1,
		ManipulationSystem:  randomManipulationSystem(),
	}
	b.origUuid = b.uuid
	b.OuterReceptor = randomOuterReceptor(b.OuterAnalyzersCount)
	b.OuterAnalyzerNet = randomOuterAnalyzerNet(b.OuterAnalyzersCount, len(b.ManipulationSystem))
	b.HealthAnalyzerNet = randomHealthAnalyzerNet(len(b.OuterAnalyzerNet))
	return b
}

func (b *Brain) randomMutation() iMutation {
	switch rand.Intn(11) {
	case 0:
		return mIncreaseVisionRange{}
	case 1:
		return mDecreaseVisionRange{}
	case 2:
		return randomChangeOuterReceptor(b)
	case 3:
		return randomAddHealthAnalyzerLink(b)
	case 4:
		return randomChangeHealthAnalyzerCorrection(b)
	case 5:
		return randomChangeHealthAnalyzerMinMax(b)
	case 6:
		return randomChangeOuterAnalyzerLinkPower(b)
	case 7:
		return randomChangeManipulatorValue(b)
	case 8:
		return reductiveMutation{}
	case 9:
		return mAddOuterAnalyzer{}
	case 10:
		return randomAddManipulatorMutation()
	default:
		return randomAddManipulatorMutation()
	}
}
