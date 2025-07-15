# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Instala dependências necessárias
RUN apk add --no-cache gcc musl-dev

# Copia os arquivos de dependências
COPY go.mod ./
COPY go.sum ./

# Baixa as dependências
RUN go mod download

# Copia o código fonte
COPY . .

# Compila a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

# Instala certificados CA
RUN apk --no-cache add ca-certificates

# Copia o binário compilado
COPY --from=builder /app/main .

# Expõe a porta da aplicação
EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Comando para executar a aplicação
CMD ["./main"] 