const but_lead = document.getElementById('leadboard')
const but_play = document.getElementById('play')
const but_quit = document.getElementById("quit")
let token = localStorage.getItem('token')
if (token==null) {
        window.location.href = "http://localhost:7010";

    }
const play = () => {
    window.location.href="/lobby/leadboard"
}
const lead = () => {
    window.location.href="/lobby/waiting/plaing"
}
const quit1 =()=>{
    localStorage.clear()
    window.location.href = "http://localhost:7010"
}
// Перевірка для кожної кнопки окремо
if (but_play) {
    but_play.addEventListener('click', lead);
}

if (but_lead) {
    but_lead.addEventListener('click', play);
}

if (but_quit) {
    but_quit.addEventListener('click', quit1);
}