# ServeGo
ServeGo is a minimal HTTP web server built with Go, designed to demonstrate a production-grade setup. It includes essential features such as graceful shutdown, structured logging, and health check endpoints.

## Features

- **Minimal HTTP Server**: Lightweight and efficient server implementation.
- **Health Check Endpoints**:
  - `/health/check`: Reports the health status of the instance.
  - `/health/alive`: Confirms the instance is alive.
- **Root Endpoint**:
  - `/root`: A welcome message for the server.
- **Graceful Shutdown**: Ensures proper cleanup of resources during termination.
- **Structured Logging**: Uses `slog` for JSON-based structured logging.

## Getting Started

### Prerequisites

- Go 1.22 or higher

### Running the Server

1. Clone the repository:
   ```sh
   git clone https://github.com/Yajanth/ServeGo.git
   cd ServeGo

2. Run the server:
   ```sh
   go run main.go
   