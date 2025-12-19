#!/bin/bash

set -e

echo "DNA Relations - Deployment Script"
echo "=================================="

if [ "$EUID" -ne 0 ]; then
    echo "Please run as root (sudo)"
    exit 1
fi

INSTALL_DIR="/opt/dna-relations"
DOMAIN="${1:-dna-relations.website}"

echo "Installing to: $INSTALL_DIR"
echo "Domain: $DOMAIN"

apt-get update
apt-get install -y docker.io docker-compose certbot

systemctl enable docker
systemctl start docker

mkdir -p $INSTALL_DIR
cd $INSTALL_DIR

if [ ! -f ".env" ]; then
    echo "Creating .env file..."
    cat > .env << EOF
POSTGRES_USER=dna
POSTGRES_PASSWORD=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 32)
POSTGRES_DB=dna

AUTH_PASSWORD_DENIS=change_me
AUTH_PASSWORD_NASTYA=change_me
JWT_SECRET=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 32)

TELEGRAM_BOT_TOKEN=your_bot_token_here

WEBAPP_URL=https://$DOMAIN/miniapp

TZ=Europe/Moscow
EOF
    echo "IMPORTANT: Edit .env and set your passwords and bot token!"
    echo "  nano $INSTALL_DIR/.env"
fi

echo "Stopping existing services..."
docker-compose down 2>/dev/null || true

echo "Obtaining SSL certificate..."
certbot certonly --standalone -d $DOMAIN --non-interactive --agree-tos --email admin@$DOMAIN || true

mkdir -p nginx/ssl
if [ -f "/etc/letsencrypt/live/$DOMAIN/fullchain.pem" ]; then
    cp /etc/letsencrypt/live/$DOMAIN/fullchain.pem nginx/ssl/cert.pem
    cp /etc/letsencrypt/live/$DOMAIN/privkey.pem nginx/ssl/key.pem
else
    echo "SSL certificate not found. Generating self-signed..."
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout nginx/ssl/key.pem \
        -out nginx/ssl/cert.pem \
        -subj "/CN=$DOMAIN"
fi

echo "Building and starting services..."
docker-compose build
docker-compose up -d

echo "Waiting for database..."
sleep 10

echo "Running migrations..."
docker exec dna-api sh -c 'for f in /app/migrations/*.sql; do
    echo "Running $f..."
    cat "$f" | grep -v "^--" | PGPASSWORD=$POSTGRES_PASSWORD psql -h postgres -U $POSTGRES_USER -d $POSTGRES_DB 2>/dev/null || true
done'

echo "Setting up backup cron job..."
cp scripts/backup.sh /usr/local/bin/dna-backup
chmod +x /usr/local/bin/dna-backup
echo "0 3 * * * /usr/local/bin/dna-backup" | crontab -

cat > /etc/systemd/system/dna-relations.service << EOF
[Unit]
Description=DNA Relations
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=$INSTALL_DIR
ExecStart=/usr/bin/docker-compose up -d
ExecStop=/usr/bin/docker-compose down

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable dna-relations

echo ""
echo "=================================="
echo "Deployment complete!"
echo ""
echo "Services:"
echo "  - API: https://$DOMAIN/api/v1"
echo "  - Web: https://$DOMAIN"
echo "  - Mini App: https://$DOMAIN/miniapp"
echo ""
echo "Next steps:"
echo "  1. Edit .env: nano $INSTALL_DIR/.env"
echo "  2. Restart: docker-compose restart"
echo "  3. Check logs: docker-compose logs -f"
echo ""
