FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/moving-schedule-backend

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/moving-schedule-backend .

EXPOSE 8080

CMD ["./moving-schedule-backend"]
