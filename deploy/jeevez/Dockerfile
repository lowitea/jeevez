FROM golang:1.14 as builder

ADD ./go.mod \
    ./go.sum \
    /build/
WORKDIR /build
RUN go mod download

ADD ./cmd/ /build/cmd
ADD ./internal/ /build/internal

ENV GOOS=linux
RUN go build -o jeevez ./cmd/jeevez/main.go
RUN go build -o updater ./cmd/updater/main.go

FROM debian:buster

ARG JEEVEZ_APP_VERSION
ENV JEEVEZ_APP_VERSION=${JEEVEZ_APP_VERSION}

# for debug
RUN apt-get update && apt-get install curl -yqq

COPY --from=builder /build/jeevez /app/
COPY --from=builder /build/updater /app/

WORKDIR /app

CMD ["./jeevez"]
