document.getElementById("signup-form").addEventListener("submit", function(event) {
    event.preventDefault();
    
    let password = document.getElementById("password").value;
    let confirmPassword = document.getElementById("confirm-password").value;
    
    if (password !== confirmPassword) {
        alert("Les mots de passe ne correspondent pas");
        return;
    }

    alert("Inscription réussie bien joué la team");
});
