# Building image
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main -ldflags="-s -w" -trimpath ./cmd/server-v4/


# Local development image
FROM golang:1.24.2-alpine AS development

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache curl
RUN go install github.com/air-verse/air@v1.52.3

EXPOSE 8080

CMD ["air", "-c", "./cmd/server-v4/.air.toml"]


# Production image
FROM gcr.io/distroless/static-debian12 AS production

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD [ "./main" ]
