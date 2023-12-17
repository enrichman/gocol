package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	WHITE = "#FFF"
	BLACK = "#000"
	GRAY  = "#AAA"

	RED         = "#E03C32"
	ORANGE      = "#FF8C01"
	YELLOW      = "#FFE733"
	LIGHT_GREEN = "#7BB662"
	GREEN       = "#006B3D"
	DARK_GREEN  = "#024E1B"

	BRIGHT_GREEN = "#00D700"
	FUCHSIA      = "#FF00FF"
)

var (
	brightGreenStyle = NewStyleWithFG(BRIGHT_GREEN)
	yellowStyle      = NewStyleWithFG(YELLOW)
	redColor         = NewStyleWithFG(RED)
	fuchsiaColor     = NewStyleWithFG(FUCHSIA)

	covNopeStyle = NewStyleWithFG(BLACK).
			Background(lipgloss.Color(GRAY))

	covNope   = NewCoverageColor(0, 0, BLACK, GRAY)
	covColors = []*CoverageColor{
		NewCoverageColor(0, 30, BLACK, RED),
		NewCoverageColor(30, 50, BLACK, ORANGE),
		NewCoverageColor(50, 70, BLACK, YELLOW),
		NewCoverageColor(70, 80, BLACK, LIGHT_GREEN),
		NewCoverageColor(80, 90, WHITE, GREEN),
		NewCoverageColor(90, 100, WHITE, DARK_GREEN),
	}
)

type CoverageColor struct {
	start float64
	end   float64
	style lipgloss.Style
}

func (c *CoverageColor) Color(v string) string {
	return c.style.Render(v)
}

func NewCoverageColor(start, end float64, fgColor, bgColor string) *CoverageColor {
	return &CoverageColor{
		start: start,
		end:   end,
		style: NewStyleWithFG(fgColor).
			Background(lipgloss.Color(bgColor)),
	}
}

func NewStyleWithFG(fgColor string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(fgColor))
}

func main() {
	checkPipe()

	// read lines and fail if something happened
	lines, err := readLines()
	if err != nil {
		log.Fatal(err)
	}

	noTestFileLines := []string{}
	colorizedLines := []string{}

	for _, line := range lines {
		// append to write [no test files] at the end
		if strings.Contains(line, "[no test files]") {
			line = reorderNoTestLine(line)
			line := colorizeLine(line)
			noTestFileLines = append(noTestFileLines, line)
			continue
		}

		// reorder columns for coverage lines
		if strings.Contains(line, "coverage:") {
			line = reorderCoverageLine(line)
		}

		line := colorizeLine(line)
		colorizedLines = append(colorizedLines, line)
	}

	// print colorized lines
	for _, line := range colorizedLines {
		fmt.Println(line)
	}

	// print [no test files] lines
	for _, line := range noTestFileLines {
		fmt.Println(line)
	}
}

func checkPipe() {
	// check if there is something to read on STDIN
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	isPipe := (stat.Mode() & os.ModeCharDevice) == 0
	if !isPipe {
		fmt.Print("Gocol colorize your coverage\n\n")
		fmt.Println("Usage:\n\tgo test -cover ./... | gocol")
		return
	}
}

func readLines() ([]string, error) {
	var lines []string
	var line string
	var err error

	// read until EOF
	reader := bufio.NewReader(os.Stdin)
	for err == nil {
		line, err = reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	if errors.Is(err, io.EOF) {
		return lines, nil
	}
	return nil, err
}

func reorderCoverageLine(line string) string {
	splitted := strings.Split(line, "\t")

	if len(splitted) > 3 {
		covColumn := splitted[3]

		// keep the [no tests to run] at the end
		if strings.HasSuffix(covColumn, " [no tests to run]") {
			covColumn = strings.TrimSuffix(covColumn, " [no tests to run]")
			splitted = append(splitted, "[no tests to run]")
		}

		// switch package and coverage
		splitted[3], splitted[1] = splitted[1], covColumn
	}

	return strings.Join(splitted, "\t")
}

func reorderNoTestLine(line string) string {
	splitted := strings.Split(line, "\t")

	if len(splitted) == 3 {
		covColumn := splitted[2]
		// switch [no test files] and package
		splitted[2], splitted[1] = splitted[1], covColumn
	}

	return strings.Join(splitted, "\t")
}

func colorizeLine(line string) string {
	if strings.Contains(line, "coverage:") {
		return colorizeCoverageLine(line)
	}

	if strings.HasPrefix(line, "?") {
		line = colorizeMatch(line, "?", fuchsiaColor)
		line = colorizeMatchBg(line, "[no test files]", covNopeStyle)
		return line
	}

	trimmed := line
	if strings.HasPrefix(line, "--- ") {
		trimmed = strings.TrimPrefix(line, "--- ")
	}

	style := lipgloss.NewStyle()

	switch {
	case strings.HasPrefix(trimmed, "PASS"):
		style = brightGreenStyle
	case strings.HasPrefix(trimmed, "SKIP"):
		style = yellowStyle
	case strings.HasPrefix(trimmed, "FAIL"):
		style = redColor
	}

	return style.Render(line)
}

func colorizeCoverageLine(line string) string {
	line = colorizeMatch(line, "ok", brightGreenStyle)
	line = colorizeMatch(line, "(cached)", yellowStyle)
	line = colorizeMatch(line, "[no tests to run]", fuchsiaColor)
	line = colorizeCoverage(line)

	return line
}

func colorizeMatch(line, match string, style lipgloss.Style) string {
	index := strings.Index(line, match)
	if index > -1 {
		return strings.Replace(line, match, style.Render(match), 1)
	}
	return line
}

func colorizeMatchBg(line, match string, style lipgloss.Style) string {
	index := strings.Index(line, match)
	if index > -1 {
		return strings.Replace(line, match, style.Render(match), 1)
	}
	return line
}

func colorizeCoverage(line string) string {
	coverageIndex := strings.Index(line, "coverage:")
	if coverageIndex == -1 {
		return line
	}

	// fix different percentage lenght
	// if strings.Contains(line, "coverage: 0.0%") {
	// 	line = strings.Replace(line, "coverage: 0.0%", "coverage:   0.0%", 1)
	// } else if !strings.Contains(line, "coverage: 100.0%") {
	// 	line = strings.Replace(line, "coverage: ", "coverage:  ", 1)
	// }

	statementsIndex := strings.Index(line, "statements")
	if statementsIndex == -1 {
		return line
	}
	end := statementsIndex + len("statements")

	// fix different percentage lenght
	switch end - coverageIndex {
	case 28:
		line = strings.Replace(line, "coverage: ", "coverage:   ", 1)
		end += 2
	case 29:
		line = strings.Replace(line, "coverage: ", "coverage:  ", 1)
		end += 1
	}

	percentage, err := findPercentageValue(line)
	if err != nil {
		fmt.Println(err)
		return line
	}

	covColor := getCoverageColor(percentage)

	return line[:coverageIndex] +
		covColor.Color(line[coverageIndex:end]) +
		line[end:]
}

func findPercentageValue(line string) (float64, error) {
	percentageStr := ""

	for _, field := range strings.Fields(line) {
		if strings.HasSuffix(field, "%") {
			percentageStr = strings.TrimSuffix(field, "%")
			break
		}
	}

	return strconv.ParseFloat(percentageStr, 32)
}

func getCoverageColor(coverage float64) *CoverageColor {
	if coverage == 0 {
		return covNope
	}

	for _, covCol := range covColors {
		if coverage >= covCol.start &&
			(coverage < covCol.end || covCol.end == 100) {
			return covCol

		}
	}

	return covNope
}
