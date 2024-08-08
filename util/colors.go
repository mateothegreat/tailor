package util

import (
	"github.com/fatih/color"
	"golang.org/x/exp/rand"
)

var Colors = []color.Attribute{
	color.FgGreen, color.FgMagenta,
	color.FgBlue, color.FgYellow, color.FgCyan,
	color.FgHiGreen, color.FgHiYellow,
	color.FgHiBlue, color.FgHiMagenta, color.FgHiCyan,
}

func GetByInt(i int) color.Attribute {
	return Colors[i%len(Colors)]
}

func RandomColor() color.Attribute {
	colors := []color.Attribute{
		color.FgGreen, color.FgYellow,
		color.FgBlue, color.FgMagenta, color.FgCyan,
		color.FgHiGreen, color.FgHiYellow,
		color.FgHiBlue, color.FgHiMagenta, color.FgHiCyan,
	}
	return colors[rand.Intn(len(colors))]
}
