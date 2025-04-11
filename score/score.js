const gameTitle = "Blind Test";
document.getElementById("game-title").textContent = `Résultats du jeu : ${gameTitle}`;

const players = [
    { name: "TEAM", score: 120 },
    { name: "OE", score: 250 },
    { name: "LA", score: 180 }
];

players.sort((a, b) => b.score - a.score);

const scoreboardDiv = document.getElementById("scoreboard");
players.forEach(player => {
    const entry = document.createElement("div");
    entry.className = "score-entry";
    entry.innerHTML = `<span>${player.name}</span><span>${player.score} pts</span>`;
    scoreboardDiv.appendChild(entry);
});
