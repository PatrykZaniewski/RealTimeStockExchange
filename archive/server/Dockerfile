FROM golang:1.18.2-alpine

COPY . /go/src/stock-exchange-app
WORKDIR /go/src/stock-exchange-app

RUN go mod download

ENV PORT_RUN 50051

CMD go run main.go

EXPOSE $PORT_RUN