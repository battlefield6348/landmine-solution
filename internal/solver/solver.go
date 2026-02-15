package solver

import (
	"landmine-solution/internal/model"
)

// Result 包含解題後的機率資訊
type Result struct {
	Probabilities [][]float64
	Solvable      bool
	Timeout       bool
}

// Solve 計算盤面上每個未知格子的地雷機率
func Solve(g *model.Grid) *Result {
	res := &Result{
		Probabilities: make([][]float64, g.Rows),
		Solvable:      true,
	}
	for i := 0; i < g.Rows; i++ {
		res.Probabilities[i] = make([]float64, g.Cols)
	}

	// 1. 找出所有未解格子
	var unknownPos [][2]int
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {
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
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {
			state := g.Cells[r][c].State
			if state < 0 {
				continue
			}

			// 找出此數字格周圍的未知格索引與已確定的旗幟數
			cst := constraint{
				target:           int(state),
				uIndices:         []int{},
				knownMines:       0,
				missingNeighbors: 0,
			}

			// 計算實際鄰居與遺失的鄰居（在盤面外的）
			for dr := -1; dr <= 1; dr++ {
				for dc := -1; dc <= 1; dc++ {
					if dr == 0 && dc == 0 {
						continue
					}
					nr, nc := r+dr, c+dc
					if nr >= 0 && nr < g.Rows && nc >= 0 && nc < g.Cols {
						// 盤面內的鄰居
						if g.Cells[nr][nc].State == model.StateFlag {
							cst.knownMines++
						} else if g.Cells[nr][nc].State == model.StateUnknown {
							for i, upos := range unknownPos {
								if upos[0] == nr && upos[1] == nc {
									cst.uIndices = append(cst.uIndices, i)
									break
								}
							}
						}
					} else {
						// 盤面外的鄰居
						cst.missingNeighbors++
					}
				}
			}
			constraints = append(constraints, cst)
		}
	}

	// 3. 找出哪些未知格是有約束的
	isConstrained := make([]bool, len(unknownPos))
	constrainedCount := 0
	for _, cst := range constraints {
		for _, uIdx := range cst.uIndices {
			if !isConstrained[uIdx] {
				isConstrained[uIdx] = true
				constrainedCount++
			}
		}
	}

	// 建立受約束未知格的索引映射
	constrainedToOriginal := make([]int, 0, constrainedCount)
	originalToConstrained := make(map[int]int)
	for i, constrained := range isConstrained {
		if constrained {
			originalToConstrained[i] = len(constrainedToOriginal)
			constrainedToOriginal = append(constrainedToOriginal, i)
		}
	}

	// 更新約束中的索引
	for i := range constraints {
		newIndices := make([]int, len(constraints[i].uIndices))
		for j, oldIdx := range constraints[i].uIndices {
			newIndices[j] = originalToConstrained[oldIdx]
		}
		constraints[i].uIndices = newIndices
	}

	// 4. 預處理：建立受約束未知格索引到受影響約束的映射
	uIdxToConstraints := make([][]int, constrainedCount)
	for i, cst := range constraints {
		for _, uIdx := range cst.uIndices {
			uIdxToConstraints[uIdx] = append(uIdxToConstraints[uIdx], i)
		}
	}

	totalValidConfigs := 0
	mineCounts := make([]int, constrainedCount)
	currentMines := make([]bool, constrainedCount)
	iterations := 0
	const maxIterations = 10000000 // 增加安全限制到 1000 萬次

	// 5. 回溯法 (僅針對受約束的格子)
	var backtrack func(uIdx int)
	backtrack = func(uIdx int) {
		iterations++
		if iterations > maxIterations {
			return
		}

		// 剪枝檢查：只檢查受上一格決策影響的約束
		if uIdx > 0 {
			lastUIdx := uIdx - 1
			for _, cIdx := range uIdxToConstraints[lastUIdx] {
				if !constraints[cIdx].isPossible(currentMines, uIdx) {
					return
				}
			}
		}

		if uIdx == constrainedCount {
			totalValidConfigs++
			for i := 0; i < constrainedCount; i++ {
				if currentMines[i] {
					mineCounts[i]++
				}
			}
			return
		}

		// 分支
		currentMines[uIdx] = false
		backtrack(uIdx + 1)

		currentMines[uIdx] = true
		backtrack(uIdx + 1)
	}

	backtrack(0)

	if iterations > maxIterations {
		res.Solvable = false // 標記為無法在時限內解出
		res.Timeout = true
	}

	// 6. 設定機率
	if totalValidConfigs > 0 {
		// 受約束的格子
		for i, originalIdx := range constrainedToOriginal {
			pos := unknownPos[originalIdx]
			res.Probabilities[pos[0]][pos[1]] = float64(mineCounts[i]) / float64(totalValidConfigs)
		}
		// 未受約束的格子，預設機率為 50%
		for i, constrained := range isConstrained {
			if !constrained {
				pos := unknownPos[i]
				res.Probabilities[pos[0]][pos[1]] = 0.5
			}
		}
	} else if len(unknownPos) > 0 {
		// 如果完全沒有合法組合且有未知格，視為矛盾 (或因超限未算出)
		res.Solvable = false
	}

	return res
}

type constraint struct {
	target           int   // 目標雷數
	uIndices         []int // 周圍未知格在 unknownPos 中的索引
	knownMines       int   // 周圍已確定的旗幟數
	missingNeighbors int   // 在盤面外的鄰居數量 (局部地圖用)
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

	// 2. 即使剩下的全放雷(內部+外部)也湊不到目標
	// 修正：必須考慮盤面外的鄰居(missingNeighbors)
	if assignedMines+unassignedCount+c.missingNeighbors < c.target {
		return false
	}

	// 如果內部所有未知格都已決定 (unassignedCount == 0)
	if unassignedCount == 0 {
		// 在內部全決定的情況下，外部必須能補足剩下的雷
		// 規則：(目標 - 內部已確定的雷) 必須小於等於 外部空間
		remainingNeeded := c.target - assignedMines
		if remainingNeeded < 0 || remainingNeeded > c.missingNeighbors {
			return false
		}
	}

	return true
}
