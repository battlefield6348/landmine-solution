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
        
        // 點擊事件
        div.onclick = () => {
            cell.state = (cell.state + 2) % 10 - 1; // -1 -> 0 -> 1 ... -> 8 -> -1
            renderGrid();
        };

        // 右鍵事件
        div.oncontextmenu = (e) => {
            e.preventDefault();
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
