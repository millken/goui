package goui

type Layout struct {
	cols, rows uint32
	cells      []*Cell
}

func NewLayout(cols, rows uint32) *Layout {
	if cols == 0 {
		cols = 1
	}
	if rows == 0 {
		rows = 1
	}
	return &Layout{
		cols: cols,
		rows: rows,
	}
}

func (l *Layout) SetLabel(label *Label) {
}
