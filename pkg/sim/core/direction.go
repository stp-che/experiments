package core

import (
	"fmt"
	"log"
	"math/rand"
)

type Direction byte

const (
	UpLeft Direction = iota + 1
	Up
	UpRight
	Right
	DownRight
	Down
	DownLeft
	Left
)

func RandomDirection() Direction {
	return Direction(rand.Intn(7) + 1)
}

func (d Direction) DeltaXY() (int, int) {
	switch d {
	case UpLeft:
		return -1, -1
	case Up:
		return 0, -1
	case UpRight:
		return 1, -1
	case Right:
		return 1, 0
	case DownRight:
		return 1, 1
	case Down:
		return 0, 1
	case DownLeft:
		return -1, 1
	case Left:
		return -1, 0
	default:
		log.Fatalf("wrong direction %v\n", d)
		return 0, 0
	}
}

func (d Direction) String() string {
	switch d {
	case UpLeft:
		return "UpLeft"
	case Up:
		return "Up"
	case UpRight:
		return "UpRight"
	case Right:
		return "Right"
	case DownRight:
		return "DownRight"
	case Down:
		return "Down"
	case DownLeft:
		return "DownLeft"
	case Left:
		return "Left"
	default:
		return fmt.Sprintf("%d", d)
	}
}
