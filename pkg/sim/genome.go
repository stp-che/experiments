package sim

import "math/rand"

type iGenome interface {
	NextAction(*World, int) *Action
	EnergyCost() int
}

type ControllerGene struct {
	ActionType ActionType
	DirValues  [8]int8
}

type Genome struct {
	Controllers []*ControllerGene
}

func (g *Genome) NextAction(world *World, pos int) *Action {
	activity := make([]int16, len(g.Controllers))
	for i := 0; i < len(activity); i++ {
		activity[i] = int16(rand.Intn(255) - 100)
	}
	return genomeComputeAction(g.Controllers, activity)
}

func genomeComputeAction(genes []*ControllerGene, activity []int16) *Action {
	actionsTable := make(map[ActionType]*[8]int16)
	// TODO: check genes and activity have the same len
	for i, gene := range genes {
		actionTotal, ok := actionsTable[gene.ActionType]
		if !ok {
			actionTotal = &[8]int16{}
			actionsTable[gene.ActionType] = actionTotal
		}
		for j, v := range gene.DirValues {
			actionTotal[j] += int16(v) * activity[i]
		}
	}
	var max int16 = 0
	var action *Action
	for aType, values := range actionsTable {
		for i, v := range values {
			if max < v {
				if action == nil {
					action = &Action{}
				}
				action.Type = aType
				action.Direction = Direction(i + 1)
				max = v
			}
		}
	}
	return action
}

func (g *Genome) EnergyCost() int {
	return 10
}

func randomGenome() *Genome {
	controllers := make([]*ControllerGene, rand.Intn(5)+1)
	for i := 0; i < len(controllers); i++ {
		dirValues := [8]int8{}
		for i := 0; i < 8; i++ {
			dirValues[i] = int8(rand.Intn(255) - 128)
		}
		controllers[i] = &ControllerGene{
			ActionType: randomActionType(),
			DirValues:  dirValues,
		}
	}
	return &Genome{
		Controllers: controllers,
	}
}
