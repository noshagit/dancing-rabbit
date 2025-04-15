document.getElementById("rules-button").addEventListener("click", function () {
    document.getElementById("rules-popup").style.display = "block";
});

document.getElementById("close-popup").addEventListener("click", function () {
    document.getElementById("rules-popup").style.display = "none";
});

const timeArrow = document.getElementById("time-arrow");
const timeOptions = document.getElementById("time-options");
const selectedTime = document.getElementById("selected-time");

timeArrow.addEventListener("click", function () {
    const isOpen = timeOptions.style.display === "block";
    timeOptions.style.display = isOpen ? "none" : "block";
    timeArrow.classList.toggle("arrow-down", !isOpen);
    timeArrow.classList.toggle("arrow-up", isOpen);
});

document.querySelectorAll(".time-button").forEach(button => {
    button.addEventListener("click", function () {
        selectedTime.textContent = this.textContent;
        timeOptions.style.display = "none";
        timeArrow.classList.remove("arrow-down");
        timeArrow.classList.add("arrow-up");
    });
});

const playersList = document.getElementById("players-list");

const players = ["Emilia", "Nathaël", "Ilian", "Natïha", "Corentin", "Goatin"];

function updatePlayersList(players) {
    playersList.innerHTML = "";
    players.forEach(player => {
        const li = document.createElement("li");
        li.textContent = player;
        playersList.appendChild(li);
    });
}

updatePlayersList(players);

document.getElementById("reset-button").addEventListener("click", function () {
    document.getElementById("rounds-input").value = "";

    selectedTime.textContent = "02:00";

    timeOptions.style.display = "none";
    modeOptions.style.display = "none";
    timeArrow.classList.remove("arrow-down");
    timeArrow.classList.add("arrow-up");
    modeArrow.classList.remove("arrow-down");
    modeArrow.classList.add("arrow-up");
});

document.getElementById("hard-button").addEventListener("click", function () {
    document.getElementById("rounds-input").value = "10";

    selectedTime.textContent = "00:30";

    timeOptions.style.display = "none";
    modeOptions.style.display = "none";
    timeArrow.classList.remove("arrow-down");
    timeArrow.classList.add("arrow-up");
    modeArrow.classList.remove("arrow-down");
    modeArrow.classList.add("arrow-up");
});

document.getElementById("easy-button").addEventListener("click", function () {
    document.getElementById("rounds-input").value = "5";

    selectedTime.textContent = "05:00";

    timeOptions.style.display = "none";
    modeOptions.style.display = "none";
    timeArrow.classList.remove("arrow-down");
    timeArrow.classList.add("arrow-up");
    modeArrow.classList.remove("arrow-down");
    modeArrow.classList.add("arrow-up");
});
document.getElementById("play_button").addEventListener("click", function () {
    window.location.href = "/petit-bac";
});

