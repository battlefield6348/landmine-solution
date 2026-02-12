package model

// Grid 代表正方形盤面
type Grid struct {
	Size  int
	Cells [][]Cell
}

// NewGrid 建立一個指定尺寸且預設為未解的盤面
func NewGrid(size int) *Grid {
	g := &Grid{
		Size:  size,
		Cells: make([][]Cell, size),
	}
	for r := 0; r < size; r++ {
		g.Cells[r] = make([]Cell, size)
		for c := 0; c < size; c++ {
			g.Cells[r][c] = Cell{State: StateUnknown}
		}
	}
	return g
}

// Resize 調整盤面尺寸並儘可能保留現有資料
func (g *Grid) Resize(newSize int) {
	newCells := make([][]Cell, newSize)
	for r := 0; r < newSize; r++ {
		newCells[r] = make([]Cell, newSize)
		for c := 0; c < newSize; c++ {
			if r < g.Size && c < g.Size {
				newCells[r][c] = g.Cells[r][c]
			} else {
				newCells[r][c] = Cell{State: StateUnknown}
			}
		}
	}
	g.Cells = newCells
	g.Size = newSize
}

// GetNeighbors 取得坐標 (r, c) 周圍的鄰居座標
func (g *Grid) GetNeighbors(r, c int) [][2]int {
	neighbors := [][2]int{}
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			nr, nc := r+dr, c+dc
			if nr >= 0 && nr < g.Size && nc >= 0 && nc < g.Size {
				neighbors = append(neighbors, [2]int{nr, nc})
			}
		}
	}
	return neighbors
}
