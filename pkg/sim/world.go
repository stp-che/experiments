package sim

import "log"

type RegionContent int

const (
	RCNone RegionContent = iota + 1
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
}

func (w *World) Region(p Pos) *Region {
	if p[0] < 0 || p[0] >= w.Cols || p[1] < 0 || p[1] >= w.Rows {
		return nil
	}
	return w.Regions[p[1]*w.Cols+p[0]]
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
}

func newWorld(cols, rows int) *World {
	w := &World{
		Cols: cols,
		Rows: rows,
	}
	w.Init()
	return w
}
