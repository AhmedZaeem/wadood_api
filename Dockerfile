# Start from the official Golang image
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o wadood

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/wadood .
EXPOSE 8080
CMD ["./wadood"]