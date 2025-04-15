document.addEventListener("DOMContentLoaded", () => {
    const ws = new WebSocket(`ws://${window.location.hostname}:8080/bac/ws`);
    let playersData = [];
    let currentIndex = 0;
    const playerNameEl = document.getElementById("player-name");
    const answersContainer = document.getElementById("answers-container");
    const playersList = document.getElementById("players-list");
    const prevButton = document.getElementById("prev-player");
    const nextButton = document.getElementById("next-player");

    ws.onopen = () => {
        console.log("Connected to server!");
        ws.send(JSON.stringify({
            type: "request_vote_data"
        }));
    };

    ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        console.log("Server says:", message);

        if (message.type === "player_assigned") {
            const playerIdElement = document.createElement('div');
            playerIdElement.id = "my-player-id";
            playerIdElement.dataset.playerId = message.content.id;
            playerIdElement.style.display = 'none';
            document.body.appendChild(playerIdElement);
        }
        else if (message.type === "vote_data") {
            console.log("Vote data received!");
            playersData = prepareVoteData(message.content);
            updatePlayerList();
            showPlayerAnswers();
        }
        else if (message.type === "vote_results") {
            updateVoteCounts(message.content);
        }
        else if (message.type === "voting_complete") {
            showMessage("Voting complete! Going to scoreboard...");
            setTimeout(() => {
                window.location.href = "/scoreboard";
            }, 2000);
        }
    };

    function showMessage(text) {
        let notificationArea = document.getElementById('notification-area');
        if (!notificationArea) {
            notificationArea = document.createElement('div');
            notificationArea.id = 'notification-area';
            notificationArea.style.position = 'fixed';
            notificationArea.style.top = '20px';
            notificationArea.style.right = '20px';
            document.body.appendChild(notificationArea);
        }
        const message = document.createElement('div');
        message.textContent = text;
        message.style.padding = '10px';
        message.style.margin = '5px';
        message.style.backgroundColor = '#333';
        message.style.color = 'white';
        message.style.borderRadius = '5px';

        notificationArea.appendChild(message);
        setTimeout(() => message.remove(), 3000);
    }

    function prepareVoteData(serverData) {
        const result = [];
        for (const [playerId, rounds] of Object.entries(serverData)) {
            if (!rounds || rounds.length === 0) continue;
            rounds.forEach((round, roundIndex) => {
                if (round && round.Answers) {
                    const playerEntry = {
                        id: playerId,
                        name: "Player " + playerId,
                        round: roundIndex + 1,
                        letter: round.Letter || "?",
                        answers: [],
                        categories: [],
                        votes: []
                    };
                    for (const [category, answer] of Object.entries(round.Answers)) {
                        playerEntry.answers.push(answer);
                        playerEntry.categories.push(category);
                        playerEntry.votes.push(null);
                    }
                    result.push(playerEntry);
                }
            });
        }
        if (result.length === 0) {
            result.push({
                id: "example",
                name: "Example Player",
                round: 1,
                letter: "A",
                answers: ["Apple", "America", "Amazing"],
                categories: ["Fruit", "Country", "Adjective"],
                votes: [null, null, null]
            });
        }

        return result;
    }

    function updatePlayerList() {
        playersList.innerHTML = "";
        const players = {};
        playersData.forEach(entry => {
            if (!players[entry.id]) {
                players[entry.id] = {
                    id: entry.id,
                    name: entry.name,
                    rounds: []
                };
            }

            players[entry.id].rounds.push({
                round: entry.round,
                done: entry.votes.every(v => v !== null)
            });
        });
        Object.values(players).forEach(player => {
            const allDone = player.rounds.every(r => r.done);

            const li = document.createElement("li");
            li.dataset.playerId = player.id;
            if (allDone) li.classList.add("validated");

            const nameSpan = document.createElement("span");
            nameSpan.classList.add("player-name");
            nameSpan.textContent = player.name;
            li.appendChild(nameSpan);

            const roundsDiv = document.createElement("div");
            roundsDiv.classList.add("rounds-status");
            player.rounds.forEach(r => {
                const badge = document.createElement("span");
                badge.classList.add("round-badge");
                if (r.done) badge.classList.add("validated");
                badge.dataset.round = r.round;
                badge.textContent = "R" + r.round;
                roundsDiv.appendChild(badge);
            });
            li.appendChild(roundsDiv);

            const statusSpan = document.createElement("span");
            statusSpan.classList.add("status");
            statusSpan.textContent = allDone ? "Validé" : "À valider";
            li.appendChild(statusSpan);

            li.addEventListener("click", () => {
                const index = playersData.findIndex(p => p.id === player.id);
                if (index >= 0) {
                    currentIndex = index;
                    showPlayerAnswers();
                }
            });

            playersList.appendChild(li);
        });
    }

    function showPlayerAnswers() {
        if (!playersData.length) {
            playerNameEl.textContent = "No players data";
            answersContainer.innerHTML = "<p>No answer data available</p>";
            return;
        }
        if (currentIndex >= playersData.length) {
            currentIndex = 0;
        }

        const player = playersData[currentIndex];

        prevButton.disabled = currentIndex <= 0;
        nextButton.disabled = currentIndex >= playersData.length - 1;

        playerNameEl.innerHTML = `${player.name} <span class="round-info">Round ${player.round} - Letter ${player.letter}</span>`;

        answersContainer.innerHTML = "";

        player.answers.forEach((answer, idx) => {
            const answerDiv = document.createElement("div");
            answerDiv.className = "answer-item";

            answerDiv.innerHTML = `
                <div class="answer-text">${player.categories[idx]}: ${answer}</div>
                <div class="vote-buttons">
                    <button class="correct">Correct</button>
                    <button class="incorrect">Incorrect</button>
                </div>
            `;

            const correctBtn = answerDiv.querySelector(".correct");
            const incorrectBtn = answerDiv.querySelector(".incorrect");

            if (player.votes[idx] === true) {
                correctBtn.style.backgroundColor = "#8FCB9B";
            } else if (player.votes[idx] === false) {
                incorrectBtn.style.backgroundColor = "#F37C7C";
            }
            correctBtn.addEventListener("click", () => {
                player.votes[idx] = true;
                correctBtn.style.backgroundColor = "#8FCB9B";
                incorrectBtn.style.backgroundColor = "";
                updatePlayerList();

                sendVote(player.id, player.categories[idx], player.round, true);
            });

            incorrectBtn.addEventListener("click", () => {
                player.votes[idx] = false;
                incorrectBtn.style.backgroundColor = "#F37C7C";
                correctBtn.style.backgroundColor = "";
                updatePlayerList();

                sendVote(player.id, player.categories[idx], player.round, false);
            });

            answersContainer.appendChild(answerDiv);
        });
        document.querySelectorAll('#players-list li').forEach(li => {
            li.classList.remove('active');
            if (li.dataset.playerId === player.id) {
                li.classList.add('active');
            }
        });
    }

    function sendVote(targetPlayerId, category, round, isValid) {
        const myId = document.getElementById("my-player-id")?.dataset.playerId || "unknown";

        ws.send(JSON.stringify({
            type: "voted",
            content: {
                playerID: myId,
                targetPlayerID: targetPlayerId,
                category: category,
                round: round,
                valid: isValid
            }
        }));
    }

    function updateVoteCounts(results) {
        const player = playersData[currentIndex];
        if (!player) return;

        const answerItems = document.querySelectorAll('.answer-item');

        answerItems.forEach((item, idx) => {
            const category = player.categories[idx];

            let voteKey = "";
            for (const key in results[player.id]) {
                if (key.includes(category)) {
                    voteKey = key;
                    break;
                }
            }

            if (voteKey && results[player.id][voteKey]) {
                const validVotes = results[player.id][voteKey].valid || 0;
                const invalidVotes = results[player.id][voteKey].invalid || 0;

                let countDiv = item.querySelector('.vote-count');
                if (!countDiv) {
                    countDiv = document.createElement('div');
                    countDiv.className = 'vote-count';
                    item.appendChild(countDiv);
                }

                countDiv.innerHTML = `
                    <span style="color: green;">✓ ${validVotes}</span>
                    <span style="color: red;">✗ ${invalidVotes}</span>
                `;

                if (validVotes > invalidVotes) {
                    item.style.backgroundColor = "rgba(143, 203, 155, 0.3)";
                } else if (invalidVotes > validVotes) {
                    item.style.backgroundColor = "rgba(243, 124, 124, 0.3)";
                } else {
                    item.style.backgroundColor = "";
                }
            }
        });
    }

    prevButton.addEventListener("click", () => {
        if (currentIndex > 0) {
            currentIndex--;
            showPlayerAnswers();
        }
    });

    nextButton.addEventListener("click", () => {
        if (currentIndex < playersData.length - 1) {
            currentIndex++;
            showPlayerAnswers();
        }
    });
    try {
        updatePlayerList();
        showPlayerAnswers();
    } catch (err) {
        console.log("Waiting for data from server...");
        playerNameEl.textContent = "Waiting for data...";
        answersContainer.innerHTML = "<p>Waiting for player data. It will appear soon!</p>";
    }
});