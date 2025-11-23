# Cloudzy Weather API

A RESTful weather API service built with Go that fetches and stores weather data from OpenWeatherMap API. The API provides endpoints to retrieve weather information and protected endpoints to create, update, and delete weather records.

## Installation

### Option 1: Local Development

1. **Clone the repository:**
   ```bash
   git clone https://github.com/zarinpy/abrnoc_weather
   cd abrnoc_weather
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database:**

    using psql
   ```bash
   psql -U postgres -c "CREATE DATABASE abrnoc_weather;"
   ```

4. **Create `.env` file:**
   ```bash
   cp .env.example .env  # If you have an example file
   # Or create manually with the variables listed above
   ```

5. **Run the application:**
   ```bash
   go run main.go
   ```

   The API will be available at `http://localhost:8080`

### Option 2: Docker Compose (Recommended)

1. **Update environment variables in `docker-compose.yml`:**
   ```yaml
   environment:
     - DB_HOST=db
     - OPENWEATHER_API_KEY=your_actual_api_key
     - JWT_SECRET=your_secure_secret
   ```

2. **Update PostgreSQL credentials in `docker-compose.yml`:**
   ```yaml
   environment:
     POSTGRES_DB: abrnoc_weather
     POSTGRES_USER: postgres
     POSTGRES_PASSWORD: yourpassword
   ```

3. **Build and run with Docker Compose:**
   ```bash
   docker-compose up --build
   ```

## Testing

The project includes comprehensive test coverage for authentication and weather endpoints.

### Prerequisites

Install test dependencies:
```bash
go get github.com/stretchr/testify/assert
go get gorm.io/driver/sqlite
```

### Running Tests

**Run all tests:**
```bash
go test ./internals/handlers/...
```
