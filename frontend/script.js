//Вікна
const loginWindow = document.getElementById('login_window')
const regWindow = document.getElementById('registration_window')
//Кнопки зміна вікна
const but_changetologin = document.getElementById('change_window_button')
const but_changetoreg = document.getElementById('change_window_button2')
//Поля Вводу
const input_check_password = document.getElementById('check_password');
const matchError = document.getElementById('match_error');
const input_regname = document.getElementById('regname')
const input_regpassword = document.getElementById('regpassword')
const input_name = document.getElementById('name')
const input_password = document.getElementById('password')
//Кнопки реєстації логіну
const but_reg = document.getElementById('register_button')
const but_login = document.getElementById('log_button')
let token = localStorage.getItem('token')
if (token!=null) {
        window.location.href = "http://localhost:7010/lobby";

    }
//кусок для зміни регістраці(логіну)----------------------------
const change_window = () => {
    if (regWindow.style.display === 'none') {
        regWindow.style.display = 'block'
        loginWindow.style.display = 'none'
    }
    else{
        regWindow.style.display = 'none'
        loginWindow.style.display = 'block'
    }
}
but_changetologin.addEventListener('click', change_window)
but_changetoreg.addEventListener('click', change_window)
// -------------------------------------------------------------


const passwordError = document.getElementById("password_error");


// Окрема функція для перевірки паролів
const validatePasswords = () => {
    const password = input_regpassword.value;
    const confirm_password = input_check_password.value;

    // Перевірка 1: якщо нема букв
    if (!/[a-zA-Z]/.test(password)) {
        passwordError.textContent = "Пароль має містити хоча б одну букву";
        passwordError.style.display = "block";
        return false; // Кажемо "Ні, є помилка"
    } else {
        passwordError.style.display = "none";
    }

    // Перевірка 2: чи співпадають паролі
    if (password !== confirm_password) {
        matchError.style.display = "block";
        return false; // Кажемо "Ні, є помилка"
    } else {
        matchError.style.display = "none";
    }

    // Якщо дійшли сюди — все ідеально
    return true;
}

const register = async () => {

    // перевірка
    if (!validatePasswords()) {
        return;
    }

    const reguser_data = {
        name: input_regname.value,
        password: input_regpassword.value, }
    const response = await fetch('http://localhost:8080/register', {
        method: 'post',
        headers: {
            'Content-type': 'application/json',

        },
        body: JSON.stringify(reguser_data)
    })
    if (!response.ok) {
        throw new Error('Помилка сервера');
        }
    const data = await response.json();
    console.log("Ответ от сервера:", data);

    // Теперь можно сохранять данные в переменную
    const userName = data.name;
    const Token = data.token;
    const status=data.status;
    localStorage.setItem('token', Token);
    localStorage.setItem('userName', userName);
    console.log("Имя пользователя:", userName);
    console.log("токен:", Token);
    console.log("статус", status);
    if (status==='success') {
        window.location.href = "http://localhost:7010/lobby";

    }
}
const login = async () => {
    const loginuser_data = {
        name: input_name.value,
        password: input_password.value,
    }
    try {
        const response = await fetch('http://localhost:8080/login', {
            method: 'POST',
            headers: { 'Content-type': 'application/json' },
            body: JSON.stringify(loginuser_data)
        });

        if (!response.ok) {

            const errorData = await response.text();
            console.error("Сервер повернув помилку:", errorData);

            return;
        }

        const data = await response.json();
        const userName = data.name;
        const Token = data.token;
        const status=data.status;
        localStorage.setItem('token', Token);
        localStorage.setItem('userName', userName);
        console.log("Успішний вхід:", data);
        if (status==='success') {
            window.location.href = "http://localhost:7010/lobby";

        }
    } catch (err) {

        console.error("Помилка запиту:", err);
    }
};

but_login.addEventListener('click', login)
but_reg.addEventListener('click', register)

// перевірка надійності паролю



















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