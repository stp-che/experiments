package sim

import (
	"log"
	"math/rand"
	"time"
)

type RegionContent int

const (
	RCNone RegionContent = iota
	RCWall
	RCFood
	RCBot
)

type Direction int

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

type Pos [2]int

func (p Pos) Next(d Direction) Pos {
	switch d {
	case UpLeft:
		return Pos{p[0] - 1, p[1] - 1}
	case Up:
		return Pos{p[0], p[1] - 1}
	case UpRight:
		return Pos{p[0] + 1, p[1] - 1}
	case Right:
		return Pos{p[0] + 1, p[1]}
	case DownRight:
		return Pos{p[0] + 1, p[1] + 1}
	case Down:
		return Pos{p[0], p[1] + 1}
	case DownLeft:
		return Pos{p[0] - 1, p[1] + 1}
	case Left:
		return Pos{p[0] - 1, p[1]}
	default:
		log.Fatalln("wrong direction")
		return p
	}
}

type Region struct {
	Content RegionContent
	Bot     *Bot
}

func (r *Region) Clear() {
	r.Content = RCNone
	r.Bot = nil
}

func (r *Region) Busy() bool {
	return r.Content == RCWall || r.Content == RCBot
}

func (r *Region) Occupy(bot *Bot) {
	r.Content = RCBot
	r.Bot = bot
}

type World struct {
	Cols    int
	Rows    int
	Regions []*Region
	// Boundless world means that going beyond the edge we appear on the opposite side
	Boundless bool
}

func (w *World) Region(p Pos) *Region {
	if !w.IncludesPos(p) {
		return nil
	}
	return w.Regions[p[1]*w.Cols+p[0]]
}

func (w *World) IncludesPos(p Pos) bool {
	return p[0] >= 0 && p[0] < w.Cols && p[1] >= 0 && p[1] < w.Rows
}

func (w *World) PossibleDirection(p Pos, d Direction) bool {
	return w.IncludesPos(p.Next(d))
}

func (w *World) Init() {
	if w.Regions != nil {
		return
	}
	n := w.Cols * w.Rows
	w.Regions = make([]*Region, n)
	for i := 0; i < n; i++ {
		w.Regions[i] = &Region{}
	}
	buildWalls(w, n/30)
}

func buildWalls(w *World, totalLenght int) {
	if totalLenght == 0 {
		return
	}
	rand.Seed(time.Now().UnixNano())
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
	rand.Seed(time.Now().UnixNano())
	currentPos := w.RandomPos()
	dirs := [4]Direction{Up, Right, Down, Left}
	currentDir := Up
	dirsChoice := make([]Direction, 0, 4)
	for step(w.Region(currentPos)) {
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
		currentPos = currentPos.Next(currentDir)
	}
}

func (w *World) RandomPos() Pos {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(w.Cols * w.Rows)
	return Pos{n % w.Cols, n / w.Cols}
}

func newWorld(cols, rows int) *World {
	w := &World{
		Cols: cols,
		Rows: rows,
	}
	w.Init()
	return w
}
