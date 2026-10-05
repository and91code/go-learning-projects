# Imagem oficial do Go
FROM golang:1.25-alpine

# Instala ferramentas básicas de build
RUN apk add --no-cache git gcc musl-dev
RUN go install github.com/air-verse/air@v1.61.7

WORKDIR /app

# Copia os arquivos de dependência do Go
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o código-fonte
COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]