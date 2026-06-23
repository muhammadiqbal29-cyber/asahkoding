# Stage 1: Kompilasi Golang menjadi Biner Tunggal
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Instal Docker CLI di builder agar Unit Test (yang memanggil docker run) bisa berjalan
RUN apk update && apk add --no-cache docker-cli
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Stage 2: Lingkungan Produksi + Instalasi Docker CLI
FROM alpine:latest
WORKDIR /app

# [SANGAT PENTING] Backend kita mengeksekusi perintah 'docker run', 
# jadi OS di dalam kontainer ini wajib memiliki Docker CLI.
RUN apk update && apk add --no-cache docker-cli

COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]
