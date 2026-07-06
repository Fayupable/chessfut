FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o chessfut-be .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/chessfut-be .
COPY --from=builder /app/db/migration ./db/migration

EXPOSE 8099

ENTRYPOINT ["./chessfut-be"]