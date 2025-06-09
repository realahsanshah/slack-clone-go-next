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
docker-compose down --remove-orphans

# Function to run database migrations
function Run-Migrations {
    Write-Status "Running database migrations..."
    
    # Wait for database to be ready
    Start-Sleep -Seconds 5
    
    # Run migrations using goose from the API container
    try {
        if ($Mode -eq "prod") {
            docker-compose exec -T api-prod goose -dir sql/schema postgres "postgres://postgres:postgres123@postgres:5432/slackclone?sslmode=disable" up
        } else {
            docker-compose exec -T api-dev goose -dir sql/schema postgres "postgres://postgres:postgres123@postgres:5432/slackclone?sslmode=disable" up
        }
        Write-Status "Database migrations completed successfully!"
    }
    catch {
        Write-Warning-Status "Migration failed, database might already be up to date."
    }
}

# Build and start services
if ($Mode -eq "prod") {
    Write-Status "Building and starting production services..."
    if ($BuildCache) {
        docker-compose build $BuildCache api-prod
    } else {
        docker-compose build api-prod
    }
    docker-compose up -d postgres
    
    # Wait for postgres to be ready
    Write-Status "Waiting for PostgreSQL to be ready..."
    Start-Sleep -Seconds 10
    
    # Run migrations
    Run-Migrations
    
    # Start API service
    docker-compose up -d api-prod
    
    # Wait for services to be healthy
    Write-Status "Waiting for services to be healthy..."
    docker-compose ps
    
    Write-Status "Production deployment complete!"
    Write-Warning-Status "Don't forget to:"
    Write-Warning-Status "1. Change JWT_SECRET in docker-compose.yml"
    Write-Warning-Status "2. Use environment-specific database credentials"
    Write-Warning-Status "3. Set up proper SSL/TLS termination"
    Write-Warning-Status "4. Configure proper backup strategy"
} else {
    Write-Status "Building and starting development services..."
    if ($BuildCache) {
        docker-compose build $BuildCache api-dev
    } else {
        docker-compose build api-dev
    }
    docker-compose up -d postgres
    
    # Wait for postgres to be ready
    Write-Status "Waiting for PostgreSQL to be ready..."
    Start-Sleep -Seconds 10
    
    # Run migrations
    Run-Migrations
    
    # Start API service
    docker-compose up api-dev
    
    Write-Status "Development deployment complete!"
    Write-Status "API available at: http://localhost:8080"
    Write-Status "PostgreSQL available at: localhost:5432"
}

# Show running containers
Write-Status "Running containers:"
docker-compose ps

# Show logs
Write-Status "Recent logs:"
docker-compose logs --tail=20

Write-Status "Deployment finished! Use 'docker-compose logs -f' to follow logs." 