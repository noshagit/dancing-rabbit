# Dancing Rabbit

Dancing Rabbit is a multiplayer web-based game platform featuring various interactive games such as **Deaf Rhythm**, **Petit Bac**, and **Blind Test**. The platform allows players to compete in real-time, guess song lyrics, and enjoy a fun gaming experience.

## Features

- **Deaf Rhythm**: A game where players guess song titles based on displayed lyrics.
- **Petit Bac**: A classic word game with categories and timed rounds.
- **Blind Test**: A music guessing game where players identify songs from short audio clips.
- Real-time multiplayer functionality using WebSockets.
- User authentication and profile management.
- Responsive design for desktop and mobile devices.

## Project Structure

### Backend

The backend is built with Go and uses the Gorilla Mux router for handling HTTP and WebSocket connections.

- **Main entry point**: [`backend/main/main.go`](backend/main/main.go)
- **Game logic**: [`backend/games/deaf_rhythm.go`](backend/games/deaf_rhythm.go)
- **Handlers**: Various handlers for user authentication, game menus, and game logic.

### Frontend

The frontend is built with HTML, CSS, and JavaScript.

- **Deaf Rhythm Game**:

  - HTML: [`frontend/deaf-rhythm/game/deaf-rhythm.html`](frontend/deaf-rhythm/game/deaf-rhythm.html)
  - JavaScript: [`frontend/deaf-rhythm/game/deaf-rhythm.js`](frontend/deaf-rhythm/game/deaf-rhythm.js)
  - CSS: `frontend/deaf-rhythm/game/deaf-rhythm.css`
- **Main Menu**:

  - HTML: `frontend/main-menu/menu.html`
  - JavaScript: `frontend/main-menu/menu.js`
  - CSS: `frontend/main-menu/menu.css`

## How to Run

### Prerequisites

- Go (latest version)
- Node.js (if additional frontend tooling is required)
- A modern web browser

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/dancing-rabbit.git
   cd dancing-rabbit/backend/main
   ```

2. Launch the game 
    ```bash
    go run main.go
    ```