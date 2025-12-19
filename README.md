# DNA Relations

A relationship quality-of-life system for couples. Manage discussion topics, shared calendar, notes, and more.

## Features

- **Topics**: Track things to discuss with priority levels (0-5)
- **Calendar**: Shared events with date/time management
- **Notes**: Rules, thoughts, resources, credentials
- **Multi-platform**: Telegram bot, Mini App, Web interface

## Quick Start (Local Development)

### Prerequisites
- Go 1.22+
- Node.js 20+
- Docker & Docker Compose
- PostgreSQL (or use Docker)

### 1. Clone and setup
```bash
git clone https://github.com/zemdenalex/dna-relations.git
cd dna-relations
cp .env.example .env
# Edit .env with your values
```

### 2. Start database
```bash
docker-compose up -d postgres
```

### 3. Run migrations
```bash
# Connect to postgres and run migrations
docker exec -i dna-postgres psql -U dna -d dna < api/migrations/001_users.sql
docker exec -i dna-postgres psql -U dna -d dna < api/migrations/002_topics.sql
docker exec -i dna-postgres psql -U dna -d dna < api/migrations/003_events.sql
docker exec -i dna-postgres psql -U dna -d dna < api/migrations/004_notes.sql
docker exec -i dna-postgres psql -U dna -d dna < api/migrations/005_favorites.sql
docker exec -i dna-postgres psql -U dna -d dna < api/migrations/006_media_sessions.sql
docker exec -i dna-postgres psql -U dna -d dna < api/migrations/007_seed_data.sql
```

### 4. Start API
```bash
cd api
go mod tidy
go run ./cmd/server
# API runs on http://localhost:9000
```

### 5. Start Bot
```bash
cd bot
go mod tidy
go run ./cmd
```

### 6. Start Mini App (for development)
```bash
cd miniapp
npm install
npm run dev
# Runs on http://localhost:5173
```

## Production Deployment (Ubuntu VPS)

### 1. Copy files to server
```bash
scp -r dna-relations user@your-server:/opt/
```

### 2. Run deployment script
```bash
ssh user@your-server
cd /opt/dna-relations
sudo bash scripts/deploy.sh your-domain.com
```

### 3. Configure
```bash
sudo nano /opt/dna-relations/.env
# Set all required values
docker-compose restart
```

### 4. Setup Telegram Bot
1. Create bot with @BotFather
2. Set Web App URL: `/setmenubutton` → your-domain.com/miniapp
3. Add token to .env

## API Endpoints

### Auth
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/auth/me` - Get current user

### Topics
- `GET /api/v1/topics` - List topics
- `POST /api/v1/topics` - Create topic
- `GET /api/v1/topics/:id` - Get topic
- `PUT /api/v1/topics/:id` - Update topic
- `DELETE /api/v1/topics/:id` - Delete topic
- `POST /api/v1/topics/:id/discuss` - Mark as discussed

### Events
- `GET /api/v1/events?from=&to=` - List events
- `POST /api/v1/events` - Create event
- `GET /api/v1/events/:id` - Get event
- `PUT /api/v1/events/:id` - Update event
- `DELETE /api/v1/events/:id` - Delete event

### Notes
- `GET /api/v1/notes` - List notes
- `POST /api/v1/notes` - Create note
- `GET /api/v1/notes/:id` - Get note
- `PUT /api/v1/notes/:id` - Update note
- `DELETE /api/v1/notes/:id` - Delete note
- `POST /api/v1/notes/:id/pin` - Toggle pin

## Priority System

| Priority | Name | Use Case |
|----------|------|----------|
| 0 | Buffer | Low priority parking |
| 1 | Emergency | Must discuss NOW |
| 2 | ASAP | Very urgent |
| 3 | Must discuss | Important |
| 4 | Soon | Can wait a bit |
| 5 | When possible | No rush |

## Note Types

- **general** - Regular notes
- **rule** - Relationship agreements
- **thought** - Important realizations
- **resource** - Links, books, media
- **credential** - Shared logins

## Backups

Automatic daily backups at 3 AM (configured during deployment).

Manual backup:
```bash
/usr/local/bin/dna-backup
```

Backups stored in `/var/backups/dna-relations/`

## Tech Stack

- **Backend**: Go, Chi router, PostgreSQL
- **Frontend**: React, TypeScript, Tailwind CSS
- **Bot**: telegram-bot-api v5
- **Infrastructure**: Docker, Nginx, Let's Encrypt
