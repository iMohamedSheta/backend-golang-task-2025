# Get the image from the official Golang image
FROM golang:1.24

# Install netcat (nc)
RUN apt-get update && apt-get install -y netcat-openbsd wget

# Set the working directory
WORKDIR /app

# Copy the current directory into the container
COPY . .

# Copy the .env file into the container
COPY .env.docker .env 

# Use the local vendor directory don't download it
ENV GOFLAGS=-mod=vendor

# Disable CGO to avoid gcc errors
ENV CGO_ENABLED=0


# Build the image
RUN go build -o app cmd/server/main.go

# Build migrate
RUN go build -o migrate cmd/migrate/main.go

# Build the worker
RUN go build  -o worker cmd/worker/main.go

# Run the container
CMD  ./app

