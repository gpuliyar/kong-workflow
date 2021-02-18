## Stage 1 - Build Plugins
FROM golang:alpine as builder

RUN apk add --no-cache git gcc libc-dev

RUN go get github.com/Kong/go-pluginserver

RUN mkdir /go-plugins

COPY go-plugins/client-auth.go /go-plugins/.
COPY go.mod /go-plugins/.
COPY go.sum /go-plugins/.

WORKDIR /go-plugins

RUN go build -buildmode plugin -o client-auth.so client-auth.go

## Stage 2 - Bundle Kong with Plugings
FROM kong:2.3-alpine

COPY --from=builder /go/bin/go-pluginserver /usr/local/bin/go-pluginserver

RUN mkdir /tmp/go-plugins

COPY --from=builder /go-plugins/client-auth.so /tmp/go-plugins/client-auth.so

COPY kong.yaml /tmp/config.yml

USER root

RUN chmod -R 775 /tmp
