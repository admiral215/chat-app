
# Chat Application - Backend

A simple backend chat application built using Go with the following technologies:
- **Gin** - A fast web framework
- **GORM** - ORM for database management
- **Gorilla WebSocket** - WebSocket handling for real-time communication

This chat app is inspired by WhatsApp and includes essential features such as private messaging, group messaging, and user authentication.

## Features

- **Real-Time Communication**: Use WebSocket to send and receive messages in real-time.
- **Private Messaging**: Send messages between users privately.
- **Group Chat**: Users can join and send messages to groups.
- **User Authentication**: Secure user login using JWT (JSON Web Tokens).
- **Message Persistence**: Messages are saved in a database using GORM.
  
## Technologies Used

- **Go** (Golang)
- **Gin** (Web Framework)
- **GORM** (Database ORM)
- **Gorilla WebSocket** (WebSocket handling)
- **JWT** (JSON Web Tokens for authentication)
- **PostgreSQL** (Database)

## Prerequisites

Before you start, ensure you have the following installed:
- **Go** (version 1.22+)
- **PostgreSQL** (or another database compatible with GORM)
- **Go modules enabled**

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/chat-app.git
   cd chat-app
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up the database:
   - Make sure you have PostgreSQL running and create a database (or use your existing one).
   - Update the `DATABASE_URL` in the `.env` file with your database credentials.

4. Run the application:
   ```bash
   go run main.go
   ```

   The app will start running on `http://localhost:8080`.

5. Access the Swagger Documentation:
   Once your app is running, navigate to `http://localhost:8080/swagger/index.html#/` in your browser.

   You should now see the Swagger UI with your API documentation, including the ability to try out the API endpoints directly from the Swagger interface.

