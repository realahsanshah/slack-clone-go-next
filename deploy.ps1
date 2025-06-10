param(
    [switch]$Prod,
    [switch]$Production,
    [switch]$Dev,
    [switch]$Development,
    [switch]$Cache,
    [switch]$Help
)

# Colors for output
$Red = "Red"
$Green = "Green"
$Yellow = "Yellow"

function Write-Status {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor $Green
}

function Write-Warning-Status {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor $Yellow
}

function Write-Error-Status {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $Red
}

# Show help
if ($Help) {
    Write-Host "Usage: .\deploy.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Prod, -Production      Deploy in production mode"
    Write-Host "  -Dev, -Development      Deploy in development mode (default)"
    Write-Host "  -Cache                  Use Docker build cache"
    Write-Host "  -Help                   Show this help message"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\deploy.ps1                    # Deploy in development mode"
    Write-Host "  .\deploy.ps1 -Prod              # Deploy in production mode"
    Write-Host "  .\deploy.ps1 -Dev -Cache        # Deploy in dev mode with cache"
    exit 0
}

# Check if Docker is running
try {
    docker info *>$null
}
catch {
    Write-Error-Status "Docker is not running. Please start Docker and try again."
    exit 1
}

# PORT=8080
# JWT_SECRET=jwt_secret
# Check if the .env file exists
if (-not (Test-Path -Path ".env")) {
    # Create a .env file with the default values
    Write-Status "Creating .env file..."
    New-Item -Path ".env" -ItemType File -Value "PORT=8080`nJWT_SECRET=jwt_secret"
    Write-Status ".env file created successfully."
    Write-Status "Please edit the .env file with the correct values."
}

# Determine mode
$Mode = "dev"
if ($Prod -or $Production) {
    $Mode = "prod"
}

# Determine build cache
$BuildCache = "--no-cache"
if ($Cache) {
    $BuildCache = ""
}

Write-Status "Starting deployment in $Mode mode..."

# Stop and remove existing containers
Write-Status "Stopping existing containers..."
docker compose down --remove-orphans

# Build and start services
if ($Mode -eq "prod") {
    Write-Status "Building and starting production services..."
    if ($BuildCache) {
        docker compose build $BuildCache api-prod
    } else {
        docker compose build api-prod
    }
    
    # Start all services - migrations will run automatically
    Write-Status "Starting services including database and migrations..."
    docker compose up -d
    
    Write-Status "Production deployment complete!"
    Write-Warning-Status "Don't forget to:"
    Write-Warning-Status "1. Change JWT_SECRET in docker-compose.yml"
    Write-Warning-Status "2. Use environment-specific database credentials"
    Write-Warning-Status "3. Set up proper SSL/TLS termination"
    Write-Warning-Status "4. Configure proper backup strategy"
} else {
    Write-Status "Building and starting development services..."
    if ($BuildCache) {
        docker compose build $BuildCache api-dev
    } else {
        docker compose build api-dev
    }
    
    # Start all services - migrations will run automatically
    Write-Status "Starting services including database and migrations..."
    docker compose up
    
    Write-Status "Development deployment complete!"
    Write-Status "API available at: http://localhost:8080"
    Write-Status "PostgreSQL available at: localhost:5432"
}

# Show running containers
Write-Status "Running containers:"
docker compose ps

# Show logs
Write-Status "Recent logs:"
docker compose logs --tail=20

Write-Status "Deployment finished! Use 'docker compose logs -f' to follow logs." 