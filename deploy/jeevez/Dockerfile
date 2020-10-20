FROM golang:1.14 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main ./cmd/jeevez/main.go

FROM debian:buster
ARG JEEVEZ_APP_VERSION
ENV JEEVEZ_APP_VERSION=${JEEVEZ_APP_VERSION}
# for debug
RUN apt-get update && apt-get install curl -yqq
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
