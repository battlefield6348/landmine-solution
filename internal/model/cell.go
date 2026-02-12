package model

import "fmt"

// CellState 代表儲存格的狀態
type CellState int

const (
	StateUnknown CellState = -1 // 未解 (u)
	StateFlag    CellState = -2 // 旗幟 (f)
	// 0-8 代表數字
)

// Cell 儲存格結構
type Cell struct {
	State       CellState
	Probability float64 // 地雷機率 (0.0 - 1.0)
}

func (s CellState) String() string {
	switch s {
	case StateUnknown:
		return "u"
	case StateFlag:
		return "f"
	default:
		return fmt.Sprintf("%d", int(s))
	}
}
