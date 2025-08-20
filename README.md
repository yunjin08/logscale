# LogScale

A modern, scalable log aggregation and analytics platform built with Go, featuring real-time event processing, service metrics, and Redis Streams.

## Features

- 📊 **Real-time Analytics**: Automatic service metrics aggregation
- 🔄 **Event Streaming**: Redis Streams for reliable event processing
- 📈 **Service Monitoring**: Error rates, log counts, and performance metrics
- 🚀 **High Performance**: Built with Go for speed and efficiency
- 🐳 **Docker Ready**: Complete containerized deployment
- 🔍 **RESTful API**: Clean HTTP endpoints for log ingestion and querying

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for local development)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/yunjin08/logscale.git
   cd logscale
   ```

2. **Start the services**
   ```bash
   docker-compose up -d
   ```

3. **Run database migrations**
   ```bash
   make migrate-up
   ```

4. **Test the API**
   ```bash
   # Health check
   curl http://localhost:8080/health
   
   # Send a test log
   curl -X POST http://localhost:8080/v1/logs \
     -H "Content-Type: application/json" \
     -d '{"log": {"service": "test-service", "level": "info", "message": "Hello LogScale!"}}'
   ```

## Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   API       │    │   Worker    │    │   Redis     │
│  (HTTP)     │───▶│ (Processor) │───▶│  Streams    │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       ▼                   ▼                   ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ PostgreSQL  │    │ PostgreSQL  │    │ PostgreSQL  │
│   (Logs)    │    │ (Analytics) │    │ (Metrics)   │
└─────────────┘    └─────────────┘    └─────────────┘
```

## API Endpoints

### Logs
- `POST /v1/logs` - Create single or batch logs
- `GET /v1/logs` - Query logs with pagination and filters

### Health
- `GET /health` - Service health check

## Development

### Running Locally

```bash
# Start dependencies
docker-compose up postgres redis -d

# Run migrations
make migrate-up

# Start API
go run cmd/api/main.go

# Start Worker (in another terminal)
go run cmd/worker/main.go
```

### Available Make Commands

```bash
make up              # Start all services
make down            # Stop all services
make migrate-up      # Run database migrations
make migrate-down    # Rollback migrations
make test            # Run tests
make lint            # Run linter
```

### Environment Variables

```bash
DATABASE_URL=postgres://user:pass@localhost:5432/db?sslmode=disable
REDIS_URL=redis://localhost:6379
STREAM_NAME=logscale:logs
```

## Docker Deployment

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Scale worker instances
docker-compose up -d --scale worker=3
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.