//go:build js && wasm

package main

import (
	"fmt"
	"landmine-solution/internal/model"
	"landmine-solution/internal/solver"
	"syscall/js"
)

func main() {
	fmt.Println("WASM Go Initialized")

	// 導出 solveMinesweeper 函式給 JavaScript 呼叫
	js.Global().Set("solveMinesweeper", js.FuncOf(solve))

	// 保持程式不退出
	select {}
}

func solve(this js.Value, args []js.Value) any {
	if len(args) < 3 {
		return nil
	}

	rows := args[0].Int()
	cols := args[1].Int()
	inputCells := args[2]

	grid := model.NewGrid(rows, cols)
	for i := 0; i < rows*cols; i++ {
		r, c := i/cols, i%cols
		stateVal := inputCells.Index(i).Int()
		grid.Cells[r][c].State = model.CellState(stateVal)
	}

	result := solver.Solve(grid)

	// 將結果轉換為 JavaScript 陣列
	probabilities := js.Global().Get("Array").New(rows * cols)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			probabilities.SetIndex(r*cols+c, result.Probabilities[r][c])
		}
	}

	// 將結果封裝成 JavaScript 物件
	response := js.Global().Get("Object").New()
	response.Set("probabilities", probabilities)
	response.Set("timeout", result.Timeout)
	response.Set("solvable", result.Solvable)

	return response
}
