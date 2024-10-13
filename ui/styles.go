package ui

import (
	"github.com/charmbracelet/lipgloss"
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

var Styles = AppStyles{
	LogMessageStyles: LogMessageStyles{
		LineCountStyle: lipgloss.NewStyle().Italic(true).Foreground(PastelYellow),
		OriginStyle:    lipgloss.NewStyle().Foreground(PastelPurple),
		TimestampStyle: lipgloss.NewStyle().Foreground(PastelBlue),
		LevelStyle:     lipgloss.NewStyle().Foreground(PastelGreen),
		MessageStyle:   lipgloss.NewStyle().Foreground(PastelWhite),
		CursorStyle:    lipgloss.NewStyle().Bold(true).Foreground(PastelRed),
	},
	NearbyLogEntryStyles: NearbyLogEntryStyles{
		OriginStyle:    lipgloss.NewStyle().Foreground(FadedPurple),
		TimestampStyle: lipgloss.NewStyle().Foreground(FadedGreen),
		LevelStyle:     lipgloss.NewStyle().Foreground(FadedBlue),
		MessageStyle:   lipgloss.NewStyle().Foreground(FadedGrey),
	},
}

type AppStyles struct {
	LogMessageStyles     LogMessageStyles
	NearbyLogEntryStyles NearbyLogEntryStyles
}

type LogMessageStyles struct {
	LineCountStyle lipgloss.Style
	OriginStyle    lipgloss.Style
	TimestampStyle lipgloss.Style
	LevelStyle     lipgloss.Style
	MessageStyle   lipgloss.Style
	CursorStyle    lipgloss.Style
}

type NearbyLogEntryStyles struct {
	OriginStyle    lipgloss.Style
	TimestampStyle lipgloss.Style
	LevelStyle     lipgloss.Style
	MessageStyle   lipgloss.Style
}
