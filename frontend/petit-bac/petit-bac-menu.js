let socket;
let playerId;

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

    if (isOpen) {
        timeArrow.classList.remove("arrow-down");
        timeArrow.classList.add("arrow-up");
    } else {
        timeArrow.classList.add("arrow-down");
        timeArrow.classList.remove("arrow-up");
    }
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
const samplePlayers = ["Emilia", "Nathaël", "Ilian", "Natïha", "Corentin", "Goatin"];

function updatePlayersList(players) {
    playersList.innerHTML = "";
    players.forEach(player => {
        const li = document.createElement("li");
        li.textContent = player;
        playersList.appendChild(li);
    });
}

updatePlayersList(samplePlayers);
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

function extractTimeValue(timeStr) {
    const parts = timeStr.split(':');

    if (parts.length === 2) {
        const minutes = parseInt(parts[0], 10);
        const seconds = parseInt(parts[1], 10);
        return (minutes * 60) + seconds;
    }
    return 30;
}

function initWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsURL = `${protocol}//${window.location.host}/ws`;
    
    socket = new WebSocket(wsURL);
    
    socket.onopen = () => {
        console.log("WebSocket connection established");
    };
    
    socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        
        if (message.type === "player_assigned") {
            playerId = message.content.id;
            console.log("Assigned player ID:", playerId);
        } else if (message.type === "player_list") {
            const players = message.content.map(player => player.name);
            updatePlayersList(players);
        } else if (message.type === "game_start") {
            const rounds = message.content.rounds;
            const duration = message.content.duration;
            window.location.href = `/game?rounds=${rounds}&duration=${duration}`;
        }
    };
    
    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
        showNotification("Connection error. Please try again.", "error");
    };
    
    socket.onclose = () => {
        console.log("WebSocket connection closed");
        setTimeout(initWebSocket, 3000);
    };
}

document.addEventListener('DOMContentLoaded', initWebSocket);

document.getElementById("start-game-button").addEventListener("click", function () {
    const roundsInput = document.getElementById("rounds-input");
    const rounds = roundsInput.value ? parseInt(roundsInput.value, 10) : 5;

    const timeStr = selectedTime.textContent;
    const duration = extractTimeValue(timeStr);

    localStorage.setItem("petit-bac-rounds", rounds);
    localStorage.setItem("petit-bac-duration", duration);

    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({
            type: "start_game",
            content: {
                rounds: rounds,
                duration: duration
            }
        }));
    } else {
        showNotification("Connection issue. Please try again.", "error");
    }
});

function showNotification(message, type) {
    const notificationArea = document.getElementById('notification-area');
    const notification = document.createElement('div');
    notification.className = `notification ${type || ''}`;
    notification.textContent = message;

    notificationArea.appendChild(notification);
    setTimeout(() => notification.remove(), 3000);
}
