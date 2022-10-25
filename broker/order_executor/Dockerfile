## Build
FROM golang:1.19.1-alpine AS build

ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /order_executor/program

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /order_executor /order_executor
COPY settings.yaml /

EXPOSE 5010

USER nonroot:nonroot

ENTRYPOINT ["/order_executor/program"]