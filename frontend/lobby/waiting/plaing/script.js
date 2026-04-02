const canvas = document.getElementById("game");
const ctx = canvas.getContext("2d");

canvas.width = window.innerWidth;
canvas.height = window.innerHeight;

// масив об'єктів
let objects = [];

// генеруємо швидкість

function getSpeed() {
    let speed = (Math.random() - 0.5) * 4;

    // мінімальна швидкість
    if (Math.abs(speed) < 1 && Math.abs(speed) > 1) {
        speed = speed < 0 ? -1 : 1;
    }

    return speed;
}

// створюємо X та O
for (let i = 0; i < 200; i++) {
    objects.push({
        type: Math.random() > 0.5 ? "X" : "O",
        x: Math.random() * canvas.width,
        y: Math.random() * canvas.height,
        dx: getSpeed(),
        dy: getSpeed(),
        size: 20
    });
}

// функція малювання піксельного X
function drawX(x, y) {
    ctx.fillStyle = "red";

    const p = [
        [0,0],[1,1],[2,2],[3,3],[4,4],
        [4,0],[3,1],[2,2],[1,3],[0,4]
    ];

    p.forEach(([px, py]) => {
        ctx.fillRect(x + px * 4, y + py * 4, 4, 4);
    });
}

// функція малювання піксельного O
function drawO(x, y) {
    ctx.fillStyle = "blue";

    const p = [
        [1,0],[2,0],
        [0,1],[3,1],
        [0,2],[3,2],
        [1,3],[2,3]
    ];

    p.forEach(([px, py]) => {
        ctx.fillRect(x + px * 4, y + py * 4, 4, 4);
    });
}

// анімація
function animate() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    objects.forEach(obj => {
        obj.x += obj.dx;
        obj.y += obj.dy;

        // відштовхування від стін
        if (obj.x <= 0 || obj.x + obj.size >= canvas.width) {
            obj.dx *= -1;
        }

        if (obj.y <= 0 || obj.y + obj.size >= canvas.height) {
            obj.dy *= -1;
        }

        // малюємо
        if (obj.type === "X") {
            drawX(obj.x, obj.y, obj.size);
        } else {
            drawO(obj.x, obj.y, obj.size);
        }
    });

    requestAnimationFrame(animate);
}

animate();












// код до хрестиків ноликів

const boardElement = document.getElementById('game-board');
const cells = document.querySelectorAll('.cell');
const statusText = document.getElementById('status');
const resetBtn = document.getElementById('reset-btn');

let currentPlayer = "X"; // Починаємо з Х
let gameState = ["", "", "", "", "", "", "", "", ""]; // Пусте поле
let isGameActive = true;

// виграшні комбінації
const winningConditions = [
    [0, 1, 2], [3, 4, 5], [6, 7, 8], // Горизонталі
    [0, 3, 6], [1, 4, 7], [2, 5, 8], // Вертикалі
    [0, 4, 8], [2, 4, 6]             // Діагоналі
];

// клік по клітинці
function handleCellClick(e) {
    const clickedCell = e.target;
    const clickedIndex = parseInt(clickedCell.getAttribute('data-index'));

    if (gameState[clickedIndex] !== "" || !isGameActive) return;

    gameState[clickedIndex] = currentPlayer;
    clickedCell.innerText = currentPlayer;

    // ДОДАЄМО ЦЕ: встановлюємо клас для кольору
    clickedCell.classList.add(currentPlayer.toLowerCase());

    checkResult();
}

function checkResult() {
    let roundWon = false;

    for (let i = 0; i < winningConditions.length; i++) {
        const [a, b, c] = winningConditions[i];
        if (gameState[a] && gameState[a] === gameState[b] && gameState[a] === gameState[c]) {
            roundWon = true;
            break;
        }
    }

    if (roundWon) {
        statusText.innerText = `Гравець ${currentPlayer} переміг!`;
        isGameActive = false;
        return;
    }

    // Перевірка на нічию
    if (!gameState.includes("")) {
        statusText.innerText = "Нічия!";
        isGameActive = false;
        return;
    }

    // Зміна гравця
    currentPlayer = currentPlayer === "X" ? "O" : "X";
    statusText.innerText = `Зараз хід: ${currentPlayer}`;
}

// Очищення гри
function restartGame() {
    currentPlayer = "X";
    gameState = ["", "", "", "", "", "", "", "", ""];
    isGameActive = true;
    statusText.innerText = "";
    cells.forEach(cell => {
        cell.innerText = "";
        cell.classList.remove('x', 'o'); // ВИДАЛЯЄМО КЛАСИ
    });
}

cells.forEach(cell => cell.addEventListener('click', handleCellClick));
resetBtn.addEventListener('click', restartGame);