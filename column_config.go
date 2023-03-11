package tabwriter

type columnConfig struct {
	columnName    string
	columnAlign   alignment
	columnWidth   uint
	enableBestFit bool
}

type ColumnConfig interface {
	name() string
	alignment() alignment
	width() uint
	setWidth(width uint)
	isBestFitEnabled() bool
}

func NewColumnConfig(name string, align alignment, width uint) ColumnConfig {
	bestFit := false
	if width == 0 {
		bestFit = true
	}

	return &columnConfig{
		columnName:    name,
		columnAlign:   align,
		columnWidth:   width,
		enableBestFit: bestFit,
	}
}

func (c *columnConfig) name() string {
	return c.columnName
}

func (c *columnConfig) alignment() alignment {
	return c.columnAlign
}

func (c *columnConfig) width() uint {
	return c.columnWidth
}

func (c *columnConfig) isBestFitEnabled() bool {
	return c.enableBestFit
}

func (c *columnConfig) setWidth(width uint) {
	c.columnWidth = width
}
