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
        { text: "Connexion", url: "/connexion/connexion.html", id: "connexion-button" },
        { text: "Inscription", url: "/inscription/inscription.html", id: "inscription-button" },
        { text: "Profil", url: "/profil/profil.html", id: "profil-button" },
];

buttons.forEach(btn => {
        const button = document.createElement("button");
        button.className = "auth-button";
        button.textContent = btn.text;
        button.id = btn.id;
        button.onclick = () => {
                if (btn.id === "logout-button") {
                        logout();
                } else {
                        window.location.href = btn.url;
                }
        };
        header.appendChild(button);
});

const sessionCookie = document.cookie.split("; ").find(row => row.startsWith("session_token="));

if (sessionCookie) {
        document.getElementById("connexion-button").style.display = "none";
        document.getElementById("inscription-button").style.display = "none";
} else {
        document.getElementById("profil-button").style.display = "none";
}
