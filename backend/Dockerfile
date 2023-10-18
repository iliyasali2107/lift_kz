# Latest golang image on Alpine Linux
FROM golang:1.21.1-alpine3.18

# Work directory
WORKDIR /app

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download



# Copying all the files
COPY . .

# Starting our application
CMD ["sh", "-c", "cd cmd/api/ && go run main.go"]


# Exposing server port
EXPOSE 8080


