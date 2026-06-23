# Stage 1: Kompilasi Golang menjadi Biner Tunggal
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Pre-install govulncheck agar image builder bersifat immutable dan cepat
RUN go install golang.org/x/vuln/cmd/govulncheck@latest
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Stage 2: Lingkungan Produksi + Instalasi Docker CLI
FROM alpine:latest
WORKDIR /app

# [SANGAT PENTING] Backend kita mengeksekusi perintah 'docker run', 
# jadi OS di dalam kontainer ini wajib memiliki Docker CLI.
# Kita juga menginstal 'go' agar bisa melakukan kompilasi secara instan (0.1s) dibanding harus 'docker run golang' (15s).
RUN apk update && apk add --no-cache docker-cli go

COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]
