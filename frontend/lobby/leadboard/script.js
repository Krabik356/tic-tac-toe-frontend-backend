async function loadLeaders() {
    try {
        const token = localStorage.getItem('token');
        const response = await fetch('http://localhost:8080/getLeaderBoard',{
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'session_token': token
            }
        });
        const dat = await response.json();
        console.log("Реальні дані з сервера:", dat.lb.data); // ОБОВ'ЯЗКОВО подивись сюди в консолі
        const data = dat.lb.data
        // Перевіряємо, чи існує data.lb і чи є він масивом

        const listContainer = document.querySelector('.field-list ul');
        listContainer.innerHTML = '';
        data.forEach((player, index) => {
        const li = document.createElement('li');
        li.innerHTML = `
            <span>${index + 1}</span>
            <span>${player.name}</span>
            <span>${player.rank}</span>
        `;
        listContainer.appendChild(li);
    });
    } catch (error) {
        console.error("Помилка при завантаженні лідерів:", error);
    }
}

document.addEventListener('DOMContentLoaded', loadLeaders);