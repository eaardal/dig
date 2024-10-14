package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/eaardal/dig/utils"
	"math/rand"
)

var PastelPink = lipgloss.Color("#F78DA7")
var PastelBlue = lipgloss.Color("#A7C7E7")
var PastelGreen = lipgloss.Color("#A8E6CF")
var PastelYellow = lipgloss.Color("#FFF6A5")
var PastelPurple = lipgloss.Color("#C4A7E7")
var PastelOrange = lipgloss.Color("#FFD3B5")
var PastelTeal = lipgloss.Color("#98D5C6")
var PastelRed = lipgloss.Color("#F4978E")
var PastelLavender = lipgloss.Color("#E8DFF5")
var PastelGrey = lipgloss.Color("#D9D9D9")
var PastelWhite = lipgloss.Color("#F0F0F0")

var FadedPink = lipgloss.Color("#E6A2B1")
var FadedBlue = lipgloss.Color("#B3C7D8")
var FadedGreen = lipgloss.Color("#B2CDB9")
var FadedYellow = lipgloss.Color("#EDE4B2")
var FadedPurple = lipgloss.Color("#B7A9C8")
var FadedOrange = lipgloss.Color("#EECAB7")
var FadedTeal = lipgloss.Color("#A4BEB9")
var FadedRed = lipgloss.Color("#E5A39B")
var FadedLavender = lipgloss.Color("#D6CCDD")
var FadedGrey = lipgloss.Color("#C2C2C2")
var FadedWhite = lipgloss.Color("#E0E0E0")

var AllPastelColors = []lipgloss.Color{
	PastelPink,
	PastelBlue,
	PastelGreen,
	PastelYellow,
	PastelPurple,
	PastelOrange,
	PastelTeal,
	PastelRed,
	PastelLavender,
	PastelGrey,
}

var AllFadedColors = []lipgloss.Color{
	FadedPink,
	FadedBlue,
	FadedGreen,
	FadedYellow,
	FadedPurple,
	FadedOrange,
	FadedTeal,
	FadedRed,
	FadedLavender,
	FadedGrey,
}

var AllColors = append(AllPastelColors, AllFadedColors...)

func RandomPastelColor() lipgloss.Color {
	return AllPastelColors[rand.Intn(len(AllPastelColors))]
}

func RandomFadedColor() lipgloss.Color {
	return AllFadedColors[rand.Intn(len(AllFadedColors))]
}

func GetPastelColorForValue(value string) lipgloss.Color {
	return utils.DeterministicItemForValue(value, AllPastelColors)
}

var Styles = AppStyles{
	Cursor: lipgloss.NewStyle().Bold(true).Foreground(PastelRed),
	LogEntry: LogMessageStyles{
		LineNumber: lipgloss.NewStyle().Italic(true).Foreground(PastelYellow),
		Origin:     lipgloss.NewStyle().Foreground(PastelPurple),
		Timestamp:  lipgloss.NewStyle().Foreground(PastelBlue),
		Level:      lipgloss.NewStyle().Foreground(PastelGreen),
		Message:    lipgloss.NewStyle().Foreground(PastelWhite),
	},
	NearbyLogEntry: NearbyLogEntryStyles{
		Origin:    lipgloss.NewStyle().Foreground(FadedPurple),
		Timestamp: lipgloss.NewStyle().Foreground(FadedGreen),
		Level:     lipgloss.NewStyle().Foreground(FadedBlue),
		Message:   lipgloss.NewStyle().Foreground(FadedGrey),
	},
}

type AppStyles struct {
	LogEntry       LogMessageStyles
	NearbyLogEntry NearbyLogEntryStyles
	Cursor         lipgloss.Style
}

type LogMessageStyles struct {
	LineNumber lipgloss.Style
	Origin     lipgloss.Style
	Timestamp  lipgloss.Style
	Level      lipgloss.Style
	Message    lipgloss.Style
}

type NearbyLogEntryStyles struct {
	Origin    lipgloss.Style
	Timestamp lipgloss.Style
	Level     lipgloss.Style
	Message   lipgloss.Style
}
