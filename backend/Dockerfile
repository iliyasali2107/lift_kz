# Latest golang image on Alpine Linux
FROM golang:1.21.1-alpine3.18

# Work directory
WORKDIR /app

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN echo 'https://dl-cdn.alpinelinux.org/alpine/v3.14/community' >> /etc/apk/repositories
RUN echo 'https://dl-cdn.alpinelinux.org/alpine/v3.14/main' >> /etc/apk/repositories
RUN apk add --no-cache wkhtmltopdf
# RUN apt-get update && \
#     apt-get install -y xvfb wkhtmltopdf ghostscript

# Copying all the files
COPY . .

# Starting our application
CMD ["sh", "-c", "cd cmd/api/ && go run main.go"]

# Exposing server port
EXPOSE 8080
