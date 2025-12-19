FROM golang:1.22-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod ./
RUN go mod download || true

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata postgresql-client

WORKDIR /app

COPY --from=builder /build/server .
COPY migrations ./migrations

EXPOSE 9000

CMD ["./server"]
