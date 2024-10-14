package unicode

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

func BracketTopWithText(bracketStyle lipgloss.Style, textStyle lipgloss.Style, format string, args ...any) string {
	bracketCorner := bracketStyle.Render(BracketTop)
	leftBracket := bracketStyle.Render("[")
	rightBracket := bracketStyle.Render("]")
	text := textStyle.Render(fmt.Sprintf(format, args...))
	return fmt.Sprintf("%s %s %s %s\n", bracketCorner, leftBracket, text, rightBracket)
}

func PrefixVerticalLine(lineStyle lipgloss.Style, format string, args ...any) string {
	return fmt.Sprintf("%s %s", lineStyle.Render(VerticalLine), fmt.Sprintf(format, args...))
}

func BracketBottomWithText(bracketStyle lipgloss.Style, textStyle lipgloss.Style, format string, args ...any) string {
	bracketCorner := bracketStyle.Render(BracketBottom)
	leftBracket := bracketStyle.Render("[")
	rightBracket := bracketStyle.Render("]")
	text := textStyle.Render(fmt.Sprintf(format, args...))
	return fmt.Sprintf("%s %s %s %s\n", bracketCorner, leftBracket, text, rightBracket)
}
