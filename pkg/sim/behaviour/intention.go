package behaviour

import (
	"experiments/pkg/sim/core"
	"fmt"
	"math/rand"
)

type ActionType byte

func (t ActionType) String() string {
	switch t {
	case AMove:
		return "Move"
	case AEat:
		return "Eat"
	default:
		return fmt.Sprintf("%d", t)
	}
}

const (
	AMove ActionType = iota + 1
	AEat
)
const ActionTypesCount = 2

type Intention struct {
	ActionType ActionType
	Direction  core.Direction
}

func RandomIntention() *Intention {
	return &Intention{
		ActionType: randomActionType(),
		Direction:  core.RandomDirection(),
	}
}

func randomActionType() ActionType {
	return ActionType(rand.Intn(ActionTypesCount) + 1)
}

func actionType(n uint8) ActionType {
	a := n % ActionTypesCount
	if a == 0 {
		a = ActionTypesCount
	}
	return ActionType(a)
}
