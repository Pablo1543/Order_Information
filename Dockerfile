FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy

RUN go build -o order_information ./cmd/order_service
RUN go build -o publisher ./cmd/publisher

FROM debian:bookworm-slim

WORKDIR /app

# Копируем бинарь и нужные файлы
COPY --from=builder /app/order_information .
COPY --from=builder /app/publisher .
COPY pages ./pages
COPY schema.sql .
COPY model.json .

EXPOSE 8080

CMD ["./order_information"]
