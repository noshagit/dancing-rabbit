let ws;
let timerInterval;
let playerID = localStorage.getItem('deafRhythmPlayerID');
let playerName = localStorage.getItem('deafRhythmPlayerName') || generateRandomName();
let roundNumber = 0;
let hasGuessedThisRound = false;
let connectionAttempts = 0;
const MAX_RECONNECT_ATTEMPTS = 5;

function generateRandomName() {
    const adjectives = ["Swift", "Clever", "Funky", "Jazzy", "Melodic"];
    const nouns = ["Singer", "Dancer", "Listener", "Player", "Musician"];
    
    const adj = adjectives[Math.floor(Math.random() * adjectives.length)];
    const noun = nouns[Math.floor(Math.random() * nouns.length)];

    return `${adj}${noun}${Math.floor(Math.random() * 100)}`;
}

document.addEventListener('DOMContentLoaded', () => {
    connectWebSocket();
    document.getElementById('submit-button').addEventListener('click', submitGuess);
    document.getElementById('song-title').addEventListener('keypress', e => e.key === 'Enter' && submitGuess());
    document.getElementById('round-counter').textContent = `${roundNumber}/10`;
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
    setTimeout(requestLyrics, 1000);
}

function requestLyrics() {
    if (ws?.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: 'request_lyrics' }));
    }
}

function handleServerMessage(message) {
    const { type, content } = message;

    switch (type) {
        case 'lyrics':
            updateLyrics(content);
            break;

        case 'timer_start':
            clearNotifications();
            startTimer(content.timerEnd);
            document.getElementById('round-counter').textContent = `${content.round}/${content.maxRounds}`;
            document.querySelectorAll('.status').forEach(status => {
                status.textContent = 'Not finished';
                status.className = 'status not-finished';
            });
            hasGuessedThisRound = false;
            break;

        case 'timer_end':
            clearNotifications();
            clearInterval(timerInterval);
            document.getElementById('timer').textContent = "00:00";
            hasGuessedThisRound = false;
            showNotification(
                content.skipped
                    ? 'Everyone guessed correctly! Skipping to the next song.'
                    : `The song was: "${content.songName}" by ${content.artistName}`,
                'info'
            );
            break;

        case 'player_guessed':
            displayGuessNotification(content);
            break;

        case 'guess_result':
            guessResult(content);
            break;

        case 'player_list':
            updatePlayerList(content);
            break;

        case 'player_assigned':
            playerID = content.id;
            localStorage.setItem('deafRhythmPlayerID', playerID);
            document.getElementById('player-name').textContent = playerName;
            break;

        case 'game_over':
            showGameOverNotif();
            setTimeout(() => window.location.href = content, 3000);
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
                <span class="player-score">${player.score || 0} pts</span>
            </div>
            <span class="status not-finished">Not finished</span>
        `;
        if (player.id === playerID) item.classList.add('current-player');
        playersList.appendChild(item);
    });
}

function updateLyrics(lyrics) {
    const lyricsDisplay = document.getElementById('lyrics-display');
    lyricsDisplay.innerHTML = '';

    if (!lyrics?.length) {
        lyricsDisplay.innerHTML = '<div>Waiting for lyrics...</div>';
        return;
    }

    lyrics.forEach(line => {
        const paragraph = document.createElement('div');
        paragraph.textContent = line;
        lyricsDisplay.appendChild(paragraph);
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
