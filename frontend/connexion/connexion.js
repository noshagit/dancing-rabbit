document.getElementById('signup-form').addEventListener('submit', async (event) => {
    event.preventDefault(); 

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    const response = await fetch('/connexion/connexion.html', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify({ email, password })
    });

    if (response.ok) {
        window.location.href = '/main-menu/menu.html';
    } else {
        const errorMessage = await response.text();
        alert("Erreur : " + errorMessage);
    }
});
