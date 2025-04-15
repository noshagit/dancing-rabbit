document.addEventListener("DOMContentLoaded", function () {
    let maxRounds = 5;
    let roundDuration = 30;

    if (window.maxRounds)
        maxRounds = window.maxRounds;
    if (window.roundDuration)
        roundDuration = window.roundDuration;

    if (localStorage.getItem("petit-bac-rounds"))
        maxRounds = parseInt(localStorage.getItem("petit-bac-rounds"), 10);
    if (localStorage.getItem("petit-bac-duration"))
        roundDuration = parseInt(localStorage.getItem("petit-bac-duration"), 10);

    console.log(`Game will have ${maxRounds} rounds, ${roundDuration} seconds each`);

    const ws = new WebSocket(`ws://${window.location.host}/ws`);
    let hasSubmitted = false;
    let playerId = "";
    const playersList = document.getElementById("players-list");
    const submitButton = document.getElementById("validate-button");
    const timerElement = document.querySelector(".timer");

    ws.onopen = () => {
        console.log('Connected to server!');
        setTimeout(() => {
            ws.send(JSON.stringify({
                type: 'request_game_params'
            }));
        }, 1000);
        registerPlayer();
    };
    ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        console.log("Server says:", message);
        handleServerMessage(message);
    };

    function handleServerMessage(message) {
        const type = message.type;
        const content = message.content;
        if (type === 'player_list') {
            updatePlayerList(content);

        } else if (type === 'player_assigned') {
            playerId = content.id;
            localStorage.setItem('deafRhythmPlayerID', playerId);

        } else if (type === 'game_params') {
            updateGameDisplay(content);

        } else if (type === 'timer_start') {
            startTimer(content.timerEnd);
            updateRoundCounter(content.round, content.maxRounds);

        } else if (type === 'end_game') {
            showMessage('Game Over! Going to voting screen...', 'end-game');
            setTimeout(() => {
                window.location.href = '/vote';
            }, 1000);

        } else if (type === 'player_submitted') {
            playerSubmitted(content);

        } else if (type === 'round_reset') {
            newRoundStarted(content);
        }
    }

    function registerPlayer() {
        if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({
                type: 'register_player',
                content: { name: "Player" + Math.floor(Math.random() * 1000) }
            }));
        }
    }

    function submitAnswers() {
        if (hasSubmitted) {
            showMessage('You already submitted answers this round', 'error');
            return;
        }
        const answers = {};
        let allFilled = true;
        document.querySelectorAll('.answer-input').forEach(input => {
            const category = input.dataset.category;
            const answer = input.value.trim();

            if (!answer) {
                allFilled = false;
            }

            answers[category] = answer;
        });
        if (!allFilled) {
            showMessage('Please fill in all categories', 'error');
            return;
        }
        if (ws.readyState !== WebSocket.OPEN) {
            showMessage('Connection lost.', 'error');
            return;
        }
        ws.send(JSON.stringify({
            type: 'guess',
            content: answers
        }));
        hasSubmitted = true;
        showMessage('Answers submitted!', 'success');
    }

    function updatePlayerList(players) {
        playersList.innerHTML = '';
        players.forEach(player => {
            const listItem = document.createElement("li");
            listItem.innerHTML = `
                <span>${player.name}</span>
                <span class="status not-finished">Not finished</span>
                <span class="score">${player.score || 0}</span>
            `;
            playersList.appendChild(listItem);
        });
    }

    function updateGameDisplay(params) {
        const letterElement = document.getElementById('current-letter');
        if (letterElement) {
            letterElement.textContent = `Lettre : ${params.letter || 'A'}`;
        } 
        const choicesWrapper = document.querySelector('.choices-wrapper');
        if (choicesWrapper && params.categories) {
            choicesWrapper.innerHTML = '';
            params.categories.forEach((category) => {
                const div = document.createElement('div');
                div.className = 'choice';
                div.innerHTML = `
                    <div class="choice-label">${category}</div>
                    <input type="text" class="answer-input" data-category="${category}" placeholder="Entrez votre réponse...">
                `;
                choicesWrapper.appendChild(div);
            });
        }
        const roundCounter = document.querySelector('.round-counter');
        if (roundCounter) {
            roundCounter.textContent = `1/${maxRounds}`;
        }
    }

    function startTimer(endTime) {
        function updateTime() {
            const now = Date.now();
            const timeLeft = Math.max(0, Math.ceil((endTime - now) / 1000));
            const seconds = timeLeft % 60;
            const minutes = Math.floor(timeLeft / 60);
            timerElement.textContent = `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
            if (timeLeft > 0) {
                requestAnimationFrame(updateTime);
            }
        }
        hasSubmitted = false;
        updateTime();
    }

    function updateRoundCounter(currentRound, maxRounds) {
        const counter = document.querySelector('.round-counter');
        if (counter) {
            counter.textContent = `${currentRound}/${maxRounds}`;
        }
    }

    function showMessage(message, type) {
        const messageArea = document.getElementById('notification-area');

        const notification = document.createElement('div');
        notification.className = `notification ${type || ''}`;
        notification.textContent = message;

        messageArea.appendChild(notification);
        setTimeout(() => notification.remove(), 3000);
    }

    function playerSubmitted(playerName) {
        document.querySelectorAll('.answer-input').forEach(input => {
            input.value = '';
            input.disabled = true;
        });
        hasSubmitted = false;
        showMessage(`${playerName} has submitted! Next round soon...`, 'info');

        const playerItems = document.querySelectorAll('#players-list li');
        playerItems.forEach(item => {
            const nameSpan = item.querySelector('span:first-child');
            if (nameSpan.textContent === playerName) {
                const statusSpan = item.querySelector('.status');
                statusSpan.textContent = 'Submitted';
                statusSpan.classList.remove('not-finished');
                statusSpan.classList.add('finished');
            }
        });
    }

    function newRoundStarted(letter) {
        document.querySelectorAll('.answer-input').forEach(input => {
            input.disabled = false;
        });
        const letterElement = document.getElementById('current-letter');
        if (letterElement) {
            letterElement.textContent = `Lettre : ${letter}`;
        }

        hasSubmitted = false;

        showMessage(`New round started with letter: ${letter}`, 'info');
    }
    submitButton.addEventListener("click", submitAnswers);
});