# Latest golang image on Alpine Linux
FROM golang:1.21.1-alpine3.18

# Work directory
WORKDIR /app

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN apk update && apk add --no-cache chromium

# Copying all the files
COPY . .

ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/

# Starting our application
CMD ["sh", "-c", "go run cmd/api/main.go"]


# Exposing server port
EXPOSE 8080


