#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if the .env file exists
if [ ! -f ".env" ]; then
    print_status "Creating .env file..."
    echo "PORT=8080" > .env
    echo "JWT_SECRET=jwt_secret" >> .env
    print_status ".env file created successfully."
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    read -n 1 -s -r -p "Press any key to continue..."
    echo ""
    exit 1
fi

# Default values
MODE="dev"
BUILD_CACHE="--no-cache"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --prod|--production)
            MODE="prod"
            shift
            ;;
        --dev|--development)
            MODE="dev"
            shift
            ;;
        --cache)
            BUILD_CACHE=""
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --prod, --production    Deploy in production mode"
            echo "  --dev, --development    Deploy in development mode (default)"
            echo "  --cache                 Use Docker build cache"
            echo "  --help, -h             Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                     # Deploy in development mode"
            echo "  $0 --prod              # Deploy in production mode"
            echo "  $0 --dev --cache       # Deploy in dev mode with cache"
            echo ""
            echo "Press any key to continue..."
            read -n 1 -s -r -p ""
            echo ""
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use --help for usage information."
            echo ""
            echo "Press any key to continue..."
            read -n 1 -s -r -p ""
            echo ""
            exit 1
            ;;
    esac
done

print_status "Starting deployment in $MODE mode..."

# Stop and remove existing containers
print_status "Stopping existing containers..."
docker-compose down --remove-orphans

# Build and start services
if [ "$MODE" = "prod" ]; then
    print_status "Building and starting production services..."
    docker-compose build $BUILD_CACHE api-prod
    docker-compose up -d postgres api-prod
    
    # Wait for services to be healthy
    print_status "Waiting for services to be healthy..."
    docker-compose ps
    
    print_status "Production deployment complete!"
    print_warning "Don't forget to:"
    print_warning "1. Change JWT_SECRET in docker-compose.yml"
    print_warning "2. Use environment-specific database credentials"
    print_warning "3. Set up proper SSL/TLS termination"
    print_warning "4. Configure proper backup strategy"
    
else
    print_status "Building and starting development services..."
    docker-compose build $BUILD_CACHE api-dev
    docker-compose up postgres api-dev
    
    print_status "Development deployment complete!"
    print_status "API available at: http://localhost:8080"
    print_status "PostgreSQL available at: localhost:5432"
fi

# Show running containers
print_status "Running containers:"
docker-compose ps

# Show logs
print_status "Recent logs:"
docker-compose logs --tail=20

print_status "Deployment finished! Use 'docker-compose logs -f' to follow logs." 