# Doorman

Authentication and authorization service.

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/zollidan/doorman.git
cd doorman
```

2. Create your configuration file:
```bash
cp .env.example .env
```

3. Edit `.env` and configure your settings (especially change `JWT_SECRET` and database password)

4. Start the services:
```bash
docker-compose up -d
```

The service will be available at `http://localhost:2222`

### Using Docker Directly

1. Build the image:
```bash
docker build -t doorman:latest .
```

2. Run with PostgreSQL (create `.env` file first):
```bash
docker run -d \
  --name doorman \
  --env-file .env \
  -p 2222:2222 \
  doorman:latest
```

### Using Pre-built Image

```bash
# Pull the image
docker pull ghcr.io/zollidan/doorman:latest

# Run with your config
docker run -d \
  --name doorman \
  --env-file .env \
  -p 2222:2222 \
  ghcr.io/zollidan/doorman:latest
```

## Configuration

Create a `.env` file based on `.env.example`:

- `SERVER_ADDRESS` - Server port (default: `:2222`)
- `APP_MODE` - `development` or `production`
- `POSTGRES_*` - PostgreSQL connection settings
- `JWT_SECRET` - **MUST be changed in production!**
- `JWT_EXPIRY` - Access token lifetime
- `REFRESH_TOKEN_EXPIRY` - Refresh token lifetime

## Database Support

- **PostgreSQL** (recommended for production)
- **SQLite** (development mode only)

## Development Mode

When running without a `.env` file, Doorman starts with SQLite database in development mode. **Do not use in production without proper configuration.**

## Health Check

The service provides a health check endpoint:
```bash
curl http://localhost:2222/health
```

## License

MIT
