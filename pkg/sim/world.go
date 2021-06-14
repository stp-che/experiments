package sim

import (
	"experiments/pkg/sim/core"
	"fmt"
	"math/rand"
	"time"
)

type RegionContent byte

const (
	RCNone RegionContent = iota + 1
	RCWall
	RCFood
	RCBot
)

func (c RegionContent) String() string {
	switch c {
	case RCNone:
		return "None"
	case RCWall:
		return "Wall"
	case RCFood:
		return "Food"
	case RCBot:
		return "Bot"
	default:
		return fmt.Sprintf("%d", c)
	}
}

type Region struct {
	Content RegionContent
	Bot     *Bot
}

func (r *Region) Busy() bool {
	return r.Content == RCWall || r.Content == RCBot
}

type World struct {
	Cols    int
	Rows    int
	Regions []*Region
	// Boundless world means that going beyond the edge we appear on the opposite side
	Boundless bool
}

func (w *World) PossibleDirection(p int, d core.Direction) bool {
	return w.NextPos(p, d) >= 0
}

func (w *World) Init() {
	if w.Regions != nil {
		return
	}
	n := w.Cols * w.Rows
	w.Regions = make([]*Region, n)
	for i := 0; i < n; i++ {
		w.Regions[i] = &Region{Content: RCNone}
	}
	buildWalls(w, n/30)
}

func buildWalls(w *World, totalLenght int) {
	if totalLenght == 0 {
		return
	}
	// TODO: check if totalLength >= world size
	built := 0
	building := true
	keepingAction := 0
	randomWalk(w, 80, func(r *Region) bool {
		if r.Content == RCNone && building {
			r.Content = RCWall
			built++
		}
		keepingAction++
		i := rand.Intn((w.Cols + w.Rows) / 2)
		if building && i/5 < keepingAction || i < keepingAction {
			building = !building
			keepingAction = 0
		}
		return built < totalLenght
	})
}

func randomWalk(w *World, keepingDirection int, step func(*Region) bool) {
	currentPos := w.RandomPos()
	dirs := [4]core.Direction{core.Up, core.Right, core.Down, core.Left}
	currentDir := core.Up
	dirsChoice := make([]core.Direction, 0, 4)
	for step(w.Regions[currentPos]) {
		// fmt.Println("---------------")
		// fmt.Printf("w.PossibleDirection(%v, %v): %v\n", currentPos, currentDir, w.PossibleDirection(currentPos, currentDir))
		if !w.PossibleDirection(currentPos, currentDir) || rand.Intn(100) > keepingDirection {
			// fmt.Println("change direction")
			dirsChoice = dirsChoice[:0]
			for _, d := range dirs {
				if w.PossibleDirection(currentPos, d) && d != currentDir {
					dirsChoice = append(dirsChoice, d)
				}
			}
			// fmt.Printf("    dirsChoice: %v\n", dirsChoice)
			currentDir = dirsChoice[rand.Intn(len(dirsChoice))]
			// fmt.Printf("    new direction: %v\n", currentDir)
		}
		currentPos = w.NextPos(currentPos, currentDir)
	}
}

func (w *World) Region(pos int) *Region {
	if pos < 0 || pos > len(w.Regions) {
		return nil
	}
	return w.Regions[pos]
}

func (w *World) RandomPos() int {
	return rand.Intn(w.Cols * w.Rows)
}

// Returns position of region next to the given one in given direction
// Returns -1 if it is out of world bounds
func (w *World) NextPos(pos int, d core.Direction) int {
	dx, dy := d.DeltaXY()
	return w.ShiftXY(pos, dx, dy)
}

func (w *World) XYPos(n int) (int, int) {
	// TODO: check bounds
	return n % w.Cols, n / w.Cols
}

func (w *World) ShiftXY(pos, dx, dy int) int {
	x, y := w.XYPos(pos)
	if x+dx < 0 || x+dx >= w.Cols || y+dy < 0 || y+dy >= w.Rows {
		return -1
	}
	return pos + dy*w.Cols + dx
}

func (w *World) RandomEmptyPositions(n int) []int {
	if n == 0 {
		return []int{}
	}

	wSize := w.Cols * w.Rows
	freeRegIdxs := make([]int, 0, wSize)
	for i, r := range w.Regions {
		if r.Content == RCNone {
			freeRegIdxs = append(freeRegIdxs, i)
		}
	}
	if n >= len(freeRegIdxs) {
		return freeRegIdxs
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(freeRegIdxs), func(i, j int) {
		freeRegIdxs[i], freeRegIdxs[j] = freeRegIdxs[j], freeRegIdxs[i]
	})
	return freeRegIdxs[:n]
}

func newWorld(cols, rows int) *World {
	w := &World{
		Cols: cols,
		Rows: rows,
	}
	w.Init()
	return w
}
