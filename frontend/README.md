# Dancing Rabbit

Dancing Rabbit is a multiplayer web-based game platform featuring various interactive games such as **Deaf Rhythm**, **Petit Bac**, and **Blind Test**. The platform allows players to compete in real-time, guess song lyrics, and enjoy a fun gaming experience.

## Features

- **Deaf Rhythm**: A game where players guess song titles based on displayed lyrics.
- **Petit Bac**: A classic word game with categories and timed rounds.
- **Blind Test**: A music guessing game where players identify songs from short audio clips.
- Real-time multiplayer functionality using WebSockets.
- User authentication and profile management.
- Responsive design for desktop and mobile devices.
- **Profile Management**: Users can check their profile.
- **Connection and Registration**: Secure user authentication with registration and login functionality.
- **Scoreboard**: A  game-specific scoreboard to track player rankings from the game you just finished.
- **Multiplayer Games**: Multiplayer experience with synchronized gameplay.

## Project Structure

### Backend

The backend is built with Go and uses the Gorilla Mux router for handling HTTP and WebSocket connections.

- **Main entry point**: [`backend/main/main.go`](backend/main/main.go)
- **Handlers**: Various handlers for user authentication, game menus, and game logic.

### Frontend

The frontend is built with HTML, CSS, and JavaScript.

## How to Run

### Prerequisites

- Go (latest version)
- A modern web browser

### Steps

1. **Clone the repository:**

   ```bash
   git clone https://github.com/noshagit/dancing-rabbit.git
   cd dancing-rabbit
   ```
2. **Launch the game:**

  At first, go to the main folder:

```bash
  cd backend/main/
```

  Then you can start the server:

```bash
  go run main.go
```

3. **Open the web page:**

*You have 3 options:*
  a. Type [http://localhost:8080](http://localhost:8080) in your browser (or click the link directly when the server is started).
  b. Upon launch, the console will display the message: "The server is running on port 8080: http://localhost:8080". Press CTRL + Click on "http://localhost:8080".
  c. If you're using VSCode, a small pop-up will appear in the bottom-right corner of your screen. Click "Open in Browser".

4. Test all our features and have fun.
