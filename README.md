# Cloudzy Weather API

A RESTful weather API service built with Go that fetches and stores weather data from OpenWeatherMap API. The API provides endpoints to retrieve weather information and protected endpoints to create, update, and delete weather records.

## Installation

### Option 1: Local Development

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd abrnoc_weather
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database:**
   ```bash
   # Create database
   createdb abrnoc_weather
   
   # Or using psql
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

   The API will be available at `http://localhost:8080`

### Option 3: Docker (Standalone)

1. **Build the Docker image:**
   ```bash
   docker build -t abrnoc-weather .
   ```

2. **Run the container:**
   ```bash
   docker run -p 8080:8080 \
     -e DB_HOST=host.docker.internal \
     -e DB_USER=postgres \
     -e DB_PASSWORD=yourpassword \
     -e DB_NAME=abrnoc_weather \
     -e DB_PORT=5432 \
     -e OPENWEATHER_API_KEY=your_api_key \
     -e JWT_SECRET=your_secret \
     abrnoc-weather
   ```

## API Endpoints

### Public Endpoints

- `GET /weather` - Get all weather records
- `GET /weather/:id` - Get weather record by ID
- `GET /weather/latest/:cityName` - Get latest weather for a specific city
- `GET /swagger/*any` - Swagger API documentation

### Protected Endpoints (Require JWT Bearer Token)

- `POST /weather` - Create a new weather record
  ```json
  {
    "cityName": "London",
    "country": "GB"
  }
  ```

- `PUT /weather/:id` - Update a weather record
- `DELETE /weather/:id` - Delete a weather record

## Authentication

Protected endpoints require a JWT Bearer token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

To generate a JWT token, you'll need to use a JWT library or tool. The token must be signed with the `JWT_SECRET` from your environment variables.

## API Documentation

Once the server is running, access the Swagger documentation at:
- `http://localhost:8080/swagger/index.html`

### Generating Swagger Documentation

```bash
swag init
```
