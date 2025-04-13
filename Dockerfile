FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o avito-app ./cmd/server/main.go

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/avito-app .
ENTRYPOINT ["./avito-app"]
