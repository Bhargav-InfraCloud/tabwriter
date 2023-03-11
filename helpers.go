package tabwriter

import "strings"

func (m *manager) bestFit(index int) uint {
	max := len(m.columnConfigs[index].name())
	for _, row := range m.rows {
		rowSize := len(row[index])
		if rowSize > max {
			max = rowSize
		}
	}

	return uint(max)
}

func (m *manager) createEmptySheet(linesCount int) [][]string {
	var (
		lines [][]string
		cell  string
	)

	for r := 0; r < linesCount; r++ {
		var line []string

		for c := 0; c < len(m.columnConfigs); c++ {
			cell = strings.Repeat(whitespace, int(m.columnConfigs[c].width()))
			cell = m.aligner.full(cell)
			line = append(line, cell)
		}

		lines = append(lines, line)
	}

	return lines
}

func splitBySize(str string, size uint) []string {
	var subStr []string

	for len(str) > int(size) {
		subStr = append(subStr, str[:size])
		str = str[size:]
	}

	subStr = append(subStr, str)

	return subStr
}

func mergeValues(emptySheet [][]string, index int, columnVals []string) [][]string {
	for r, cell := range columnVals {
		emptySheet[r][index] = cell
	}

	return emptySheet
}
