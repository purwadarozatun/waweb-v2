#!/bin/bash
# Deploy waweb-v2 to 10.10.10.8 via SSH + Docker
# Usage: ./deploy.sh
# Requires: SSH access to theakistea@10.10.10.8

set -e

REMOTE_HOST="theakistea@10.10.10.8"
REMOTE_DIR="/home/theakistea/waweb-v2"
APP_PORT="8093"

echo "==> Packaging project..."
cd "$(dirname "$0")"

# Create tarball excluding unwanted files
tar --exclude='.git' \
    --exclude='published_sites' \
    --exclude='uploads' \
    --exclude='*.log' \
    --exclude='.env.deploy' \
    -czf /tmp/waweb-v2.tar.gz .

echo "==> Copying to remote server..."
scp /tmp/waweb-v2.tar.gz "${REMOTE_HOST}:/tmp/waweb-v2.tar.gz"

echo "==> Deploying on remote server..."
ssh "${REMOTE_HOST}" bash <<EOF
set -e

# Create deploy directory
mkdir -p ${REMOTE_DIR}
cd ${REMOTE_DIR}

# Extract
tar -xzf /tmp/waweb-v2.tar.gz
rm /tmp/waweb-v2.tar.gz

# Write .env
cat > .env <<'ENVEOF'
PORT=3000
JWT_SECRET=supersecret-change-in-production-waweb-v2-2024
HOST_FOLDER=./published_sites
TEMPLATES_FOLDER=./templates
DB_HOST=10.10.10.160
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourPassword
DB_NAME=wawebv2
ENVEOF

# Create required directories
mkdir -p uploads published_sites

# Stop & remove old container if exists
docker compose down 2>/dev/null || docker-compose down 2>/dev/null || true

# Build and start
if docker compose version &>/dev/null; then
    docker compose up -d --build
else
    docker-compose up -d --build
fi

echo ""
echo "==> Done! App running at http://10.10.10.8:${APP_PORT}"
docker ps --filter "name=waweb-v2" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
EOF

rm -f /tmp/waweb-v2.tar.gz

echo ""
echo "==> Deployment complete!"
echo "    URL: http://10.10.10.8:${APP_PORT}"
