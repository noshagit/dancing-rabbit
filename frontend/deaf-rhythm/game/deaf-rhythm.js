document.addEventListener("DOMContentLoaded", function () {
    let timeLeft = 30;
    const timerElement = document.querySelector(".timer");
    const lyricsDisplay = document.getElementById("lyrics-display");
    const playersList = document.getElementById("players-list");
    const validateButton = document.getElementById("validate-button");

    const lyrics = [
        "OEEEEEEEEEEEEEEEEEEE LA TEAM",
        "OEEEEEEEEEEEEEEEEEE LA TEAM",
        "OEEEEEEEEEE LA TEAM",
        "OEEEEEEEEEEEEEEE LA TEAM",
        "OEEEEEEEEEEEEEEEE LA TEAM",
        "OEEEEEEEEEEEEEEEEEE LA TEAM",
        "OEEEEEEEEEEEEE LA TEAM",
        "OEEEEEEEEEEEEEEEE LA TEAM",
        "OEEEEEEEEEEEEEEE LA TEAM",
        "OEEEEEEEEEEEEEEEEE LA TEAM"
    ];

    lyricsDisplay.innerHTML = lyrics.map(line => `<p class="lyrics-line">${line}</p>`).join('');

    const players = ["Emilia", "Quentin", "Nathaël", "Ilian", "Corentin"];

    function generatePlayerList() {
        players.forEach(player => {
            const listItem = document.createElement("li");
            listItem.innerHTML = `
                <span>${player}</span>
                <span class="status not-finished">Not finished</span>
            `;
            playersList.appendChild(listItem);
        });
    }

    function finishTurn(playerName) {
        const items = document.querySelectorAll("#players-list li");
        items.forEach(item => {
            if (item.textContent.includes(playerName)) {
                const status = item.querySelector(".status");
                status.textContent = "Finished";
                status.classList.remove("not-finished");
                status.classList.add("finished");
            }
        });
    }

    generatePlayerList();

    const countdown = setInterval(() => {
        timeLeft--;
        timerElement.textContent = `00:${timeLeft < 10 ? '0' : ''}${timeLeft}`;

        if (timeLeft <= 0) {
            clearInterval(countdown);
            timerElement.textContent = "00:00";
        }
    }, 1000);

    validateButton.addEventListener("click", function () {
        finishTurn(players[Math.floor(Math.random() * players.length)]);
    });
});
