package goui

type CellType int

const (
	CellTypeEmpty CellType = iota
	CellTypeComponent
	CellTypeLayout
)

type cellDimension struct {
	ForcedSize    float32
	PaddingBefore float32
	PaddingAfter  float32
	Size          float32
}

type Cell struct {
	// The type of the cell.
	Type    CellType
	Visible bool
	Enabled bool
	Parrent *Layout
}
