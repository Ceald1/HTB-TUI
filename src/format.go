package src

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// based on https://github.com/silofy/hackthebox/tree/master
	// Base background color
	BaseBG  = lipgloss.Color("#141D2B")

	// Primary palette
	Purple     = lipgloss.Color("#A4B1CD")
	Red        = lipgloss.Color("#FF8484")
	Cyan       = lipgloss.Color("#5CECC6")
	Pink       = lipgloss.Color("#CFBDFB")
	Yellow     = lipgloss.Color("#FFCC5C")
	LightBlue  = lipgloss.Color("#C5D1EB")
	Blue       = lipgloss.Color("#5CB2FF")
	LightGreen = lipgloss.Color("#C5F467")

	// Optional: Text color for dark backgrounds (use primary palette for accents)
	TextDefault = Purple
	TextRed     = Red
	TextCyan    = Cyan
	TextPink    = Pink
	TextYellow  = Yellow
	TextLightBlue = LightBlue
	TextBlue    = Blue
	TextLightGreen = LightGreen
	TextTitle = LightGreen

)

