# Etapa de construcción
FROM golang:1.21-alpine AS builder

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar los archivos del módulo y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o httpserver ./cmd/httpserver

# Etapa de producción
FROM alpine:latest

# Instalar certificados CA
RUN apk --no-cache add ca-certificates

# Establecer el directorio de trabajo
WORKDIR /root/

# Copiar el binario desde la etapa de construcción
COPY --from=builder /app/httpserver .

# Puerto expuesto (debe coincidir con el puerto en tu código o variable de entorno)
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./httpserver"]
