function init(data) {
    const gameData = JSON.parse(data);
    const gameTitle = gameData.game;
    document.getElementById("game-title").textContent = `Résultats du jeu : ${gameTitle}`;

    const players = gameData.players || [];
    players.sort((a, b) => b.score - a.score);

    const scoreboardDiv = document.getElementById("scoreboard");
    players.forEach(player => {
        const entry = document.createElement("div");
        entry.className = "score-entry";
        entry.innerHTML = `<span>${player.name}</span><span>${player.score} pts</span>`;
        scoreboardDiv.appendChild(entry);
    });
}