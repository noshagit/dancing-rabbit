document.getElementById('startButton').addEventListener('click', function() {

    const bunnyImage = document.createElement('img');
    bunnyImage.id = 'bunnyImage';
    bunnyImage.src = '/images/bunny.jpg';
    document.body.appendChild(bunnyImage);

    setTimeout(() => {
        bunnyImage.style.opacity = '1';
        bunnyImage.style.width = '2000px';
        bunnyImage.style.height = '2000px';
    }, 100);

    setTimeout(() => {
        window.location.href = '/main-menu/menu.html';
    }, 400);
});
