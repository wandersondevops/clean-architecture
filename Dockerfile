# Dockerfile

# Etapa 1: Construção da aplicação Go
FROM golang:1.22.5 as builder

WORKDIR /app

# Copia o go.mod e go.sum e instala as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte da aplicação
COPY . .

# Compila a aplicação Go
RUN go build -o ordersystem ./cmd/ordersystem

# Instala o migrate para migrações de banco de dados
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Etapa 2: Imagem Final
FROM golang:1.22.5

WORKDIR /app

# Copia o binário compilado da etapa 1
COPY --from=builder /app/ordersystem .

# Copia o migrate da etapa 1
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Copia o arquivo .env e o diretório de migrações
COPY .env /app/.env
COPY ./internal/infra/database/migrations /app/migrations
