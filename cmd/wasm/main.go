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
	if len(args) < 2 {
		return nil
	}

	size := args[0].Int()
	inputCells := args[1]

	grid := model.NewGrid(size)
	for i := 0; i < size*size; i++ {
		r, c := i/size, i%size
		stateVal := inputCells.Index(i).Int()
		grid.Cells[r][c].State = model.CellState(stateVal)
	}

	result := solver.Solve(grid)

	// 將結果轉換為 JavaScript 陣列
	probabilities := js.Global().Get("Array").New(size * size)
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			probabilities.SetIndex(r*size+c, result.Probabilities[r][c])
		}
	}

	return probabilities
}
