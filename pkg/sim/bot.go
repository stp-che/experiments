package sim

import (
	"math/rand"
	"time"
)

type Bot struct {
	world *World
	Pos   Pos
}

func (b *Bot) Settle(world *World, pos Pos) {
	reg := world.Region(pos)
	if reg == nil {
		return
	}
	b.world = world
	b.Pos = pos
	reg.Occupy(b)
}

func (b *Bot) Move() bool {
	newPos := b.Pos.Next(randomDirection())
	newReg := b.world.Region(newPos)
	if newReg == nil {
		return false
	}
	b.world.Region(b.Pos).Clear()
	b.Pos = newPos
	newReg.Occupy(b)
	return true
}

func randomDirection() Direction {
	rand.Seed(time.Now().UnixNano())
	return Direction(rand.Intn(7) + 1)
}
