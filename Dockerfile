# Stage 1: Build frontend assets with Node.js
FROM node:22-alpine AS frontend-builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY assets/ ./assets/
COPY public/ ./public/
COPY app/views/ ./app/views/
COPY vite.config.js ./
RUN npm run build

# Stage 2: Build Go application
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

RUN apk add curl
RUN mkdir -p bin && \
    curl -fsSL -o bin/dbmate https://github.com/amacneil/dbmate/releases/download/v2.26.0/dbmate-linux-amd64 && \
    chmod +x bin/dbmate

COPY go.mod go.sum ./
RUN go mod download

# Embedded files
COPY --from=frontend-builder /app/dist ./dist
COPY sql/pragmas.sql ./sql/
COPY locales ./locales/
# Go source files
COPY app ./app/
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/app .

# Stage 3: Runtime image
FROM alpine:latest

RUN apk add --no-cache sqlite

COPY --from=backend-builder /app/bin/dbmate /dbmate
COPY sql/migrations /migrations
COPY --from=backend-builder /app/bin/app /app
COPY deploy/run.sh /run.sh

EXPOSE 3000

CMD ["/run.sh"]
