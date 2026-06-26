#!/bin/bash

# Build and package script for Linux deployment
# Creates output folder with different environment packages

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
OUTPUT_DIR="$PROJECT_DIR/output"

APP_NAME="smileyan-backend"

echo "Building and packaging $APP_NAME..."

cd "$PROJECT_DIR"

# Clean and create output directory first
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Build for Linux amd64
echo "Building Linux amd64..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o "$OUTPUT_DIR/smileyan-backend" \
    main.go

# Create environment-specific packages
for env in dev prod; do
    ENV_DIR="$OUTPUT_DIR/$env"
    mkdir -p "$ENV_DIR"

    # Copy binary
    cp "$OUTPUT_DIR/smileyan-backend" "$ENV_DIR/"

    # Create config.yaml from template
    sed "s/mode: debug/mode: ${env}/" "$PROJECT_DIR/config.yaml" > "$ENV_DIR/config.yaml"

    # Create .env.example
    cat > "$ENV_DIR/.env.example" << 'ENVEOF'
# Database
SMILEYAN_BACKEND_DB_HOST=localhost
SMILEYAN_BACKEND_DB_USER=root
SMILEYAN_BACKEND_DB_PASSWORD=your_password
SMILEYAN_BACKEND_DB_NAME=smileyan

# Redis
SMILEYAN_BACKEND_REDIS_HOST=localhost
SMILEYAN_BACKEND_REDIS_PASSWORD=
SMILEYAN_BACKEND_REDIS_USERNAME=

# Email
SMILEYAN_BACKEND_EMAIL_PASSWORD=

# JWT
SMILEYAN_BACKEND_JWT_SECRET=your_secret_key
ENVEOF

    # Create run script
    cat > "$ENV_DIR/start.sh" << 'STRTEOF'
#!/bin/bash
cd "$(dirname "$0")"

# Load environment variables from .env file if it exists
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Start the backend
./smileyan-backend
STRTEOF
    chmod +x "$ENV_DIR/start.sh"

    # Create stop script
    cat > "$ENV_DIR/stop.sh" << 'STOPEOF'
#!/bin/bash
pkill -f smileyan-backend || true
STOPEOF
    chmod +x "$ENV_DIR/stop.sh"

    # Create restart script
    cat > "$ENV_DIR/restart.sh" << 'RSTREOF'
#!/bin/bash
cd "$(dirname "$0")"
./stop.sh
sleep 1
./start.sh &
RSTREOF
    chmod +x "$ENV_DIR/restart.sh"

    echo "Created package for: $env"
done

# Create tar.gz packages
cd "$OUTPUT_DIR"
tar --no-xattr -czvf "${APP_NAME}-dev.tar.gz" dev
tar --no-xattr -czvf "${APP_NAME}-prod.tar.gz" prod

# Cleanup intermediate files
rm -f smileyan-backend
rm -rf dev prod

# Create all-in-one package
mkdir -p all
cp *.tar.gz all/
tar --no-xattr -czvf "${APP_NAME}-all.tar.gz" all
rm -rf all

echo ""
echo "Build complete!"
echo "Output directory: $OUTPUT_DIR"
echo ""
ls -la "$OUTPUT_DIR"
echo ""
echo "Package contents:"
tar -tzvf "${APP_NAME}-dev.tar.gz" | head -10
echo "..."
echo ""
echo "To deploy:"
echo "  1. Copy package to server: scp ${APP_NAME}-dev.tar.gz user@server:/opt/"
echo "  2. Extract: tar -xzvf ${APP_NAME}-dev.tar.gz"
echo "  3. Edit dev/.env with real credentials"
echo "  4. Run: cd dev && ./start.sh"