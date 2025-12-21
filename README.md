# DNA Relations Fix - Instructions

## Issues Fixed

1. **Database tables missing** - migrations now create all required tables
2. **Web auth broken** - token now stored in localStorage and sent with requests
3. **Bot Russian text** - all bot messages in Russian
4. **Russian language UI** - all web UI text translated to Russian
5. **Calendar UI** - improved weekly/monthly view with better UX

## Quick Deploy (on server)

```bash
cd /opt/dna-relations

# Download and extract fixes
# (upload dna-fix.tar.gz to server first)
tar -xzvf dna-fix.tar.gz

# Copy web files
cp -r dna-fix/web/src/* web/src/

# Copy bot files  
cp dna-fix/bot/go.mod bot/
cp -r dna-fix/bot/internal/* bot/internal/

# Run migrations
chmod +x dna-fix/deploy.sh
./dna-fix/deploy.sh

# Rebuild containers
docker compose rm -f web bot
docker compose up -d --build
```

## Manual Steps if Needed

### 1. Run Database Migrations
```bash
docker exec -i dna-postgres psql -U dna -d dna < dna-fix/deploy.sh
# Or run the SQL directly from deploy.sh
```

### 2. Fix Web Auth (client.ts)
The key change: store JWT token in localStorage and add to all requests:
```typescript
let authToken: string | null = localStorage.getItem('dna_token')

export function setToken(token: string) {
  authToken = token
  localStorage.setItem('dna_token', token)
}

// In request function:
if (authToken) {
  headers['Authorization'] = `Bearer ${authToken}`
}
```

### 3. Rebuild
```bash
docker compose rm -f web bot
docker compose up -d --build
```

## Files Changed

### Web
- `web/src/api/client.ts` - token storage + auth header
- `web/src/contexts/AuthContext.tsx` - auth state management
- `web/src/App.tsx` - router setup
- `web/src/pages/Login.tsx` - Russian login page
- `web/src/pages/Dashboard.tsx` - Russian dashboard
- `web/src/pages/Topics.tsx` - topics management
- `web/src/pages/Notes.tsx` - notes management
- `web/src/pages/Calendar.tsx` - improved calendar
- `web/src/index.css` - dark theme styles
- `web/src/main.tsx` - entry point

### Bot
- `bot/go.mod` - standard telegram-bot-api v5.5.1
- `bot/internal/keyboards/keyboards.go` - URL buttons + Russian text
- `bot/internal/handlers/handlers.go` - Russian messages
- `bot/internal/api/client.go` - API client

## Testing

1. Open https://dna-relations.site
2. Login with denis/3141
3. Test: create topic, add event, create note
4. Test bot: /start, /login, /topics

## Troubleshooting

**Still getting 401?**
- Clear browser localStorage
- Check API logs: `docker compose logs api`
- Verify token in browser DevTools > Application > localStorage

**500 errors?**
- Tables not created: run migrations
- Check: `docker exec dna-postgres psql -U dna -d dna -c "\dt"`

**Bot not building?**
- Ensure go.mod has standard library (no replace directive)
- Run: `docker compose rm -f bot && docker compose up -d --build bot`

## Note on Telegram Mini App

The bot uses URL buttons to open the web app (opens in browser). Native Telegram Mini App (WebAppInfo) requires a newer telegram-bot-api version than v5.5.1. For native Mini App support, you would need to:
1. Use `github.com/OvyFlash/telegram-bot-api` or wait for v6
2. Or configure a Main Mini App in @BotFather
