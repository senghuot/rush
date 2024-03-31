# Use a specific platform for the base image
FROM golang:1.19 AS build

WORKDIR /app

# Copy the go.mod and go.sum files first to leverage Docker layer caching
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the rest of the source code
COPY main.go .
COPY payload.go .

# Cross-compile the application for Linux x86_64
# This will ensure the binary is compatible with the x86_64 architecture
RUN GOOS=linux GOARCH=amd64 go build -o ./bin/dist .

# Final stage: Use a distroless image for security and a smaller footprint
# You might want to use the nonroot tag for the image if your app doesn't require root permissions
FROM gcr.io/distroless/static-debian11

# Copy the compiled binary from the builder stage
COPY --from=build /app/bin/dist /app/

# Run the binary
CMD ["/app/dist"]

# RUN THIS TO GET IT WORK
# docker buildx build --platform linux/amd64 . -t go-container:latest