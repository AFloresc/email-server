# 🟦 Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Copia los archivos de dependencias primero
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copia el resto del código
COPY . .

# Cambiamos al directorio donde está main.go
WORKDIR /app/cmd/server

# Compila el binario para Linux sin CGO
RUN CGO_ENABLED=0 GOOS=linux go build -o /server

# 🟩 Run stage
FROM alpine:latest

WORKDIR /root/

# Instala certificados TLS para HTTPS
RUN apk --no-cache add ca-certificates

# Copia el binario compilado
COPY --from=builder /server .

# Expone el puerto (Render detecta automáticamente)
EXPOSE 10000

# Ejecuta el servidor
CMD ["./server"]