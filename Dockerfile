# Stage 1: Build Svelte frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY src ./src
COPY index.html ./
COPY vite.config.js ./
COPY tsconfig.json ./
COPY jsconfig.json ./
COPY components.json ./

RUN mkdir -p static

RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# Stage 3: Final runtime image
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=backend-builder /app/server .

COPY --from=frontend-builder /app/public ./public

RUN mkdir -p uploads

# Expose port
EXPOSE 8088

ENV UPLOAD_PASSWORD=demo \
    MAX_FILE_SIZE_MB=100 \
    FILE_EXPIRY_MINUTES=10 \
    PORT=:8088

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8088/ || exit 1

CMD ["./server"]