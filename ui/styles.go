package ui

import (
	"github.com/charmbracelet/lipgloss"
	"hash/fnv"
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
	PastelWhite,
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
	FadedWhite,
}

var AllColors = append(AllPastelColors, AllFadedColors...)

func RandomPastelColor() lipgloss.Color {
	return AllPastelColors[rand.Intn(len(AllPastelColors))]
}

func RandomFadedColor() lipgloss.Color {
	return AllFadedColors[rand.Intn(len(AllFadedColors))]
}

func hashStringToUint32(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

func RandomPastelColorForValue(value string) lipgloss.Color {
	hashValue := hashStringToUint32(value)
	colorIndex := int(hashValue) % len(AllPastelColors)
	return AllPastelColors[colorIndex]
}

var Styles = AppStyles{
	CursorStyle: lipgloss.NewStyle().Bold(true).Foreground(PastelRed),
	LogEntryStyles: LogMessageStyles{
		LineNumberStyle: lipgloss.NewStyle().Italic(true).Foreground(PastelYellow),
		OriginStyle:     lipgloss.NewStyle().Foreground(PastelPurple),
		TimestampStyle:  lipgloss.NewStyle().Foreground(PastelBlue),
		LevelStyle:      lipgloss.NewStyle().Foreground(PastelGreen),
		MessageStyle:    lipgloss.NewStyle().Foreground(PastelWhite),
	},
	NearbyLogEntryStyles: NearbyLogEntryStyles{
		OriginStyle:    lipgloss.NewStyle().Foreground(FadedPurple),
		TimestampStyle: lipgloss.NewStyle().Foreground(FadedGreen),
		LevelStyle:     lipgloss.NewStyle().Foreground(FadedBlue),
		MessageStyle:   lipgloss.NewStyle().Foreground(FadedGrey),
	},
}

type AppStyles struct {
	LogEntryStyles       LogMessageStyles
	NearbyLogEntryStyles NearbyLogEntryStyles
	CursorStyle          lipgloss.Style
}

type LogMessageStyles struct {
	LineNumberStyle lipgloss.Style
	OriginStyle     lipgloss.Style
	TimestampStyle  lipgloss.Style
	LevelStyle      lipgloss.Style
	MessageStyle    lipgloss.Style
}

type NearbyLogEntryStyles struct {
	OriginStyle    lipgloss.Style
	TimestampStyle lipgloss.Style
	LevelStyle     lipgloss.Style
	MessageStyle   lipgloss.Style
}
