const go = new Go();
let rows = 3;
let cols = 3;
let cells = []; // 儲存每個格子的狀態
let focusedIndex = null;

// 初始化 WASM
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
    console.log("Go WASM Loaded");
});

function initGrid() {
    renderGrid();
}

function renderGrid() {
    const gridEl = document.getElementById('grid');
    gridEl.style.gridTemplateColumns = `repeat(${cols}, 50px)`;
    gridEl.innerHTML = '';

    cells.forEach((cell, index) => {
        const div = document.createElement('div');
        div.className = `cell ${getCellClass(cell.state)}`;
        div.innerText = getCellText(cell.state);
        div.setAttribute('tabindex', '0');
        div.dataset.index = index;

        // 點擊：僅聚焦，不改數字
        div.onclick = () => {
            focusedIndex = index;
            div.focus();
        };

        // 鍵盤事件
        div.onkeydown = (e) => {
            focusedIndex = index;
            let handled = true;

            if (e.key >= '0' && e.key <= '8') {
                cell.state = parseInt(e.key);
            } else if (e.key.toLowerCase() === 'f') {
                cell.state = -2;
            } else if (e.key.toLowerCase() === 'u' || e.key === 'Backspace' || e.key === 'Delete') {
                cell.state = -1;
            } else if (e.key.toLowerCase() === 'e') {
                cell.state = 0;
            } else if (e.key === 'ArrowRight') {
                focusedIndex = Math.min(cells.length - 1, index + 1);
            } else if (e.key === 'ArrowLeft') {
                focusedIndex = Math.max(0, index - 1);
            } else if (e.key === 'ArrowDown') {
                focusedIndex = Math.min(cells.length - 1, index + cols);
            } else if (e.key === 'ArrowUp') {
                focusedIndex = Math.max(0, index - cols);
            } else {
                handled = false;
            }

            if (handled) {
                e.preventDefault();
                renderGrid();
                // 渲染完後自動聚焦到新的位置
                const allCells = document.querySelectorAll('.cell');
                if (allCells[focusedIndex]) {
                    allCells[focusedIndex].focus();
                }
            }
        };

        // 右鍵：旗幟
        div.oncontextmenu = (e) => {
            e.preventDefault();
            focusedIndex = index;
            div.focus();
            cell.state = (cell.state === -2) ? -1 : -2;
            renderGrid();
            // 保持聚焦
            const allCells = document.querySelectorAll('.cell');
            allCells[focusedIndex].focus();
        };

        if (cell.probability !== undefined && cell.state === -1) {
            const probSpan = document.createElement('span');
            probSpan.className = 'prob';
            probSpan.innerText = (cell.probability * 100).toFixed(1) + '%';
            div.appendChild(probSpan);

            // 根據機率變色 (紅色越高越危險)
            if (cell.probability > 0) {
                const alpha = cell.probability * 0.5;
                div.style.backgroundColor = `rgba(248, 81, 73, ${alpha})`;
            }
        }

        gridEl.appendChild(div);
    });
}

// 初始化資料
for (let r = 0; r < rows; r++) {
    for (let c = 0; c < cols; c++) {
        cells.push({ r, c, state: -1 });
    }
}

function addRow(top) {
    const newRow = [];
    for (let c = 0; c < cols; c++) {
        newRow.push({ state: -1 });
    }
    if (top) {
        cells = [...newRow, ...cells];
    } else {
        cells = [...cells, ...newRow];
    }
    rows++;
    renderGrid();
}

function removeRow(top) {
    if (rows <= 1) return;
    if (top) {
        cells.splice(0, cols);
    } else {
        cells.splice((rows - 1) * cols, cols);
    }
    rows--;
    renderGrid();
}

function addCol(left) {
    for (let r = 0; r < rows; r++) {
        const index = left ? r * (cols + 1) : r * (cols + 1) + cols;
        cells.splice(index, 0, { state: -1 });
    }
    cols++;
    renderGrid();
}

function removeCol(left) {
    if (cols <= 1) return;
    for (let r = 0; r < rows; r++) {
        const index = left ? r * (cols - 1) : r * (cols - 1) + (cols - 1);
        cells.splice(index, 1);
    }
    cols--;
    renderGrid();
}

function resetGrid() {
    cells.forEach(c => c.state = -1);
    cells.forEach(c => c.probability = undefined);
    renderGrid();
}

function solve() {
    if (typeof solveMinesweeper !== 'function') {
        alert("WASM 尚未載入完成");
        return;
    }

    // 將狀態轉為 Go 預期的格式
    const states = cells.map(c => c.state);
    const result = solveMinesweeper(rows, cols, states);

    if (!result.solvable) {
        if (result.timeout) {
            alert("💤 運算超時！由於未知格連接過於複雜，計算量已超過上限。請嘗試先標記一些旗幟或縮小範圍。");
        } else {
            // 檢查是否是因為 len(unknownPos) 太大 (雖然現在有了優化，但還是保留一個保險)
            if (states.filter(s => s === -1).length > 100) {
                alert("💤 運算範圍過大！請嘗試先標記一些旗幟或縮小範圍。");
            } else {
                alert("⚠️ 偵測到邏輯矛盾！目前的盤面配置在踩地雷規則下是不可能的，請檢查數字與旗幟是否正確。");
            }
        }
        cells.forEach(c => c.probability = undefined);
        renderGrid();
        return;
    }

    const { probabilities } = result;

    // 更新結果
    cells.forEach((cell, i) => {
        cell.probability = probabilities[i];
    });
    renderGrid();
}

function getCellClass(state) {
    if (state === -1) return 'unknown';
    if (state === -2) return 'flag';
    if (state === 0) return 'empty';
    return `n${state}`;
}

function getCellText(state) {
    if (state === -1) return '';
    if (state === -2) return '🚩';
    if (state === 0) return '';
    return state;
}

// 啟動
initGrid();
