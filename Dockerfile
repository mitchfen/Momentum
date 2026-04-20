# Build Stage
FROM golang:1.24-alpine AS build
WORKDIR /source

COPY src/go.mod ./
RUN go mod download

COPY src/ .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/momentum .

# Runtime Stage
FROM scratch
WORKDIR /app

COPY --from=build /app/momentum .

EXPOSE 80

ENTRYPOINT ["./momentum"]

LABEL org.opencontainers.image.description="Momentum - a daily habit stacking, task tracking application"
