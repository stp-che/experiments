package sim

import (
	"experiments/pkg/sim/behaviour"
	"experiments/pkg/sim/core"
	"experiments/pkg/test_helpers"
	"reflect"
	"testing"
)

type testBotBrain struct {
	intentions  []*behaviour.Intention
	visionRange int
	i           int
}

func (b *testBotBrain) SetActions(actions []*behaviour.Intention) {
	b.intentions = actions
	b.i = 0
}

func (b *testBotBrain) Process(_ []behaviour.OuterInput, _ behaviour.InnerInput) *behaviour.ProcessingResult {
	i := b.intentions[b.i%len(b.intentions)]
	b.i++
	return &behaviour.ProcessingResult{Decision: i}
}

func (b *testBotBrain) VisionRange() int {
	return b.visionRange
}

func (b *testBotBrain) Mutate(_ int) behaviour.IBrain {
	return b
}

func (b *testBotBrain) SetVisionRange(r int) {
	b.visionRange = r
}

type testSimConfig struct {
	W       int
	H       int
	BotsPos []int
	Walls   []int
	Food    []int
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func prepare(c testSimConfig) (*Experiment, []*testBotBrain) {
	world := newWorld(max(c.W, 5), max(c.H, 5))

	if c.Walls != nil {
		for _, pos := range c.Walls {
			world.Regions[pos].Content = RCWall
		}
	}

	if c.Food != nil {
		for _, pos := range c.Food {
			world.Regions[pos].Content = RCFood
		}
	}

	botsCount := len(c.BotsPos)
	brains := make([]*testBotBrain, botsCount)
	bots := make([]*Bot, botsCount)
	for i := 0; i < botsCount; i++ {
		brains[i] = &testBotBrain{}
		bots[i] = (&Bot{Brain: &BotBrain{Brain: brains[i]}}).Init(world)
		pos := c.BotsPos[i]
		(&putBot{
			Bot: bots[i],
			Reg: world.Regions[pos],
			Pos: pos,
		}).Apply()
	}

	ex := &Experiment{
		World: world,
		Bots:  bots,
	}

	return ex, brains
}

type moveTestCase struct {
	Desc       string
	InitialPos int
	Dir        core.Direction
	NewPos     int
	Before     func()
}

func assertBotPos(t *testing.T, bot *Bot, pos int) {
	if bot.Pos != pos {
		t.Errorf("Expected bot position is %d, got %d", pos, bot.Pos)
	}
}

func assertRegContainsBot(t *testing.T, pos int, reg *Region, bot *Bot) {
	if reg.Content != RCBot {
		t.Errorf("Expected region[%d] content is %d, got %d", pos, RCBot, reg.Content)
	}
	if reg.Bot != bot {
		t.Errorf("Expected region[%d] bot to be the bot", pos)
	}
}

func assertRegContentExceptBot(t *testing.T, pos int, reg *Region, content RegionContent) {
	if reg.Content != content {
		t.Errorf("Expected region[%d] content is %d, got %d", pos, content, reg.Content)
	}
	if reg.Bot != nil {
		t.Errorf("Expected region[%d] bot to be nil", pos)
	}
}

type testWorldState struct {
	Walls []int
	Bots  map[int]*Bot
}

func checkWorldState(t *testing.T, w *World, s testWorldState) {
	nonEmptyRegs := map[int]bool{}
	if s.Bots != nil {
		for pos, bot := range s.Bots {
			nonEmptyRegs[pos] = true
			assertBotPos(t, bot, pos)
			assertRegContainsBot(t, pos, w.Regions[pos], bot)
		}
	}
	if s.Walls != nil {
		for _, pos := range s.Walls {
			nonEmptyRegs[pos] = true
			assertRegContentExceptBot(t, pos, w.Regions[pos], RCWall)
		}
	}
	for pos, reg := range w.Regions {
		if _, nonEmpty := nonEmptyRegs[pos]; nonEmpty {
			continue
		}
		assertRegContentExceptBot(t, pos, reg, RCNone)
	}
}

func testMoveActions(t *testing.T) {
	c := func(s string, p1 int, dir core.Direction, p2 int) *moveTestCase {
		return &moveTestCase{
			Desc:       s,
			InitialPos: p1,
			Dir:        dir,
			NewPos:     p2,
		}
	}
	cases := []*moveTestCase{
		c("Going UpLeft", 12, core.UpLeft, 6),
		c("Going Up", 12, core.Up, 7),
		c("Going UpRight", 12, core.UpRight, 8),
		c("Going Right", 12, core.Right, 13),
		c("Going DownRight", 12, core.DownRight, 18),
		c("Going Down", 12, core.Down, 17),
		c("Going DownLeft", 12, core.DownLeft, 16),
		c("Going Left", 12, core.Left, 11),
		c("Going UpLeft out of bounds", 0, core.UpLeft, 0),
		c("Going Up out of bounds", 1, core.Up, 1),
		c("Going UpRight out of bounds", 1, core.UpRight, 1),
		c("Going Right out of bounds", 9, core.Right, 9),
		c("Going DownRight out of bounds", 9, core.DownRight, 9),
		c("Going Down out of bounds", 22, core.Down, 22),
		c("Going DownLeft out of bounds", 22, core.DownLeft, 22),
		c("Going Left out of bounds", 15, core.Left, 15),
	}
	for _, c := range cases {
		ex, genomes := prepare(testSimConfig{BotsPos: []int{c.InitialPos}})
		genomes[0].SetActions([]*behaviour.Intention{
			{ActionType: behaviour.AMove, Direction: c.Dir},
		})
		ex.Step()
		checkWorldState(t, ex.World, testWorldState{
			Bots: map[int]*Bot{c.NewPos: ex.Bots[0]},
		})
	}

	cases = []*moveTestCase{
		c("Going UpLeft into the wall", 12, core.UpLeft, 12),
		c("Going Up into the wall", 12, core.Up, 12),
		c("Going UpRight into the wall", 12, core.UpRight, 12),
		c("Going Right into the wall", 12, core.Right, 12),
		c("Going DownRight into the wall", 12, core.DownRight, 12),
		c("Going Down into the wall", 12, core.Down, 12),
		c("Going DownLeft into the wall", 12, core.DownLeft, 12),
		c("Going Left into the wall", 12, core.Left, 12),
	}
	for _, c := range cases {
		walls := []int{6, 7, 8, 11, 13, 16, 17, 18}
		sim, genomes := prepare(testSimConfig{
			Walls:   walls,
			BotsPos: []int{c.InitialPos},
		})
		genomes[0].SetActions([]*behaviour.Intention{
			{ActionType: behaviour.AMove, Direction: c.Dir},
		})
		sim.Step()
		checkWorldState(t, sim.World, testWorldState{
			Walls: walls,
			Bots:  map[int]*Bot{c.NewPos: sim.Bots[0]},
		})
	}

	sim, genomes := prepare(testSimConfig{BotsPos: []int{0, 1}})
	genomes[0].SetActions([]*behaviour.Intention{
		{ActionType: behaviour.AMove, Direction: core.Right},
	})
	genomes[1].SetActions([]*behaviour.Intention{
		{ActionType: behaviour.AMove, Direction: core.Down},
	})
	sim.Step()
	checkWorldState(t, sim.World, testWorldState{
		Bots: map[int]*Bot{0: sim.Bots[0], 6: sim.Bots[1]},
	})

	sim, genomes = prepare(testSimConfig{BotsPos: []int{0, 1}})
	genomes[0].SetActions([]*behaviour.Intention{
		{ActionType: behaviour.AMove, Direction: core.DownRight},
	})
	genomes[1].SetActions([]*behaviour.Intention{
		{ActionType: behaviour.AMove, Direction: core.Down},
	})
	sim.Step()
	checkWorldState(t, sim.World, testWorldState{
		Bots: map[int]*Bot{0: sim.Bots[0], 1: sim.Bots[1]},
	})

	sim, genomes = prepare(testSimConfig{BotsPos: []int{0}, Food: []int{1}})
	genomes[0].SetActions([]*behaviour.Intention{
		{ActionType: behaviour.AMove, Direction: core.Right},
	})
	sim.Step()
	checkWorldState(t, sim.World, testWorldState{
		Bots: map[int]*Bot{1: sim.Bots[0]},
	})
}

func assertBotEnergy(t *testing.T, bot *Bot, val int) {
	if bot.Energy != val {
		t.Errorf("Expected bot Energy to be %d, got %d", val, bot.Energy)
	}
}

func testEatActions(t *testing.T) {
	sim, genomes := prepare(testSimConfig{
		BotsPos: []int{0},
		Walls:   []int{1},
		Food:    []int{5},
	})
	genomes[0].SetActions([]*behaviour.Intention{
		{ActionType: behaviour.AEat, Direction: core.Up},
		{ActionType: behaviour.AEat, Direction: core.Right},
		{ActionType: behaviour.AEat, Direction: core.DownRight},
		{ActionType: behaviour.AEat, Direction: core.Down},
	})
	sim.Step()
	assertBotEnergy(t, sim.Bots[0], 500)
	sim.Step()
	assertBotEnergy(t, sim.Bots[0], 500)
	sim.Step()
	assertBotEnergy(t, sim.Bots[0], 500)
	sim.Step()
	assertBotEnergy(t, sim.Bots[0], 1500)
	checkWorldState(t, sim.World, testWorldState{
		Walls: []int{1},
		Bots:  map[int]*Bot{0: sim.Bots[0]},
	})

	sim, genomes = prepare(testSimConfig{
		BotsPos: []int{0, 1, 5},
		Food:    []int{6},
	})
	genomes[0].SetActions([]*behaviour.Intention{
		{ActionType: behaviour.AEat, Direction: core.DownRight},
	})
	genomes[1].SetActions([]*behaviour.Intention{
		{ActionType: behaviour.AEat, Direction: core.Down},
	})
	genomes[2].SetActions([]*behaviour.Intention{
		{ActionType: behaviour.AEat, Direction: core.Right},
	})
	sim.Step()
	assertBotEnergy(t, sim.Bots[0], 833)
	assertBotEnergy(t, sim.Bots[1], 833)
	assertBotEnergy(t, sim.Bots[2], 833)
}

func testNilActions(t *testing.T) {
	sim, genomes := prepare(testSimConfig{
		BotsPos: []int{0},
	})
	genomes[0].SetActions([]*behaviour.Intention{nil})
	sim.Step()
	checkWorldState(t, sim.World, testWorldState{
		Bots: map[int]*Bot{0: sim.Bots[0]},
	})
}

func TestStep(t *testing.T) {
	testMoveActions(t)
	testEatActions(t)
	testNilActions(t)
}

func TestBotsChart(t *testing.T) {
	e := Experiment{
		Bots: []*Bot{
			{Movements: 1, Age: 10, Energy: 0},
			{Movements: 1, Age: 100, Energy: 50},
			{Movements: 1, Age: 10, Energy: -5},
			{Movements: 1, Age: 100, Energy: 0},
			{Movements: 2, Age: 60, Energy: 0},
			{Movements: 100, Age: 50, Energy: 0},
			{Movements: 50, Age: 99, Energy: 0},
			{Movements: 2, Age: 99, Energy: 1},
		},
	}

	expectedChart := []*Bot{
		{Movements: 2, Age: 99, Energy: 0},
		{Movements: 50, Age: 99, Energy: 0},
		{Movements: 2, Age: 60, Energy: 0},
		{Movements: 100, Age: 50, Energy: 0},
		{Movements: 1, Age: 100, Energy: 50},
		{Movements: 1, Age: 100, Energy: 0},
		{Movements: 1, Age: 10, Energy: 0},
		{Movements: 1, Age: 10, Energy: -5},
	}

	actualChart := e.BotsChart()

	if !reflect.DeepEqual(expectedChart, actualChart) {
		t.Errorf("Expected chart to eq %s, got %s", test_helpers.Inspect(expectedChart), test_helpers.Inspect(actualChart))
	}
}
