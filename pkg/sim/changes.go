package sim

type change interface {
	Apply()
}

type clearReg struct {
	Reg *Region
}

func (e *clearReg) Apply() {
	e.Reg.Content = RCNone
	e.Reg.Bot = nil
}

type putBot struct {
	Bot *Bot
	Reg *Region
	Pos int
}

func (e *putBot) Apply() {
	e.Bot.Pos = e.Pos
	e.Reg.Bot = e.Bot
	e.Reg.Content = RCBot
}
