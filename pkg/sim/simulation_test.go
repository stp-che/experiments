package sim

import "testing"

type testBotGenome struct {
	actions []*Action
	i       int
}

func (c *testBotGenome) SetActions(actions []*Action) {
	c.actions = actions
	c.i = 0
}

func (c *testBotGenome) NextAction(_ *World, _ int) *Action {
	a := c.actions[c.i%len(c.actions)]
	c.i++
	return a
}
func (c *testBotGenome) EnergyCost() int {
	return 0
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

func prepare(c testSimConfig) (*Simulation, []*testBotGenome) {
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
	genomes := make([]*testBotGenome, botsCount)
	bots := make([]*Bot, botsCount)
	for i := 0; i < botsCount; i++ {
		genomes[i] = &testBotGenome{}
		bots[i] = (&Bot{Genome: genomes[i]}).Init(world)
		pos := c.BotsPos[i]
		(&putBot{
			Bot: bots[i],
			Reg: world.Regions[pos],
			Pos: pos,
		}).Apply()
	}

	sim := &Simulation{
		World: world,
		Bots:  bots,
	}

	return sim, genomes
}

type moveTestCase struct {
	Desc       string
	InitialPos int
	Dir        Direction
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
	c := func(s string, p1 int, dir Direction, p2 int) *moveTestCase {
		return &moveTestCase{
			Desc:       s,
			InitialPos: p1,
			Dir:        dir,
			NewPos:     p2,
		}
	}
	cases := []*moveTestCase{
		c("Going UpLeft", 12, UpLeft, 6),
		c("Going Up", 12, Up, 7),
		c("Going UpRight", 12, UpRight, 8),
		c("Going Right", 12, Right, 13),
		c("Going DownRight", 12, DownRight, 18),
		c("Going Down", 12, Down, 17),
		c("Going DownLeft", 12, DownLeft, 16),
		c("Going Left", 12, Left, 11),
		c("Going UpLeft out of bounds", 0, UpLeft, 0),
		c("Going Up out of bounds", 1, Up, 1),
		c("Going UpRight out of bounds", 1, UpRight, 1),
		c("Going Right out of bounds", 9, Right, 9),
		c("Going DownRight out of bounds", 9, DownRight, 9),
		c("Going Down out of bounds", 22, Down, 22),
		c("Going DownLeft out of bounds", 22, DownLeft, 22),
		c("Going Left out of bounds", 15, Left, 15),
	}
	for _, c := range cases {
		sim, genomes := prepare(testSimConfig{BotsPos: []int{c.InitialPos}})
		genomes[0].SetActions([]*Action{
			{Type: AMove, Direction: c.Dir},
		})
		sim.Step()
		checkWorldState(t, sim.World, testWorldState{
			Bots: map[int]*Bot{c.NewPos: sim.Bots[0]},
		})
	}

	cases = []*moveTestCase{
		c("Going UpLeft into the wall", 12, UpLeft, 12),
		c("Going Up into the wall", 12, Up, 12),
		c("Going UpRight into the wall", 12, UpRight, 12),
		c("Going Right into the wall", 12, Right, 12),
		c("Going DownRight into the wall", 12, DownRight, 12),
		c("Going Down into the wall", 12, Down, 12),
		c("Going DownLeft into the wall", 12, DownLeft, 12),
		c("Going Left into the wall", 12, Left, 12),
	}
	for _, c := range cases {
		walls := []int{6, 7, 8, 11, 13, 16, 17, 18}
		sim, genomes := prepare(testSimConfig{
			Walls:   walls,
			BotsPos: []int{c.InitialPos},
		})
		genomes[0].SetActions([]*Action{
			{Type: AMove, Direction: c.Dir},
		})
		sim.Step()
		checkWorldState(t, sim.World, testWorldState{
			Walls: walls,
			Bots:  map[int]*Bot{c.NewPos: sim.Bots[0]},
		})
	}

	sim, genomes := prepare(testSimConfig{BotsPos: []int{0, 1}})
	genomes[0].SetActions([]*Action{
		{Type: AMove, Direction: Right},
	})
	genomes[1].SetActions([]*Action{
		{Type: AMove, Direction: Down},
	})
	sim.Step()
	checkWorldState(t, sim.World, testWorldState{
		Bots: map[int]*Bot{0: sim.Bots[0], 6: sim.Bots[1]},
	})

	sim, genomes = prepare(testSimConfig{BotsPos: []int{0, 1}})
	genomes[0].SetActions([]*Action{
		{Type: AMove, Direction: DownRight},
	})
	genomes[1].SetActions([]*Action{
		{Type: AMove, Direction: Down},
	})
	sim.Step()
	checkWorldState(t, sim.World, testWorldState{
		Bots: map[int]*Bot{0: sim.Bots[0], 1: sim.Bots[1]},
	})

	sim, genomes = prepare(testSimConfig{BotsPos: []int{0}, Food: []int{1}})
	genomes[0].SetActions([]*Action{
		{Type: AMove, Direction: Right},
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
	genomes[0].SetActions([]*Action{
		{Type: AEat, Direction: Up},
		{Type: AEat, Direction: Right},
		{Type: AEat, Direction: DownRight},
		{Type: AEat, Direction: Down},
	})
	sim.Step()
	assertBotEnergy(t, sim.Bots[0], 300)
	sim.Step()
	assertBotEnergy(t, sim.Bots[0], 300)
	sim.Step()
	assertBotEnergy(t, sim.Bots[0], 300)
	sim.Step()
	assertBotEnergy(t, sim.Bots[0], 500)
	checkWorldState(t, sim.World, testWorldState{
		Walls: []int{1},
		Bots:  map[int]*Bot{0: sim.Bots[0]},
	})

	sim, genomes = prepare(testSimConfig{
		BotsPos: []int{0, 1, 5},
		Food:    []int{6},
	})
	genomes[0].SetActions([]*Action{
		{Type: AEat, Direction: DownRight},
	})
	genomes[1].SetActions([]*Action{
		{Type: AEat, Direction: Down},
	})
	genomes[2].SetActions([]*Action{
		{Type: AEat, Direction: Right},
	})
	sim.Step()
	assertBotEnergy(t, sim.Bots[0], 366)
	assertBotEnergy(t, sim.Bots[1], 366)
	assertBotEnergy(t, sim.Bots[2], 366)
}

func TestStep(t *testing.T) {
	testMoveActions(t)
	testEatActions(t)
}
