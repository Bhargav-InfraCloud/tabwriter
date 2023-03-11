package tabwriter

import "fmt"

const (
	whitespace             = ` `
	defaultColumnSeperator = whitespace
	defaultHeaderSeperator = `-`
)

var ErrUnexpectedExtraColumns error = fmt.Errorf("unexpected extra columns")
