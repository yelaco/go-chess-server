# Build stage
FROM golang:1.22-alpine

WORKDIR /app

# Cache module downloads
ENV GOCACHE=/go-cache
ENV GOMODCACHE=/gomod-cache
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/gomod-cache \
	go mod download

# Copy source code and build
COPY . .
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
	go build -v -o gchess-server ./cmd/server/

# Expose the necessary ports
EXPOSE 80
EXPOSE 443
EXPOSE 7201
EXPOSE 7202

# Run the server
CMD ["./gchess-server"]
