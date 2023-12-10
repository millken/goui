package ui

type Layout struct {
	row, col int
}

func NewLayout(row, col int) *Layout {
	return &Layout{
		row: row,
		col: col,
	}
}
