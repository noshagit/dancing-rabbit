document.addEventListener("DOMContentLoaded", function () {
    const players = ["Emilia", "Quentin", "Nathaël", "Ilian", "Corentin"];
    const playersList = document.getElementById("players-list");
    const validateButton = document.getElementById("validate-button");

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

    validateButton.addEventListener("click", function () {
        const randomPlayer = players[Math.floor(Math.random() * players.length)];
        finishTurn(randomPlayer);
    });

    generatePlayerList();
});
