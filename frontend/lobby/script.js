const but_lead = document.getElementById('leadboard')
const but_play = document.getElementById('play')
const play = () => {
    window.location.href="/lobby/leadboard"
}
const lead = () => {
    window.location.href="/lobby/waiting/plaing"
}
but_play.addEventListener('click', play)
but_lead.addEventListener('click', lead)