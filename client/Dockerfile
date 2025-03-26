FROM golang:latest as builder


WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN  go build -C .  -o  ./bin


FROM debian:bookworm
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=builder /app/bin .
EXPOSE 1236


CMD ["./bin"]
