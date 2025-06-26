# Go-HTMX Production-Ready Application

A production-ready personal portfolio and blog application built with Go, HTMX, and TailwindCSS. This application demonstrates modern web development practices with server-side rendering, progressive enhancement, and comprehensive production features.

## 🚀 Features

### Frontend
- **HTMX Integration** - Dynamic user interactions without JavaScript complexity
- **TailwindCSS** - Utility-first styling with optimized builds
- **Progressive Enhancement** - Works without JavaScript, enhanced with HTMX
- **Responsive Design** - Mobile-first responsive layouts
- **SEO Optimized** - Proper meta tags and semantic HTML

### Backend
- **Structured Logging** - JSON logging in production, colorized in development
- **Configuration Management** - YAML-based config with environment variable overrides
- **Graceful Shutdown** - Proper server shutdown handling
- **Health Checks** - Built-in health monitoring endpoint
- **Security Headers** - Comprehensive security headers middleware
- **Rate Limiting** - Configurable rate limiting to prevent abuse
- **Error Handling** - Structured error responses with proper HTTP status codes
- **CORS Support** - Configurable CORS middleware
- **Request Timeout** - Prevents hanging requests
- **Recovery Middleware** - Panic recovery with logging

### DevOps & Production
- **Docker Support** - Multi-stage builds with security best practices
- **Configuration Profiles** - Separate configs for development and production
- **Comprehensive Testing** - Unit tests with coverage reporting
- **Code Quality** - Linting, security checks, and formatting
- **CI/CD Ready** - Makefile with all necessary targets
- **Documentation** - Comprehensive setup and deployment guides

## 📁 Project Structure

```
go-htmx/
├── cmd/
│   └── server/           # Server implementation
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP handlers
│   ├── logger/          # Structured logging
│   ├── middleware/      # HTTP middleware
│   ├── static/          # Static assets
│   ├── template/        # HTML templates
│   └── utils/           # Utility functions
├── configs/             # Configuration files
├── docker-compose.yaml  # Docker composition
├── Dockerfile          # Production Docker build
├── Makefile           # Build and deployment targets
└── README.md
```

## 🛠 Quick Start

### Prerequisites
- Go 1.21 or later
- Node.js 18+ (for TailwindCSS)
- Docker (optional)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-htmx
   ```

2. **Install dependencies**
   ```bash
   make deps
   ```

3. **Run in development mode**
   ```bash
   make dev
   ```

   The application will start at `http://localhost:3000`

### Production Deployment

1. **Build the application**
   ```bash
   make build
   ```

2. **Run with Docker**
   ```bash
   make docker-run
   ```

   Or using docker-compose:
   ```bash
   docker-compose up app
   ```

## ⚙️ Configuration

The application uses a hierarchical configuration system:

1. **Default values** (in code)
2. **Configuration files** (`config.yaml`, `config.dev.yaml`)
3. **Environment variables** (prefixed with `GOHTMX_`)

### Environment Variables

```bash
# Server Configuration
GOHTMX_SERVER_PORT=8080
GOHTMX_SERVER_HOST=0.0.0.0
GOHTMX_SERVER_READ_TIMEOUT=10
GOHTMX_SERVER_WRITE_TIMEOUT=10

# Application Configuration
GOHTMX_APP_ENVIRONMENT=production
GOHTMX_APP_LOG_LEVEL=info

# Security Configuration
GOHTMX_SECURITY_RATE_LIMIT_RPM=100
```

See `.env.example` for a complete list of available environment variables.

## 🧪 Testing

Run the test suite:

```bash
# Run all tests with coverage
make test

# Run short tests only
make test-short

# Run benchmarks
make bench
```

## 🔍 Code Quality

This project includes comprehensive code quality tools:

```bash
# Format code and tidy modules
make tidy

# Run linters
make lint

# Run security checks
make security-check

# Run full CI pipeline
make ci
```

## 🔧 Available Make Targets

- `make dev` - Run in development mode
- `make build` - Build production binary
- `make test` - Run tests with coverage
- `make docker-build` - Build Docker image
- `make docker-run` - Run with Docker
- `make health-check` - Check application health
- `make ci` - Run full CI pipeline
- `make help` - Show all available targets

## 🏗 Architecture

### Server Architecture
- **Graceful Shutdown** - Handles SIGINT/SIGTERM signals
- **Middleware Pipeline** - Recovery → Security → Logging → Rate Limiting → CORS → Timeout
- **Error Handling** - Structured error responses for both HTML and JSON clients
- **Health Monitoring** - `/health` endpoint for load balancers and monitoring

### HTMX Integration
- **Progressive Enhancement** - Application works without JavaScript
- **Dynamic Loading** - Live markdown preview with HTMX
- **Error Handling** - Global HTMX error handling with user feedback
- **Loading States** - Visual feedback during requests

### Security Features
- **Security Headers** - X-Content-Type-Options, X-Frame-Options, etc.
- **Rate Limiting** - Configurable requests per minute
- **Input Validation** - Proper form validation and sanitization
- **CSRF Protection** - Built-in CSRF protection (can be enabled)

## 📊 Monitoring

### Health Check
```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Logs
- **Development** - Colorized console output
- **Production** - Structured JSON logs
- **Request Logging** - All HTTP requests logged with timing

## 🚀 Deployment

### Docker Deployment
```bash
# Build and run
docker-compose up app

# Scale to multiple instances
docker-compose up --scale app=3
```

### Manual Deployment
```bash
# Build for production
make build

# Copy binary and assets to server
scp bin/gohtmx user@server:/opt/gohtmx/
scp -r internal/static user@server:/opt/gohtmx/
scp -r internal/template user@server:/opt/gohtmx/
scp config.yaml user@server:/opt/gohtmx/

# Run on server
GOHTMX_APP_ENVIRONMENT=production ./gohtmx
```

## 🔒 Security

- **Non-root containers** - Docker containers run as non-root user
- **Read-only filesystem** - Production containers use read-only root filesystem
- **Security scanning** - `gosec` integration for security vulnerability scanning
- **Dependency updates** - Regular dependency updates with `go mod tidy`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run the full test suite: `make ci`
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [HTMX](https://htmx.org/) - For making HTML great again
- [TailwindCSS](https://tailwindcss.com/) - For utility-first styling
- [Gorilla Mux](https://github.com/gorilla/mux) - For powerful routing
- [Logrus](https://github.com/sirupsen/logrus) - For structured logging
