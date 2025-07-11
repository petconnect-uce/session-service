# Etapa 1: Build
FROM golang:1.21.5-alpine3.19 AS builder

WORKDIR /app

# Configura proxy para modulos y desactiva HTTP2 para evitar problemas en algunos proxys
ENV GOPROXY=https://proxy.golang.org,direct
ENV GODEBUG=http2client=0

# Copia solo archivos para dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copia el código fuente
COPY . .

# Build estático (opcionalmente podrías agregar CGO_ENABLED=0 si no usas Cgo)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o session-service .

# Etapa 2: Imagen runtime ligera
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/session-service .

# Permisos ejecución
RUN chmod +x ./session-service

# Expone el puerto que usas
ENV PORT=3014
EXPOSE 3014

# Comando para iniciar
CMD ["./session-service"]
