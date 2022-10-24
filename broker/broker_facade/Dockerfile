#FROM golang:1.19.1-alpine
#
#ENV GO111MODULE=on
#WORKDIR /app
#
#COPY go.mod .
#COPY go.sum .
#RUN go mod download
#
##COPY *.go ./broker/broker_facade/
##COPY *.yaml ./broker/broker_facade/
#
#COPY . .
#
##RUN go build -o facade
#
##EXPOSE 5012
#
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o facade
#
#EXPOSE 5012
##ENTRYPOINT ["/app/broker/broker_facade"]
#
#CMD [ "/app/facade" ]


# syntax=docker/dockerfile:1

## Build
FROM golang:1.19.1-alpine AS build

ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /facade/program

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /facade /facade
COPY settings.yaml /

EXPOSE 5012

USER nonroot:nonroot

ENTRYPOINT ["/facade/program"]