function getCookie(name) {
    let value = "; " + document.cookie;
    let parts = value.split("; " + name + "=");
    if (parts.length === 2) {
        return parts.pop().split(";").shift();
    }
    return null;
}

function loadProfile() {
    const sessionToken = getCookie("session_token");

    if (sessionToken) {
        fetch("/api/get-profile", {
            method: "GET",
            credentials: "include"
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                document.getElementById("pseudo").textContent = data.profile.pseudo;
                document.getElementById("email").textContent = data.profile.email;
            } else {
                alert("Erreur lors de la récupération des informations du profil.");
            }
        })

        .catch(error => {
            console.error("Erreur:", error);
            alert("Erreur lors de la récupération des données.");
        });
    } else {
        alert("Vous n'êtes pas connecté.");
        window.location.href = "/connexion/connexion.html";
    }
}

function logout() {
    fetch("/logout", {
        method: "POST"
    })
    .then(() => {
        alert("Déconnexion réussie !");
        window.location.href = "/main-menu/menu.html";
    })
    .catch(error => {
        console.error("Erreur de déconnexion:", error);
        alert("Erreur lors de la déconnexion.");
    });
}

document.addEventListener("DOMContentLoaded", loadProfile);