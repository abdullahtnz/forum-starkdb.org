# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/ ./
RUN npm ci --legacy-peer-deps
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.22-alpine AS backend-builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# Stage 3: Final Image - Run Both
FROM alpine:3.20

RUN apk add --no-cache nginx ca-certificates && \
    mkdir -p /run/nginx && \
    mkdir -p /var/log/nginx

WORKDIR /app

# Copy backend
COPY --from=backend-builder /app/server .
RUN mkdir -p /app/uploads

# Copy frontend files
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/http.d/default.conf

EXPOSE 80 8080

# Start both services directly (no start.sh needed!)
CMD sh -c "nginx -g 'daemon off;' & ./server"
