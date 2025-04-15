document.addEventListener("DOMContentLoaded", function () {
    const playersList = document.getElementById("players-list");
    const validateButton = document.getElementById("validate-button");
    const songInput = document.getElementById("song-input");
    const timerElement = document.querySelector(".timer");
    const roundCounter = document.querySelector(".round-counter");
    const audioPlayer = document.getElementById("audio-player");
    const leaveButton = document.getElementById("leave-button");

    let playerName = "Nathaël";
    let timeLeft = 30;
    let countdown;

    function updateGameState() {
        fetch("/state")
            .then(response => response.json())
            .then(data => {
                updateUI(data);
            })
            .catch(error => console.error("Error fetching game state:", error));
    }

    function updateUI(state) {
        if (state.roundActive) {
            timeLeft = state.timeLeft;
            timerElement.textContent = `00:${timeLeft < 10 ? '0' : ''}${timeLeft}`;
        }

        roundCounter.textContent = `${state.currentRound}/${state.totalRounds}`;

        playersList.innerHTML = "";
        state.players.forEach(player => {
            const listItem = document.createElement("li");
            listItem.innerHTML = `
                <span>${player.name}</span>
                <span class="status ${player.status === "correct" ? "finished" : "not-finished"}">
                    ${player.status === "correct" ? "Correct" : player.status === "wrong" ? "Incorrect" : player.status === "timeout" ? "....." : "..."}
                </span>
                <span class="score">${player.score} pts</span>
            `;
            playersList.appendChild(listItem);
        });

        if (audioPlayer.src !== state.previewUrl && state.previewUrl) {
            audioPlayer.src = state.previewUrl;
            audioPlayer.play().catch(e => console.log("Audio play prevented:", e));
        }
    }

    validateButton.addEventListener("click", function() {
        const guess = songInput.value.trim();
        if (guess) {
            fetch("/guess", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    player: playerName,
                    answer: guess
                })
            })
            .then(response => response.json())
            .then(data => {
                updateUI(data);
                songInput.value = "";
            })
            .catch(error => console.error("Error submitting guess:", error));
        }
    });

    function startGame(players) {
        fetch("/start", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(players)
        })
        .then(response => response.json())
        .then(data => {
            updateUI(data);
        })
        .catch(error => console.error("Error starting game:", error));
    }

    startGame(["Quentin", "Ilian", "Nathaël"]);

    setInterval(updateGameState, 100);
});