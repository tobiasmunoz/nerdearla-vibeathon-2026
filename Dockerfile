# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o live-subs-app .

# Final minimal production container
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/live-subs-app /app/live-subs-app
COPY index.html /app/index.html
COPY samples/ /app/samples/

ENV PORT=8080
EXPOSE 8080

CMD ["/app/live-subs-app"]
