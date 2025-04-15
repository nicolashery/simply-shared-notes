# Stage 1: Build frontend assets with Node.js
FROM node:22-alpine AS frontend-builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY assets/ ./assets/
COPY public/ ./public/
COPY views/ ./views/
COPY vite.config.js ./
RUN npm run build

# Stage 2: Build Go application
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Embedded files
COPY --from=frontend-builder /app/dist ./dist
COPY sql/pragmas.sql ./sql/
# Go source files
COPY db ./db/
COPY handlers ./handlers/
COPY middlewares ./middlewares/
COPY server ./server/
COPY views ./views/
COPY main.go ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app .

# Stage 3: Runtime image
FROM alpine:latest

COPY --from=backend-builder /app/bin/app /app

EXPOSE 3000

CMD ["/app"]
