function showPopup(title, text) {
        document.getElementById("popup-title").textContent = title;
        document.getElementById("popup-text").textContent = text;
        document.getElementById("popup").style.display = "block";
}

function closePopup() {
        document.getElementById("popup").style.display = "none";
}

const header = document.querySelector("header");

const buttons = [
        { text: "Connexion", url: "/connexion/connexion.html" },
        { text: "Inscription", url: "/inscription/inscription.html" },
        { text: "Profil", url: "/profil/profil.html" }
];

buttons.forEach(btn => {
        const button = document.createElement("button");
        button.className = "auth-button";
        button.textContent = btn.text;
        button.onclick = () => {
                window.location.href = btn.url;
        };
        header.appendChild(button);
});
