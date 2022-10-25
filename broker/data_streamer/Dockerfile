## Build
FROM golang:1.19.1-alpine AS build

ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /data_streamer/program

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /data_streamer /data_streamer
COPY settings.yaml /

EXPOSE 5014

USER nonroot:nonroot

ENTRYPOINT ["/data_streamer/program"]