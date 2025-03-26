FROM golang:latest as builder


WORKDIR /app
ENV GOOS=linux
ENV CGO_ENABLED=1
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN  go build -C ./cmd  -o ../bin


FROM debian:bookworm
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=builder /app/bin .
COPY --from=builder /app/.env .

# Открываем порты
EXPOSE 50051
EXPOSE 9090

# Запускаем сервис
CMD ["./bin"]