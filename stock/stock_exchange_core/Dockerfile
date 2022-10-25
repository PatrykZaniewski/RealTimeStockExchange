## Build
FROM golang:1.19.1-alpine AS build

ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /stock_exchange_core/program

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /stock_exchange_core /stock_exchange_core
COPY settings.yaml /

EXPOSE 5005

USER nonroot:nonroot

ENTRYPOINT ["/stock_exchange_core/program"]