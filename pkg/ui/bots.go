package ui

import (
	"experiments/pkg/sim"
	"fmt"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type BotsComponent struct {
	groups  []*BotsGroupView
	TopLeft pixel.Vec
}

func (b *BotsComponent) Init(groups []sim.BotsGroup) *BotsComponent {
	b.groups = make([]*BotsGroupView, len(groups))
	top := b.TopLeft.Y
	for i, g := range groups {
		b.groups[i] = (&BotsGroupView{TopLeft: pixel.V(b.TopLeft.X, top)}).Init(g)
		top -= b.groups[i].Height() + 20
	}
	return b
}

func (b *BotsComponent) Render(imd *imdraw.IMDraw) {
	for _, g := range b.groups {
		g.Render(imd)
	}
}

func (b *BotsComponent) RenderText(win *pixelgl.Window) {
	for _, g := range b.groups {
		g.RenderText(win)
	}
}

type BotsGroupView struct {
	bots    []*BotView
	TopLeft pixel.Vec
}

func (b *BotsGroupView) Init(group sim.BotsGroup) *BotsGroupView {
	botOuterWidth, botOuterHeight := botWidth+botMargin, botHeight+botMargin
	b.bots = make([]*BotView, len(group.Bots))
	for i, bot := range group.Bots {
		col, row := i%4, i/4
		b.bots[i] = &BotView{
			Bot:     bot,
			TopLeft: b.TopLeft.Add(pixel.V(float64(botOuterWidth*col), -float64(botOuterHeight*row))),
		}
	}
	return b
}

func (b *BotsGroupView) Height() float64 {
	return float64((len(b.bots) + 1) / 4 * (botHeight + botMargin))
}

func (b *BotsGroupView) Render(imd *imdraw.IMDraw) {
	for _, bot := range b.bots {
		bot.Render(imd)
	}
}

func (b *BotsGroupView) RenderText(win *pixelgl.Window) {
	for _, bot := range b.bots {
		bot.RenderText(win)
	}
}

type BotView struct {
	Bot       *sim.Bot
	TopLeft   pixel.Vec
	textAtlas *text.Atlas
}

func (b *BotView) Render(imd *imdraw.IMDraw) {
	imd.Color = aliveBotColor
	if !b.Bot.IsAlive() {
		imd.Color = deadBotColor
	}
	imd.Push(b.TopLeft)
	imd.Push((b.TopLeft.Add(pixel.V(float64(botWidth)-1, -float64(botHeight)+1))))
	imd.Rectangle(0)
}

func (b *BotView) RenderText(win *pixelgl.Window) {
	b.putString(win, pixel.V(3, 14), strconv.Itoa(b.Bot.Age))
	b.putString(win, pixel.V(3, 34), strconv.Itoa(b.Bot.Energy))
}

func (b *BotView) putString(win *pixelgl.Window, at pixel.Vec, str string) {
	if b.textAtlas == nil {
		b.textAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
	}
	at.Y *= -1
	basicTxt := text.New(b.TopLeft.Add(at), b.textAtlas)
	fmt.Fprintln(basicTxt, str)
	basicTxt.Draw(win, pixel.IM)
}
