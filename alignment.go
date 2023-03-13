package tabwriter

import "strings"

type alignment uint8

const (
	RightAlign = iota
	LeftAlign
	CenterAlign
)

type alignmentFunc func(cell string, width uint) []string

type alignmentOperator struct {
	columnSeperator string
}

type aligner interface {
	center(cell string, width uint) []string
	right(cell string, width uint) []string
	left(cell string, width uint) []string
	full(cell string) string
	alignFuncMap() map[alignment]alignmentFunc
}

func newAligner(colSep string) aligner {
	return &alignmentOperator{
		columnSeperator: colSep,
	}
}

func (a *alignmentOperator) center(cell string, width uint) []string {
	var (
		result            []string
		left, right, diff int
	)

	for _, part := range splitBySize(cell, width) {
		diff = int(width) - len(part)

		if diff%2 == 0 {
			left, right = diff/2, diff/2
		} else {
			left, right = diff/2, diff/2+1
		}

		result = append(result, a.columnSeperator+
			strings.Repeat(whitespace, left)+
			part+
			strings.Repeat(whitespace, right)+
			a.columnSeperator)
	}

	return result
}

func (a *alignmentOperator) right(cell string, width uint) []string {
	var (
		result []string
		diff   int
	)

	for _, part := range splitBySize(cell, width) {
		diff = int(width) - len(part)

		result = append(result, a.columnSeperator+
			strings.Repeat(whitespace, diff)+
			part+
			a.columnSeperator)
	}

	return result
}

func (a *alignmentOperator) left(cell string, width uint) []string {
	var (
		result []string
		diff   int
	)

	for _, part := range splitBySize(cell, width) {
		diff = int(width) - len(part)

		result = append(result, a.columnSeperator+
			part+
			strings.Repeat(whitespace, diff)+
			a.columnSeperator)
	}

	return result
}

func (a *alignmentOperator) full(cell string) string {
	return a.columnSeperator +
		cell +
		a.columnSeperator
}

func (a *alignmentOperator) alignFuncMap() map[alignment]alignmentFunc {
	return map[alignment]alignmentFunc{
		RightAlign:  a.right,
		LeftAlign:   a.left,
		CenterAlign: a.center,
	}
}
