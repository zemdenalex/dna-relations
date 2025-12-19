# DNA Relations - Progress Tracker

## Project Status: Skeleton Complete

### Infrastructure
- [x] Docker Compose setup
- [x] PostgreSQL configuration
- [x] Nginx reverse proxy
- [ ] SSL certificates
- [ ] Domain configuration

### API (Go + tg)
- [x] Project structure
- [x] tg contracts defined
- [x] Types/models
- [x] Stub implementations
- [ ] Generate transport layer (`tg transport`)
- [ ] Auth service implementation
- [ ] Topic service implementation
- [ ] Event service implementation
- [ ] Note service implementation
- [ ] Category service implementation
- [ ] Favorite service implementation
- [ ] Media service implementation (P1)
- [ ] JWT middleware
- [ ] Database layer (pgx)

### Database
- [x] Users migration
- [x] Topics migration
- [x] Events migration
- [x] Notes migration
- [x] Favorites migration
- [x] Media sessions migration
- [ ] Run migrations
- [ ] Seed initial data (Denis, Nastya users)

### Telegram Bot
- [x] Project structure
- [x] Main menu
- [x] Command handlers (stubs)
- [x] Inline keyboards
- [x] Web App button
- [ ] API client integration
- [ ] Topic CRUD via bot
- [ ] Event viewing via bot
- [ ] Note quick-add
- [ ] Notifications

### Mini App (Telegram)
- [x] Project structure
- [x] Routing
- [x] Layout with bottom nav
- [x] Home page (stub)
- [x] Topics page (stub)
- [x] Calendar page (stub)
- [x] Notes page (stub)
- [x] API client
- [x] Telegram WebApp SDK integration
- [ ] Auth flow
- [ ] Topic CRUD
- [ ] Event CRUD
- [ ] Note CRUD

### Website
- [x] Project structure
- [x] Routing
- [x] Layout with sidebar
- [x] Home page (stub)
- [ ] All other pages
- [ ] Auth flow
- [ ] Full CRUD for all entities

## Next Steps

1. Initialize Go modules:
```bash
cd api && go mod tidy
cd ../bot && go mod tidy
```

2. Run API locally: `cd api && go run ./cmd/server`

3. Test health endpoint: `curl http://localhost:9000/health`

4. Implement AuthService (check password from .env, return JWT)

5. Implement TopicService with PostgreSQL

6. Test bot + miniapp flow

7. Deploy to VPS
