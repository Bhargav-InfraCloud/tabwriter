package tabwriter

import (
	"fmt"
	"io"
	"strings"
)

type manager struct {
	columnConfigs     []ColumnConfig
	rows              [][]string
	headerSeperator   string
	isHeaderSeperated bool
	formattedSheet    [][]string
	aligner           aligner
}

type TableCreator interface {
	AddRows(rows ...[]string) TableCreator
	AddRow(row []string) TableCreator
	ColumnSeperator(colSep string) TableCreator
	HeaderSeperator(headerSep string) TableCreator
	EnableHeaderSeperator() TableCreator
	Process() error
	Print(w io.Writer)
}

func NewTableCreator(configs ...ColumnConfig) TableCreator {
	return &manager{
		columnConfigs:     configs,
		rows:              [][]string{},
		headerSeperator:   defaultHeaderSeperator,
		isHeaderSeperated: false,
		aligner:           newAligner(defaultColumnSeperator),
	}
}

func (m *manager) AddRows(rows ...[]string) TableCreator {
	m.rows = append(m.rows, rows...)

	return m
}

func (m *manager) AddRow(row []string) TableCreator {
	m.rows = append(m.rows, row)

	return m
}

func (m *manager) ColumnSeperator(colSep string) TableCreator {
	m.aligner = newAligner(colSep)

	return m
}

func (m *manager) HeaderSeperator(headerSep string) TableCreator {
	m.headerSeperator = headerSep
	m.isHeaderSeperated = true

	return m
}

func (m *manager) EnableHeaderSeperator() TableCreator {
	m.isHeaderSeperated = true

	return m
}

func (m *manager) Process() error {

	m.addHeaders()

	if m.isHeaderSeperated {
		m.addHeaderSeperator()
	}

	err := m.addRows()
	if err != nil {
		return err
	}

	return nil
}

func (m *manager) Print(w io.Writer) {
	for _, row := range m.formattedSheet {
		for _, cell := range row {
			fmt.Fprint(w, cell)
		}

		fmt.Fprintln(w)
	}
}

func (m *manager) addHeaders() {
	var (
		alignFunc           alignmentFunc
		config              ColumnConfig
		index               int
		paddedParts, header []string
		maxExtraCount       int

		alignmentFuncMap = m.aligner.alignFuncMap()
		extraLines       = map[int][]string{}
	)

	for index, config = range m.columnConfigs {
		if config.isBestFitEnabled() {
			config.setWidth(m.bestFit(index))
		}

		alignFunc = alignmentFuncMap[config.alignment()]
		paddedParts = alignFunc(config.name(), config.width())
		if len(paddedParts) > 1 {
			extraLines[index] = paddedParts[1:]
			if len(paddedParts[1:]) > maxExtraCount {
				maxExtraCount = len(paddedParts[1:])
			}
		}

		header = append(header, paddedParts[0])

	}

	m.formattedSheet = append(m.formattedSheet, header)

	if len(extraLines) > 0 {
		emptySheet := m.createEmptySheet(maxExtraCount)
		for colIndex, extras := range extraLines {
			emptySheet = mergeValues(emptySheet, colIndex, extras)
		}

		m.formattedSheet = append(m.formattedSheet, emptySheet...)
	}
}

func (m *manager) addHeaderSeperator() {
	var (
		seperator  string
		config     ColumnConfig
		seperators []string
	)

	for _, config = range m.columnConfigs {
		seperator = strings.Repeat(m.headerSeperator, int(config.width()))
		seperator = m.aligner.full(seperator)
		seperators = append(seperators, seperator)
	}

	m.formattedSheet = append(m.formattedSheet, seperators)
}

func (m *manager) addRows() error {
	var (
		alignFunc        alignmentFunc
		cell             string
		config           ColumnConfig
		index            int
		row, paddedParts []string
		maxExtraCount    int

		alignmentFuncMap = m.aligner.alignFuncMap()
		extraLines       = map[int][]string{}
	)

	for _, row = range m.rows {
		if len(row) != len(m.columnConfigs) {
			return ErrUnexpectedExtraColumns
		}

		var record []string
		for index, cell = range row {
			config = m.columnConfigs[index]
			alignFunc = alignmentFuncMap[config.alignment()]
			paddedParts = alignFunc(cell, config.width())

			if len(paddedParts) > 1 {
				extraLines[index] = paddedParts[1:]
				if len(paddedParts[1:]) > maxExtraCount {
					maxExtraCount = len(paddedParts[1:])
				}
			}

			record = append(record, paddedParts[0])

		}

		m.formattedSheet = append(m.formattedSheet, record)

		if len(extraLines) > 0 {
			emptySheet := m.createEmptySheet(maxExtraCount)
			for colIndex, extras := range extraLines {
				emptySheet = mergeValues(emptySheet, colIndex, extras)
			}

			m.formattedSheet = append(m.formattedSheet, emptySheet...)
		}
	}

	return nil
}
