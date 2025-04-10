document.addEventListener("DOMContentLoaded", () => {
    const players = [
        { name: "Ilian", answers: ["OE", "LA", "TEAM"], votes: [null, null, null] },
        { name: "Quentin", answers: ["Quentin", "LE", "GOAAAAT"], votes: [null, null, null] },
        { name: "Nathaël", answers: ["LAPIN", "BUNNY", "RABBIT"], votes: [null, null, null] },
        { name: "Léo", answers: ["LEO", "LE", "NPC"], votes: [null, null, null] }
    ];

    let currentIndex = 0;

    const playerNameEl = document.getElementById("player-name");
    const answersContainer = document.getElementById("answers-container");
    const playersList = document.getElementById("players-list");

    function renderPlayerList() {
        playersList.innerHTML = "";
        players.forEach((player) => {
            const validated = player.votes.every(v => v !== null);
            const li = document.createElement("li");
            li.innerHTML = `
                <span>${player.name}</span>
                <span class="status">${validated ? "Validé" : "À valider"}</span>
            `;
            playersList.appendChild(li);
        });
    }

    function renderAnswers() {
        const player = players[currentIndex];
        playerNameEl.textContent = player.name;
        answersContainer.innerHTML = "";

        player.answers.forEach((answer, idx) => {
            const answerDiv = document.createElement("div");
            answerDiv.className = "answer-item";
            const currentVote = player.votes[idx];

            answerDiv.innerHTML = `
                <div class="answer-text">${answer}</div>
                <div class="vote-buttons">
                    <button class="correct">Correct</button>
                    <button class="incorrect">Mauvais</button>
                </div>
            `;

            const correctBtn = answerDiv.querySelector(".correct");
            const incorrectBtn = answerDiv.querySelector(".incorrect");

            if (currentVote === true) {
                correctBtn.style.backgroundColor = "#8FCB9B";
            } else if (currentVote === false) {
                incorrectBtn.style.backgroundColor = "#F37C7C";
            }

            correctBtn.addEventListener("click", () => {
                player.votes[idx] = true;
                correctBtn.style.backgroundColor = "#8FCB9B";
                incorrectBtn.style.backgroundColor = "#A888B5";
                renderPlayerList();
            });

            incorrectBtn.addEventListener("click", () => {
                player.votes[idx] = false;
                incorrectBtn.style.backgroundColor = "#F37C7C";
                correctBtn.style.backgroundColor = "#A888B5";
                renderPlayerList();
            });

            answersContainer.appendChild(answerDiv);
        });
    }

    document.getElementById("prev-player").addEventListener("click", () => {
        if (currentIndex > 0) {
            currentIndex--;
            renderAnswers();
        }
    });

    document.getElementById("next-player").addEventListener("click", () => {
        if (currentIndex < players.length - 1) {
            currentIndex++;
            renderAnswers();
        }
    });

    renderPlayerList();
    renderAnswers();
});
