FROM golang:1.22-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod ./
RUN go mod download || true

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/bot .

CMD ["./bot"]
