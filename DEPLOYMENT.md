# Go-HTMX Deployment Guide

Deploy your Go-HTMX application with Let's Encrypt SSL. No API credentials required.

## 🚀 **Quick Start**

```bash
# One command deployment with SSL options
./scripts/deploy.sh
```

The script will ask you to choose:
1. **Let's Encrypt SSL** (free, requires domain)
2. **HTTP only** (development/testing)

## 📋 **Deployment Options**

### Option 1: Let's Encrypt SSL ✅ **Recommended for Production**

**Requirements:**
- Domain name pointing to your server
- Email address for certificate notifications
- Ports 80 and 443 open

**Setup:**
```bash
# 1. Run deployment script
./scripts/deploy.sh

# 2. Choose option 1 (Let's Encrypt)

# 3. Edit .env with your details:
DOMAIN=yourdomain.com
LETSENCRYPT_EMAIL=your-email@example.com
```

**Features:**
- ✅ Free SSL certificates
- ✅ Automatic renewal (every 60 days)
- ✅ Works with any domain provider
- ✅ Full security headers
- ✅ HTTP to HTTPS redirect

**Files used:**
- `docker-compose.yaml`
- `traefik/traefik.yml`

### Option 2: HTTP Only (No SSL)

**Perfect for:**
- Local development
- Testing
- Internal networks
- Before DNS setup

**Setup:**
```bash
# Run deployment script
./scripts/deploy.sh

# Choose option 2 (HTTP only)
```

**Access:**
- App: `http://localhost`
- Dashboard: `http://traefik.localhost:8080`

**Files used:**
- `docker-compose.no-ssl.yaml`
- `traefik/traefik.no-ssl.yml`

## 🔧 **Manual Deployment**

### Let's Encrypt Manual Setup

```bash
# 1. Create environment file
cp .env.example .env

# 2. Edit .env with your domain and email
nano .env

# 3. Update Traefik config email
sed -i 's/your-email@example.com/youremail@example.com/g' traefik/traefik.yml

# 4. Deploy
docker-compose up -d
```

### HTTP-Only Manual Setup

```bash
# Deploy without SSL
docker-compose -f docker-compose.no-ssl.yaml up -d
```

## 🔒 **SSL Certificate Process**

### Let's Encrypt HTTP Challenge

**How it works:**
1. Traefik requests certificate from Let's Encrypt
2. Let's Encrypt sends challenge to `http://yourdomain.com/.well-known/acme-challenge/`
3. Traefik responds with proof of domain ownership
4. Certificate issued and installed automatically

**Timeline:**
- First certificate: 2-5 minutes
- Renewal: Automatic every 60 days
- No manual intervention needed

**Requirements:**
- Domain DNS pointing to your server
- Port 80 accessible (for challenge)
- Port 443 open (for HTTPS traffic)

## 🌐 **Domain Setup**

### DNS Configuration

Point your domain to your server:
```
A record: yourdomain.com → YOUR_SERVER_IP
A record: traefik.yourdomain.com → YOUR_SERVER_IP (optional, for dashboard)
```

### Subdomain Support

Add multiple domains in docker-compose labels:
```yaml
- "traefik.http.routers.app.rule=Host(`yourdomain.com`) || Host(`www.yourdomain.com`)"
```

## 📊 **Comparison: SSL Options**

| Feature | Let's Encrypt | HTTP Only | Cloudflare DNS |
|---------|---------------|-----------|----------------|
| **SSL/HTTPS** | ✅ Free | ❌ No | ✅ Free |
| **Setup Complexity** | Easy | Very Easy | Moderate |
| **Requirements** | Domain only | None | CF account + API |
| **Production Ready** | ✅ Yes | ❌ No | ✅ Yes |
| **Auto Renewal** | ✅ Yes | N/A | ✅ Yes |
| **Wildcard Certs** | ❌ No | N/A | ✅ Yes |

## 🔍 **Troubleshooting**

### Let's Encrypt Issues

**Certificate not generating:**
```bash
# Check Traefik logs
docker-compose logs traefik

# Common issues:
# 1. Domain not pointing to server
nslookup yourdomain.com

# 2. Port 80 blocked
sudo ufw allow 80
sudo ufw allow 443

# 3. Previous certificate files
rm -rf traefik/data/acme.json
```

**Rate limiting (too many attempts):**
- Let's Encrypt limits: 5 failures per account/hour
- Wait 1 hour before retrying
- Use staging first: add `caServer: https://acme-staging-v02.api.letsencrypt.org/directory`

### HTTP-Only Issues

**Can't access application:**
```bash
# Check if containers are running
docker ps

# Check application logs
docker-compose -f docker-compose.no-ssl.yaml logs app

# Test direct connection
curl http://localhost:8080/health
```

### General Issues

**Port conflicts:**
```bash
# Check what's using ports
sudo netstat -tlnp | grep :80
sudo netstat -tlnp | grep :443

# Stop conflicting services
sudo systemctl stop apache2
sudo systemctl stop nginx
```

**Memory issues on 1GB VPS:**
```bash
# Check memory usage
free -h

# Restart with cleanup
docker-compose down
docker system prune -f
docker-compose up -d
```

## ⚡ **Performance Optimization**

### Resource Limits Applied

Both configurations include:
- **Memory limits**: App=250MB, Traefik=150MB
- **CPU limits**: Optimized for 1GB VPS
- **Connection limits**: Prevents resource exhaustion
- **Log rotation**: Prevents disk bloat

### Additional Optimizations

```bash
# Enable swap (recommended for 1GB VPS)
sudo fallocate -l 1G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# Make permanent
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

## 🔄 **Migration Between Options**

### From HTTP to Let's Encrypt

```bash
# Stop HTTP version
docker-compose -f docker-compose.no-ssl.yaml down

# Setup environment
cp .env.example .env
# Edit .env with your domain

# Start with SSL
docker-compose up -d
```

### From Let's Encrypt to Cloudflare

```bash
# Stop Let's Encrypt version
docker-compose -f docker-compose.letsencrypt.yaml down

# Backup existing certificates (optional)
cp -r traefik/data traefik/data.backup

# Setup Cloudflare credentials
cp traefik.env.example .env
# Add CF_API_EMAIL and CF_DNS_API_TOKEN

# Deploy with Cloudflare
docker-compose -f docker-compose.production.yaml up -d
```

## 📋 **Monitoring**

Use the included monitoring script:
```bash
./monitor.sh
```

**Key metrics:**
- Memory usage should stay under 800MB
- SSL certificate expiry (auto-renewed)
- Application response time
- Error rates in logs

## 🛡️ **Security Considerations**

### Let's Encrypt Security

- Certificates are industry-standard
- Automatic renewal prevents expiry
- Rate limiting prevents abuse
- Domain validation ensures ownership

### HTTP-Only Security

- **NOT for production**: No encryption
- OK for development/testing
- Use behind VPN if needed
- Consider reverse proxy with SSL termination

## 📦 **File Structure**

```
go-htmx/
├── docker-compose.yaml             # Main deployment (Let's Encrypt SSL)
├── docker-compose.no-ssl.yaml     # HTTP only alternative
├── .env.example                    # Environment template
├── scripts/deploy.sh               # Interactive deployment
├── traefik/
│   ├── traefik.yml                 # Main Traefik config (Let's Encrypt)
│   └── traefik.no-ssl.yml          # HTTP-only config
└── monitor.sh                      # Generated monitoring script
```

## ✅ **Success Checklist**

### Let's Encrypt Deployment
- [ ] Domain DNS points to server IP
- [ ] Ports 80 and 443 open
- [ ] `.env` file configured with domain and email
- [ ] Application loads at `https://yourdomain.com`
- [ ] SSL certificate shows as valid (green lock)
- [ ] Automatic HTTP→HTTPS redirect works

### HTTP-Only Deployment
- [ ] Application loads at `http://localhost`
- [ ] Traefik dashboard accessible
- [ ] All containers running without errors
- [ ] Memory usage under 800MB

**🎉 Your Go-HTMX application is now deployed with free SSL certificates using Let's Encrypt!** 