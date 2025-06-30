#!/bin/bash

# Simple deployment without Cloudflare API credentials
# Offers multiple SSL options

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}    Go-HTMX Deployment          ${NC}"
    echo -e "${BLUE}================================${NC}"
    echo ""
}

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check for required commands
check_requirements() {
    command -v docker >/dev/null 2>&1 || { print_error "Docker is required but not installed. Aborting." >&2; exit 1; }
    command -v docker-compose >/dev/null 2>&1 || { print_error "Docker Compose is required but not installed. Aborting." >&2; exit 1; }
}

# SSL options menu
show_ssl_options() {
    echo -e "${BLUE}Choose your SSL/HTTPS option:${NC}"
    echo ""
    echo "1. 🔒 Let's Encrypt (free SSL, requires domain)"
    echo "2. 🌐 HTTP only (no SSL, for development/testing)"
    echo "3. ❌ Cancel deployment"
    echo ""
}

deploy_letsencrypt() {
    print_status "Deploying with Let's Encrypt SSL..."
    
    # Create environment file
    if [ ! -f ".env" ]; then
        cp .env.example .env
        print_warning "Please edit .env file with your domain and email!"
        echo "   - Set DOMAIN to your actual domain name"
        echo "   - Set LETSENCRYPT_EMAIL to your email address"
        read -p "Press Enter when you've updated the .env file..."
    fi
    
    # Load environment variables
    if [ -f ".env" ]; then
        export $(cat .env | grep -v '#' | awk '/=/ {print $1}')
    fi
    
    # Validate domain
    if [ -z "$DOMAIN" ] || [ "$DOMAIN" == "yourdomain.com" ]; then
        print_error "Please set DOMAIN in .env file to your actual domain"
        exit 1
    fi
    
    # Validate email
    if [ -z "$LETSENCRYPT_EMAIL" ] || [ "$LETSENCRYPT_EMAIL" == "your-email@example.com" ]; then
        print_error "Please set LETSENCRYPT_EMAIL in .env file to your actual email"
        exit 1
    fi
    
    # Update Traefik config with user's email
    sed -i "s/your-email@example.com/$LETSENCRYPT_EMAIL/g" traefik/traefik.yml
    
    print_status "Domain: $DOMAIN"
    print_status "Email: $LETSENCRYPT_EMAIL"
    
    # Set up directories
    mkdir -p traefik/data
    chmod 600 traefik/data
    
    # Stop any existing containers
    docker-compose down --remove-orphans 2>/dev/null || true
    
    # Build and start
    print_status "Building and starting services..."
    docker-compose up -d --build
    
    print_status "✅ Deployment with Let's Encrypt completed!"
    echo ""
    echo "🌐 Your application will be available at:"
    echo "   https://$DOMAIN"
    echo "   https://traefik.$DOMAIN (Traefik dashboard)"
    echo ""
    print_warning "Note: SSL certificate generation may take 2-5 minutes on first run"
    print_warning "Make sure your domain's DNS points to this server's IP address"
}

deploy_http_only() {
    print_status "Deploying HTTP-only version..."
    
    # Create minimal environment file
    if [ ! -f ".env" ]; then
        echo "DOMAIN=localhost" > .env
        echo "GOHTMX_APP_ENVIRONMENT=production" >> .env
        echo "GOHTMX_SERVER_HOST=0.0.0.0" >> .env
        echo "GOHTMX_SERVER_PORT=8080" >> .env
    fi
    
    # Set up directories
    mkdir -p traefik/data
    chmod 600 traefik/data
    
    # Stop any existing containers
    docker-compose -f docker-compose.no-ssl.yaml down --remove-orphans 2>/dev/null || true
    
    # Build and start
    print_status "Building and starting services..."
    docker-compose -f docker-compose.no-ssl.yaml up -d --build
    
    print_status "✅ HTTP-only deployment completed!"
    echo ""
    echo "🌐 Your application is available at:"
    echo "   http://localhost"
    echo "   http://traefik.localhost:8080 (Traefik dashboard - user: admin, pass: admin)"
    echo ""
    print_warning "This is HTTP-only - not secure for production use"
}

# System optimization for 1GB VPS
optimize_system() {
    print_status "Applying basic system optimizations..."
    
    # Create swap if it doesn't exist and we have sudo
    if ! swapon --show | grep -q swap && command -v sudo >/dev/null 2>&1; then
        print_status "Creating swap file..."
        sudo fallocate -l 1G /swapfile 2>/dev/null || true
        sudo chmod 600 /swapfile 2>/dev/null || true
        sudo mkswap /swapfile 2>/dev/null || true
        sudo swapon /swapfile 2>/dev/null || true
    fi
}

# Monitoring setup
setup_monitoring() {
    print_status "Creating monitoring script..."
    cat << 'EOF' > monitor.sh
#!/bin/bash
echo "=== System Resources ==="
free -h
echo ""
echo "=== Docker Container Stats ==="
docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}"
echo ""
echo "=== Application Status ==="
curl -s http://localhost/health && echo " ✅ App healthy" || echo " ❌ App not responding"
EOF
    chmod +x monitor.sh
    print_status "Created monitoring script: ./monitor.sh"
}

# Main deployment flow
main() {
    print_header
    
    # Check if running as root
    if [[ $EUID -eq 0 ]]; then
        print_error "This script should not be run as root for security reasons"
        exit 1
    fi
    
    check_requirements
    optimize_system
    
    # Show SSL options
    while true; do
        show_ssl_options
        read -p "Enter your choice (1-3): " choice
        
        case $choice in
            1)
                deploy_letsencrypt
                break
                ;;
            2)
                deploy_http_only
                break
                ;;
            3)
                print_status "Deployment cancelled."
                exit 0
                ;;
            *)
                print_error "Invalid option. Please choose 1, 2, or 3."
                ;;
        esac
    done
    
    setup_monitoring
    
    echo ""
    echo "📋 Useful commands:"
    echo "   Monitor resources: ./monitor.sh"
    echo "   View logs:        docker-compose logs -f"
    echo "   Stop services:    docker-compose down"
    echo "   Restart app:      docker-compose restart app"
    echo ""
    print_status "🚀 Deployment completed successfully!"
}

main "$@" 