document.getElementById("signup-form").addEventListener("submit", async function(event) {
    event.preventDefault();

    let pseudo = document.getElementById("pseudo").value;
    let email = document.getElementById("email").value;
    let password = document.getElementById("password").value;
    let confirmPassword = document.getElementById("confirm-password").value;

    if (password !== confirmPassword) {
        alert("Les mots de passe ne correspondent pas");
        return;
    }

    let userData = {
        pseudo: pseudo,
        email: email,
        password: password,
        confirmPassword: confirmPassword
    };

    try {
        let response = await fetch("/inscription/inscription.html", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(userData)
        });

        if (!response.ok) {
            let errorText = await response.text();
            throw new Error(errorText);
        } else {
            alert("Inscription réussie !");
            window.location.href = "/main-menu/menu.html";
        }

    } catch (error) {
        alert("Erreur : " + error.message);
    }
});
