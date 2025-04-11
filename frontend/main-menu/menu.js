function showPopup(title, text) {
    document.getElementById("popup-title").textContent = title;
    document.getElementById("popup-text").textContent = text;
    document.getElementById("popup").style.display = "block";
}

function closePopup() {
    document.getElementById("popup").style.display = "none";
}

