## Build
FROM golang:1.19.1-alpine AS build

ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /broker_core/program

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /broker_core /broker_core
COPY settings.yaml /

EXPOSE 5011

USER nonroot:nonroot

ENTRYPOINT ["/broker_core/program"]