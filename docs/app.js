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

let focusedIndex = null;

function renderGrid() {
    const gridEl = document.getElementById('grid');
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
                focusedIndex = Math.min(cells.length - 1, index + size);
            } else if (e.key === 'ArrowUp') {
                focusedIndex = Math.max(0, index - size);
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
    const result = solveMinesweeper(size, states);

    if (!result.solvable) {
        alert("⚠️ 偵測到邏輯矛盾！目前的盤面配置在踩地雷規則下是不可能的，請檢查數字與旗幟是否正確。");
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

// 啟動
initGrid();
