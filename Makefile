.PHONY: dev up down logs build generate migrate

dev:
	docker-compose up --build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

build:
	docker-compose build

generate:
	cd api && tg transport --services ./contracts --out ./internal/transport
	cd api && tg client -js --services ./contracts --outPath ../miniapp/src/lib/generated

migrate:
	docker exec dna-api ./server migrate

api-dev:
	cd api && go run ./cmd/server

bot-dev:
	cd bot && go run ./cmd

web-dev:
	cd web && npm run dev

miniapp-dev:
	cd miniapp && npm run dev
