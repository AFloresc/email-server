# 🟦 Build stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copia los archivos de dependencias primero
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copia el resto del código
COPY . .

# Compila el binario para Linux sin CGO
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# 🟩 Run stage
FROM alpine:latest

WORKDIR /root/

# Instala certificados TLS para HTTPS
RUN apk --no-cache add ca-certificates

# Copia el binario compilado
COPY --from=builder /app/server .

# Expone el puerto (Render detecta automáticamente)
EXPOSE 10000

# Ejecuta el servidor
CMD ["./server"]