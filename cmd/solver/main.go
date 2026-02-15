package main

import (
	"bufio"
	"fmt"
	"landmine-solution/internal/model"
	"landmine-solution/internal/solver"
	"os"
	"strconv"
	"strings"
)

func main() {
	grid := model.NewGrid(3, 3)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		displayResult(grid)
		fmt.Println("\n[指令說明]")
		fmt.Println("- 's c r v' : 修改座標 (欄 c, 列 r) 值為 v (例: s 1 0 2)")
		fmt.Println("- 直接輸入全部符號 : 批次更新 (例: 1 u u 1 f ...)")
		fmt.Println("- 擴展: 'at' (上), 'ab' (下), 'al' (左), 'ar' (右)")
		fmt.Println("- 縮減: 'rt' (上), 'rb' (下), 'rl' (左), 'rr' (右)")
		fmt.Println("- 指令: 'empty' (清空), 'exit' (結束)")
		fmt.Println("- 符號: 0-8, e (空白/0), u (未解), f (旗幟)")

		fmt.Printf("\n[尺寸 %dx%d] 請輸入指令或內容: ", grid.Cols, grid.Rows)

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "exit" {
			break
		}

		// 擴展指令
		switch input {
		case "at":
			grid.AddRow(true)
			fmt.Println("✅ 頂部已增加一列")
			continue
		case "ab":
			grid.AddRow(false)
			fmt.Println("✅ 底部已增加一列")
			continue
		case "al":
			grid.AddCol(true)
			fmt.Println("✅ 左側已增加一欄")
			continue
		case "ar":
			grid.AddCol(false)
			fmt.Println("✅ 右側已增加一欄")
			continue
		case "rt":
			grid.RemoveRow(true)
			fmt.Println("✅ 頂部已移除一列")
			continue
		case "rb":
			grid.RemoveRow(false)
			fmt.Println("✅ 底部已移除一列")
			continue
		case "rl":
			grid.RemoveCol(true)
			fmt.Println("✅ 左側已移除一欄")
			continue
		case "rr":
			grid.RemoveCol(false)
			fmt.Println("✅ 右側已移除一欄")
			continue
		}

		if input == "empty" {
			grid = model.NewGrid(grid.Rows, grid.Cols)
			fmt.Printf("✅ 盤面 (%dx%d) 已清空\n", grid.Cols, grid.Rows)
			continue
		}

		tokens := strings.Fields(input)
		if len(tokens) == 0 {
			continue
		}

		// 1. 處理明確的座標修改: s c r v (如 s 1 0 2)
		if strings.ToLower(tokens[0]) == "s" {
			if len(tokens) == 4 {
				c, errC := strconv.Atoi(tokens[1])
				r, errR := strconv.Atoi(tokens[2])
				if errR == nil && errC == nil && c >= 0 && c < grid.Cols && r >= 0 && r < grid.Rows {
					if updateCell(grid, r, c, tokens[3]) {
						fmt.Printf("✅ 已將座標 (欄%d, 列%d) 更新為 %s\n", c, r, tokens[3])
						continue
					}
				}
			}
			fmt.Println("❌ 's' 指令格式錯誤。範例: s 1 0 2 (將 欄1, 列0 改為 2)")
			continue
		}

		// 2. 處理自動識別的座標修改: c r v (3 tokens)
		if len(tokens) == 3 {
			c, errC := strconv.Atoi(tokens[0])
			r, errR := strconv.Atoi(tokens[1])
			if errR == nil && errC == nil && c >= 0 && c < grid.Cols && r >= 0 && r < grid.Rows {
				if updateCell(grid, r, c, tokens[2]) {
					fmt.Printf("✅ 已將座標 (欄%d, 列%d) 更新為 %s\n", c, r, tokens[2])
					continue
				}
			}
		}

		// 3. 處理批次更新 (tokens 數量必須等於 Rows * Cols)
		expectedCount := grid.Rows * grid.Cols
		if len(tokens) == expectedCount {
			valid := true
			for i, t := range tokens {
				r, c := i/grid.Cols, i%grid.Cols
				if !updateCell(grid, r, c, t) {
					fmt.Printf("❌ 符號 '%s' 無效 (索引 %d)\n", t, i)
					valid = false
					break
				}
			}
			if valid {
				fmt.Println("✅ 盤面已整批更新完成")
				continue
			}
		} else {
			fmt.Printf("❌ 輸入格式不符：目前 %dx%d 盤面需要 %d 個符號，但你輸入了 %d 個\n", grid.Cols, grid.Rows, expectedCount, len(tokens))
			fmt.Println("   提示：單格修改請用 's c r v' (例: s 1 0 2)")
			continue
		}
	}
}

// updateCell 輔助函式，更新特定格子的狀態
func updateCell(g *model.Grid, r, c int, val string) bool {
	switch strings.ToLower(val) {
	case "u":
		g.Cells[r][c].State = model.StateUnknown
	case "f":
		g.Cells[r][c].State = model.StateFlag
	case "e":
		g.Cells[r][c].State = model.CellState(0)
	default:
		num, err := strconv.Atoi(val)
		if err != nil || num < 0 || num > 8 {
			return false
		}
		g.Cells[r][c].State = model.CellState(num)
	}
	return true
}

func displayResult(grid *model.Grid) {
	res := solver.Solve(grid)
	fmt.Println("\n--- 當前盤面與地雷機率 ---")

	// 1. 打印頂部列索引
	fmt.Print("    ") // 對應左側列索引的空白
	for c := 0; c < grid.Cols; c++ {
		fmt.Printf("   C%-2d  ", c) // 固定寬度 8
	}
	fmt.Println()

	// 2. 打印分隔線
	fmt.Print("    ")
	for c := 0; c < grid.Cols; c++ {
		fmt.Print("--------")
	}
	fmt.Println()

	// 3. 打印每一列內容
	for r := 0; r < grid.Rows; r++ {
		fmt.Printf("R%d |", r) // 左側列索引
		for c := 0; c < grid.Cols; c++ {
			state := grid.Cells[r][c].State
			var content string
			if state == model.StateUnknown {
				content = fmt.Sprintf("%5.1f%%", res.Probabilities[r][c]*100)
			} else if state == model.StateFlag {
				content = " FLAG "
			} else {
				content = fmt.Sprintf("  %d   ", int(state))
			}
			fmt.Printf(" %s |", content)
		}
		fmt.Println()

		// 打印列與列之間的分隔線
		fmt.Print("   ")
		for c := 0; c < grid.Cols; c++ {
			fmt.Print(" --------")
		}
		fmt.Println()
	}
}
