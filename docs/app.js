const go = new Go();
let size = 3;
let cells = []; // 儲存每個格子的狀態

// 初始化 WASM
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance);
    console.log("Go WASM Loaded");
});

function initGrid() {
    const gridEl = document.getElementById('grid');
    gridEl.style.gridTemplateColumns = `repeat(${size}, 50px)`;
    gridEl.innerHTML = '';

    // 初始化資料
    const newCells = [];
    for (let r = 0; r < size; r++) {
        for (let c = 0; c < size; c++) {
            // 嘗試保留舊有的值
            let existing = cells.find(it => it.r === r && it.c === c);
            newCells.push({
                r, c,
                state: existing ? existing.state : -1 // 預設未解 (-1)
            });
        }
    }
    cells = newCells;
    renderGrid();
}

function renderGrid() {
    const gridEl = document.getElementById('grid');
    gridEl.innerHTML = '';

    cells.forEach((cell, index) => {
        const div = document.createElement('div');
        div.className = `cell ${getCellClass(cell.state)}`;
        div.innerText = getCellText(cell.state);
        div.setAttribute('tabindex', '0'); // 使格子可以被聚焦（鍵盤操作）

        // 點擊事件：聚焦並循環切換（保留原意）
        div.onclick = () => {
            div.focus();
            cell.state = (cell.state + 2) % 10 - 1; // -1 -> 0 -> 1 ... -> 8 -> -1
            renderGrid();
        };

        // 鍵盤事件：支援直接輸入
        div.onkeydown = (e) => {
            if (e.key >= '0' && e.key <= '8') {
                cell.state = parseInt(e.key);
            } else if (e.key.toLowerCase() === 'f') {
                cell.state = -2; // 旗幟
            } else if (e.key.toLowerCase() === 'u' || e.key === 'Backspace' || e.key === 'Delete') {
                cell.state = -1; // 未解
            } else if (e.key.toLowerCase() === 'e') {
                cell.state = 0; // 空白
            } else {
                return; // 其他按鍵不處理
            }
            e.preventDefault();
            renderGrid();

            // 保持聚焦在下一個格子（選配：優化體驗）
            const nextIdx = index + 1;
            if (nextIdx < cells.length) {
                setTimeout(() => {
                    const allCells = document.querySelectorAll('.cell');
                    allCells[nextIdx].focus();
                }, 0);
            }
        };

        // 右鍵事件：切換旗幟
        div.oncontextmenu = (e) => {
            e.preventDefault();
            div.focus();
            if (cell.state === -2) cell.state = -1;
            else cell.state = -2;
            renderGrid();
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

function changeSize(delta) {
    size = Math.max(1, size + delta);
    initGrid();
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
    const probs = solveMinesweeper(size, states);

    // 更新結果
    cells.forEach((cell, i) => {
        cell.probability = probs[i];
    });
    renderGrid();
}

// 啟動
initGrid();
