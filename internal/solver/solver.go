package solver

import (
	"landmine-solution/internal/model"
)

// Result 包含解題後的機率資訊
type Result struct {
	Probabilities [][]float64
}

// Solve 計算盤面上每個未知格子的地雷機率
func Solve(g *model.Grid) *Result {
	res := &Result{
		Probabilities: make([][]float64, g.Size),
	}
	for i := 0; i < g.Size; i++ {
		res.Probabilities[i] = make([]float64, g.Size)
	}

	// 1. 找出所有未解格子
	var unknownPos [][2]int
	for r := 0; r < g.Size; r++ {
		for c := 0; c < g.Size; c++ {
			if g.Cells[r][c].State == model.StateUnknown {
				unknownPos = append(unknownPos, [2]int{r, c})
			}
		}
	}

	if len(unknownPos) == 0 {
		return res
	}

	// 2. 預處理：找出所有會被未知格影響的「數字約束格」
	constraints := []constraint{}
	for r := 0; r < g.Size; r++ {
		for c := 0; c < g.Size; c++ {
			state := g.Cells[r][c].State
			if state < 0 {
				continue
			}

			// 找出此數字格周圍的未知格索引與已確定的旗幟數
			cst := constraint{
				target: int(state),
				uIndices: []int{},
				knownMines: 0,
			}
			neighbors := g.GetNeighbors(r, c)
			for _, n := range neighbors {
				if g.Cells[n[0]][n[1]].State == model.StateFlag {
					cst.knownMines++
				} else if g.Cells[n[0]][n[1]].State == model.StateUnknown {
					// 找對應的 unknownPos 索引
					for i, upos := range unknownPos {
						if upos[0] == n[0] && upos[1] == n[1] {
							cst.uIndices = append(cst.uIndices, i)
							break
						}
					}
				}
			}
			constraints = append(constraints, cst)
		}
	}

	totalValidConfigs := 0
	mineCounts := make([]int, len(unknownPos))
	currentMines := make([]bool, len(unknownPos))

	// 3. 回溯法 (只針對受影響的約束進行快速檢查)
	var backtrack func(uIdx int)
	backtrack = func(uIdx int) {
		// 剪枝檢查：只檢查受當前決策影響的約束
		if uIdx > 0 {
			// lastUIdx := uIdx - 1 // This variable is not used.
			for i := range constraints {
				if !constraints[i].isPossible(currentMines, uIdx) {
					return
				}
			}
		}

		if uIdx == len(unknownPos) {
			// 到達末端，所有未知格都已決定，進行最終驗證
			// 這裡的 isPossible 已經包含了最終驗證邏輯，
			// 只要所有約束都通過 isPossible(..., len(unknownPos)) 檢查，就代表配置有效
			// 因為在 uIdx == len(unknownPos) 時，unassignedCount 會是 0，
			// isPossible 會檢查 assignedMines == c.target
			totalValidConfigs++
			for i := range unknownPos {
				if currentMines[i] {
					mineCounts[i]++
				}
			}
			return
		}

		// 分支：放雷 or 不放雷
		currentMines[uIdx] = false
		backtrack(uIdx + 1)

		currentMines[uIdx] = true
		backtrack(uIdx + 1)
	}

	backtrack(0)

	if totalValidConfigs > 0 {
		for i, pos := range unknownPos {
			res.Probabilities[pos[0]][pos[1]] = float64(mineCounts[i]) / float64(totalValidConfigs)
		}
	}

	return res
}

type constraint struct {
	target     int   // 目標雷數
	uIndices   []int // 周圍未知格在 unknownPos 中的索引
	knownMines int   // 周圍已確定的旗幟數
}

// isPossible 檢查在目前的配置下，此約束是否還有可能成立
func (c *constraint) isPossible(mines []bool, currentUIdx int) bool {
	assignedMines := c.knownMines
	unassignedCount := 0

	for _, uIdx := range c.uIndices {
		if uIdx < currentUIdx { // 已經決定的未知格
			if mines[uIdx] {
				assignedMines++
			}
		} else { // 尚未決定的未知格
			unassignedCount++
		}
	}

	// 1. 已經確定的雷超過目標
	if assignedMines > c.target {
		return false
	}
	// 2. 即使剩下的全放雷也湊不到目標
	if assignedMines+unassignedCount < c.target {
		return false
	}
	
	// 如果所有未知格都已決定 (unassignedCount == 0)，則必須精確匹配
	if unassignedCount == 0 && assignedMines != c.target {
		return false
	}

	return true
}
