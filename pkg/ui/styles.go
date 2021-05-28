package ui

import (
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

var (
	sceneBgColor = colornames.Azure
	scenePadding = 5.0

	worlMapRegionBorderColor = pixel.RGB(0.6, 0.6, 0.6)

	wallColor = colornames.Brown
	botColor  = colornames.Darkcyan
	foodColor = colornames.Green

	botWidth      = 65
	botHeight     = 40
	botMargin     = 3
	aliveBotColor = colornames.Darkcyan
	deadBotColor  = colornames.Darkgray
)
