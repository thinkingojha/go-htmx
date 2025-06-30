# Go-HTMX Application

A modern web application built with Go and HTMX, featuring Traefik reverse proxy with automatic SSL certificates.

## ✨ Features

- **Go Backend**: Fast, efficient server built with Go
- **HTMX Frontend**: Dynamic user interfaces without complex JavaScript
- **TailwindCSS**: Beautiful, responsive styling
- **Traefik Reverse Proxy**: Automatic SSL certificates with Let's Encrypt
- **Docker Deployment**: Containerized for easy deployment
- **Security Headers**: Production-ready security configuration
- **Rate Limiting**: Built-in request rate limiting
- **Health Checks**: Application monitoring and health endpoints

## 🚀 Quick Deployment

Deploy your application with SSL certificates in one command:

```bash
./scripts/deploy.sh
```

Choose from:
1. **Let's Encrypt SSL** (recommended for production)
2. **HTTP only** (development/testing)

## 📋 Requirements

- **Server**: 1GB RAM, 1 vCPU (DigitalOcean $6/month droplet works perfectly)
- **Domain**: Any domain provider (GoDaddy, Namecheap, etc.)
- **Ports**: 80, 443, and 22 (SSH)
- **Software**: Docker and Docker Compose

## 🔧 Local Development

### Setup
```bash
# Clone the repository
git clone <your-repo-url>
cd go-htmx

# Install dependencies
make deps

# Start development server
make dev
```

### Development Commands
```bash
# Build the application
make build

# Run tests
make test

# Lint code
make lint

# Build CSS
make build-css
```

## 🌐 Production Deployment

### 1. Server Setup
```bash
# On your server (Ubuntu/Debian)
sudo apt update && sudo apt upgrade -y
sudo apt install docker.io docker-compose-v2

# Add your user to docker group
sudo usermod -aG docker $USER
```

### 2. Deploy Application
```bash
# Clone and deploy
git clone <your-repo-url>
cd go-htmx

# Run deployment script
./scripts/deploy.sh

# Choose Let's Encrypt SSL
# Enter your domain and email
```

### 3. DNS Configuration
Point your domain to your server:
```
A record: yourdomain.com → YOUR_SERVER_IP
```

That's it! Your application will be available at `https://yourdomain.com`

## 📊 What You Get

### SSL & Security
- ✅ **Free SSL certificates** from Let's Encrypt
- ✅ **Automatic renewal** every 60 days
- ✅ **Security headers** (HSTS, CSP, XSS protection)
- ✅ **Rate limiting** (100 requests/minute)
- ✅ **Request validation** and filtering

### Performance
- ✅ **Optimized for 1GB VPS** (uses ~500MB)
- ✅ **Fast response times** (<300ms)
- ✅ **Efficient resource usage**
- ✅ **HTTP/2 support**

### Monitoring
- ✅ **Health checks** built-in
- ✅ **Access logs** and error tracking
- ✅ **Resource monitoring** script
- ✅ **Traefik dashboard** (optional)

## 🔍 Monitoring & Management

### Check Application Status
```bash
# Monitor resources
./monitor.sh

# View logs
docker-compose logs -f

# Check specific service
docker-compose logs app
docker-compose logs traefik
```

### Common Management Tasks
```bash
# Restart application
docker-compose restart app

# Update application
git pull && docker-compose up -d --build

# Stop services
docker-compose down

# View SSL certificate info
openssl s_client -connect yourdomain.com:443 -servername yourdomain.com
```

## 🛠️ Configuration

### Environment Variables (.env)
```bash
# Your domain
DOMAIN=yourdomain.com

# Email for SSL certificates
LETSENCRYPT_EMAIL=your-email@example.com

# Application settings
GOHTMX_APP_ENVIRONMENT=production
GOHTMX_SERVER_HOST=0.0.0.0
GOHTMX_SERVER_PORT=8080
```

### Application Configuration (config.yaml)
```yaml
server:
  host: "0.0.0.0"
  port: "8080"
  read_timeout: 30
  write_timeout: 30
  idle_timeout: 120
  shutdown_timeout: 30

app:
  environment: "production"
  title: "Go-HTMX App"
  debug: false
```

## 🔒 Security Features

- **Request Validation**: Only valid HTTP methods allowed
- **Rate Limiting**: 100 requests/minute per IP
- **Security Headers**: Complete security header suite
- **Non-root Containers**: All containers run as non-privileged users
- **Read-only Filesystem**: Application containers use read-only root filesystem
- **Network Isolation**: Services communicate through isolated Docker network

## 🚨 Troubleshooting

### SSL Certificate Issues
```bash
# Check certificate generation
docker-compose logs traefik | grep acme

# Verify domain DNS
nslookup yourdomain.com

# Check ports
sudo ufw allow 80
sudo ufw allow 443
```

### Application Issues
```bash
# Check application health
curl https://yourdomain.com/health

# View application logs
docker-compose logs app

# Check resource usage
free -h
docker stats --no-stream
```

### Common Solutions
- **Certificate not generating**: Verify domain points to server IP
- **Application not responding**: Check Docker containers are running
- **High memory usage**: Restart services with `docker-compose restart`

## 📈 Scaling Up

When you outgrow your 1GB server:

### Upgrade Server
- **2GB VPS**: Handle 200-500 concurrent users
- **4GB VPS**: Handle 500-1000+ concurrent users

### Add Services
```yaml
# Add Redis for caching
redis:
  image: redis:alpine
  restart: unless-stopped

# Add database
postgres:
  image: postgres:15-alpine
  restart: unless-stopped
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

- **Documentation**: [DEPLOYMENT.md](DEPLOYMENT.md)
- **Issues**: Create an issue on GitHub
- **Discussions**: Use GitHub Discussions for questions

---

**🎉 Your Go-HTMX application with automatic SSL certificates is ready for production!**
