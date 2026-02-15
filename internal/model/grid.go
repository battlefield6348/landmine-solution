package model

// Grid 代表盤面
type Grid struct {
	Rows  int
	Cols  int
	Cells [][]Cell
}

// NewGrid 建立一個指定尺寸且預設為未解的盤面
func NewGrid(rows, cols int) *Grid {
	g := &Grid{
		Rows:  rows,
		Cols:  cols,
		Cells: make([][]Cell, rows),
	}
	for r := 0; r < rows; r++ {
		g.Cells[r] = make([]Cell, cols)
		for c := 0; c < cols; c++ {
			g.Cells[r][c] = Cell{State: StateUnknown}
		}
	}
	return g
}

// AddRow 在頂部 (top=true) 或底部 (top=false) 新增一列
func (g *Grid) AddRow(top bool) {
	newCells := make([][]Cell, g.Rows+1)
	if top {
		newCells[0] = make([]Cell, g.Cols)
		for c := 0; c < g.Cols; c++ {
			newCells[0][c] = Cell{State: StateUnknown}
		}
		copy(newCells[1:], g.Cells)
	} else {
		copy(newCells, g.Cells)
		newCells[g.Rows] = make([]Cell, g.Cols)
		for c := 0; c < g.Cols; c++ {
			newCells[g.Rows][c] = Cell{State: StateUnknown}
		}
	}
	g.Cells = newCells
	g.Rows++
}

// RemoveRow 從頂部 (top=true) 或底部 (top=false) 移除一列
func (g *Grid) RemoveRow(top bool) {
	if g.Rows <= 1 {
		return
	}
	newCells := make([][]Cell, g.Rows-1)
	if top {
		copy(newCells, g.Cells[1:])
	} else {
		copy(newCells, g.Cells[:g.Rows-1])
	}
	g.Cells = newCells
	g.Rows--
}

// AddCol 在左側 (left=true) 或右側 (left=false) 新增一欄
func (g *Grid) AddCol(left bool) {
	for r := 0; r < g.Rows; r++ {
		newRow := make([]Cell, g.Cols+1)
		if left {
			newRow[0] = Cell{State: StateUnknown}
			copy(newRow[1:], g.Cells[r])
		} else {
			copy(newRow, g.Cells[r])
			newRow[g.Cols] = Cell{State: StateUnknown}
		}
		g.Cells[r] = newRow
	}
	g.Cols++
}

// RemoveCol 從左側 (left=true) 或右側 (left=false) 移除一欄
func (g *Grid) RemoveCol(left bool) {
	if g.Cols <= 1 {
		return
	}
	for r := 0; r < g.Rows; r++ {
		newRow := make([]Cell, g.Cols-1)
		if left {
			copy(newRow, g.Cells[r][1:])
		} else {
			copy(newRow, g.Cells[r][:g.Cols-1])
		}
		g.Cells[r] = newRow
	}
	g.Cols--
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
			if nr >= 0 && nr < g.Rows && nc >= 0 && nc < g.Cols {
				neighbors = append(neighbors, [2]int{nr, nc})
			}
		}
	}
	return neighbors
}
