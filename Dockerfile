FROM golang:1.22.4 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN cp .env.example .env

RUN go build -o ./out/go-server .

EXPOSE 8080

CMD ["./out/go-server"]