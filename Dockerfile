FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o main ./cmd/server/main.go

EXPOSE $PORT

CMD ["./main"]