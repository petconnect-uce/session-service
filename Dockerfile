FROM golang:1.21.5-alpine3.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./

# Configuramos variables para mejorar la conexión al proxy de módulos
ENV GOPROXY=https://proxy.golang.org,direct
ENV GODEBUG=http2client=0

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3014

CMD ["./main"]
