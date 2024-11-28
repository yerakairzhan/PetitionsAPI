FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY configs/config.json /app/configs/config.json

RUN go build -o main .

EXPOSE 5556

CMD ["./main"]
