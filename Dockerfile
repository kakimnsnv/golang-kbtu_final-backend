FROM golang:1.23.3-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o e-commerce-gin ./cmd/server/main.go

EXPOSE 8080

CMD ["./e-commerce-gin"]
