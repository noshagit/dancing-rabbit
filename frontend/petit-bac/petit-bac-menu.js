document.getElementById("rules-button").addEventListener("click", function () {
    document.getElementById("rules-popup").style.display = "block";
});

document.getElementById("close-popup").addEventListener("click", function () {
    document.getElementById("rules-popup").style.display = "none";
});

const timeArrow = document.getElementById("time-arrow");
const timeOptions = document.getElementById("time-options");
const selectedTime = document.getElementById("selected-time");
let maxRounds = 5;
let duration;

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
const MAX_RECONNECT_ATTEMPTS = 5;
let connectionAttempts = 0;
let ws;
let playerID;
let playerName;
let timerInterval;
let hasGuessedThisRound = false;

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
    timeArrow.classList.remove("arrow-down");
    timeArrow.classList.add("arrow-up");
});

document.getElementById("hard-button").addEventListener("click", function () {
    document.getElementById("rounds-input").value = "10";

    selectedTime.textContent = "00:30";

    timeOptions.style.display = "none";
    timeArrow.classList.remove("arrow-down");
    timeArrow.classList.add("arrow-up");
});

document.getElementById("easy-button").addEventListener("click", function () {
    document.getElementById("rounds-input").value = "5";

    selectedTime.textContent = "05:00";

    timeOptions.style.display = "none";
    timeArrow.classList.remove("arrow-down");
    timeArrow.classList.add("arrow-up");
});

document.addEventListener("DOMContentLoaded", () => {
    connectWebSocket();
});

function connectWebSocket() {
    if (connectionAttempts >= MAX_RECONNECT_ATTEMPTS) return;

    connectionAttempts++;
    console.log(`Connecting... (${connectionAttempts}/${MAX_RECONNECT_ATTEMPTS})`);

    ws = new WebSocket("ws://" + location.host + "/ws");

    ws.onopen = () => {
        console.log('Connected!');
        connectionAttempts = 0;
        registerPlayer();
    };

    ws.onmessage = event => {
        const message = JSON.parse(event.data);
        console.log('Received:', message.type);
        handleServerMessage(message);
    };

    ws.onclose = () => {
        console.log('Connection closed');
        setTimeout(connectWebSocket, 2000);
    };

    ws.onerror = error => console.error('WebSocket error:', error);
}

function registerPlayer() {
    ws.send(JSON.stringify({
        type: 'register_player',
        content: { id: playerID || "", name: playerName }
    }));
}

function handleServerMessage(message) {
    const { type, content } = message;

    switch (type) {

        case 'player_list':
            updatePlayerList(content);
            break;

        case 'player_assigned':
            playerID = content.id;
            localStorage.setItem('deafRhythmPlayerID', playerID);
            break;
    }
}

function submitGuess() {
    const input = document.getElementById('song-title');
    const guess = input.value.trim();

    if (hasGuessedThisRound) {
        showNotification('You already guessed the answer', 'error');
        return;
    }
    if (!guess) {
        showNotification('Please enter a guess', 'error');
        return;
    }
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        showNotification('Connection lost. Reconnecting...', 'error');
        connectWebSocket();
        return;
    }

    ws.send(JSON.stringify({ type: 'guess', content: guess }));
    input.value = '';
}

function updatePlayerList(players) {
    const playersList = document.getElementById('players-list');
    playersList.innerHTML = '';

    players.forEach(player => {
        const item = document.createElement('li');
        item.dataset.playerId = player.id;
        item.innerHTML = `
            <div class="player-info-row">
                <span class="player-name">${player.name}</span>
            </div>
        `;
        if (player.id === playerID) item.classList.add('current-player');
        playersList.appendChild(item);
    });
}

function startTimer(endTime) {
    if (timerInterval)
        clearInterval(timerInterval);
    hasGuessedThisRound = false;

    updateTimer(endTime);
    timerInterval = setInterval(() => updateTimer(endTime), 100);
}

function updateTimer(endTime) {
    const timerElement = document.getElementById('timer');
    const remaining = Math.max(0, Math.ceil((endTime - Date.now()) / 1000));
    const minutes = Math.floor(remaining / 60);
    const seconds = remaining % 60;

    let minutesDisplay = minutes < 10 ? "0" + minutes : minutes;
    let secondsDisplay = seconds < 10 ? "0" + seconds : seconds;
    timerElement.textContent = minutesDisplay + ":" + secondsDisplay;

    if (remaining <= 0) {
        clearInterval(timerInterval);
        timerElement.textContent = "00:00";
    }
}

function clearNotifications() {
    const notificationArea = document.getElementById('notification-area');
    const gameOverNotification = document.querySelector('.game-over-notification');
    notificationArea.innerHTML = '';
    if (gameOverNotification)
        notificationArea.appendChild(gameOverNotification);
}

function showNotification(message, type) {
    const notificationArea = document.getElementById('notification-area');
    const notification = document.createElement('div');
    notification.className = `notification ${type || ''}`;
    notification.textContent = message;
    notificationArea.appendChild(notification);
    setTimeout(() => notification.remove(), 3000);
}

function guessResult(data) {
    const notificationText = data.correct
        ? `Correct guess! +${data.points} points!`
        : 'Incorrect guess!';
    const notifClass = data.correct ? 'correct-guess' : 'incorrect-guess';

    showCustomNotification(notificationText, notifClass);

    if (data.correct) {
        hasGuessedThisRound = true;
        updatePlayerStatus(playerID, 'Finished', data.score);
    }
}

function displayGuessNotification(data) {
    if (data.playerID == playerID)
        return;

    showCustomNotification(`${data.playerName} made a guess!`);
    updatePlayerStatus(data.playerID, 'Finished', data.score);
}

function updatePlayerStatus(id, status, score) {
    const playerItem = document.querySelector(`#players-list li[data-player-id="${id}"]`);
    if (!playerItem)
        return;

    const statusElement = playerItem.querySelector('.status');
    if (statusElement) {
        statusElement.textContent = status;
        statusElement.className = `status ${status.toLowerCase()}`;
    }

    const scoreElement = playerItem.querySelector('.player-score');
    if (scoreElement && score !== undefined) {
        scoreElement.textContent = `${score} pts`;
    }
}

function showCustomNotification(message, className = '') {
    const notificationArea = document.getElementById('notification-area');
    const notification = document.createElement('div');
    notification.className = `notification ${className}`;
    notification.textContent = message;
    notificationArea.appendChild(notification);
    setTimeout(() => notification.remove(), 3000);
}

function showGameOverNotif() {
    const notification = document.createElement('div');
    notification.className = 'notification game-over-notification';
    notification.textContent = 'Game over! Redirecting to final scores...';
    document.getElementById('notification-area').appendChild(notification);
}
