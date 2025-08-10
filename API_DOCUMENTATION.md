# LogScale API Documentation

## Endpoints

### Health Check
**GET** `/health`

Returns the health status of the API and database connectivity.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Create Logs
**POST** `/v1/logs`

Accepts either a single log or batch of logs.

#### Single Log
**Request:**
```json
{
  "log": {
    "service": "user-service",
    "level": "info",
    "message": "User login successful",
    "timestamp": "2024-01-15T10:30:00Z",
    "meta": {
      "user_id": "12345",
      "ip": "192.168.1.1"
    }
  }
}
```

#### Batch Logs
**Request:**
```json
{
  "logs": [
    {
      "service": "user-service",
      "level": "info",
      "message": "User login successful",
      "timestamp": "2024-01-15T10:30:00Z",
      "meta": {
        "user_id": "12345"
      }
    },
    {
      "service": "payment-service",
      "level": "error",
      "message": "Payment failed",
      "timestamp": "2024-01-15T10:31:00Z",
      "meta": {
        "transaction_id": "tx_123",
        "error_code": "INSUFFICIENT_FUNDS"
      }
    }
  ]
}
```

**Response:**
```json
{
  "id": 1,
  "service": "user-service",
  "level": "info",
  "message": "User login successful",
  "timestamp": "2024-01-15T10:30:00Z",
  "meta": {
    "user_id": "12345",
    "ip": "192.168.1.1"
  }
}
```

### Query Logs
**GET** `/v1/logs`

Query logs with filtering and pagination.

**Query Parameters:**
- `service` (optional): Filter by service name
- `level` (optional): Filter by log level
- `start_time` (optional): Filter logs from this time (ISO 8601 format)
- `end_time` (optional): Filter logs until this time (ISO 8601 format)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Number of logs per page (default: 50, max: 100)

**Example Request:**
```
GET /v1/logs?service=user-service&level=error&start_time=2024-01-15T00:00:00Z&page=1&limit=20
```

**Response:**
```json
{
  "logs": [
    {
      "id": 1,
      "service": "user-service",
      "level": "error",
      "message": "Database connection failed",
      "timestamp": "2024-01-15T10:30:00Z",
      "meta": {
        "error": "connection timeout"
      }
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 20
}
```

## Running the API

### Local Development
1. Copy `.env.example` to `.env` and update the database URL
2. Start PostgreSQL: `make up`
3. Run migrations: `make migrate-sql`
4. Start the API: `make dev`

### Using Docker Compose (Full Stack)
```bash
make up
```

### Development Commands
```bash
docker-compose up
```

The API will be available at `http://localhost:8080` 