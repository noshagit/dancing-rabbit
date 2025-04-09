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
        { text: "Déconnexion", url: "/main-menu/menu.html", id: "logout-button" },
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
        document.getElementById("logout-button").style.display = "none";
}

function logout() {
        fetch("/logout", {
                method: "POST",
                credentials: "include",
        })
        .then(response => {
                if (response.ok) {
                        document.cookie = "session_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
                        window.location.href = "/main-menu/menu.html";
                } else {
                        showPopup("Erreur", "Erreur lors de la déconnexion.");
                }
        })
        .catch(error => {
                console.error("Erreur lors de la déconnexion:", error);
                showPopup("Erreur", "Erreur lors de la déconnexion.");
        });
}