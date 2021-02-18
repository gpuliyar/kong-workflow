## Stage 1 - Build Plugins
FROM golang:1.15-alpine3.13 as builder

RUN apk add --no-cache git gcc libc-dev make

RUN go get github.com/Kong/go-pluginserver

RUN mkdir /go-plugins

WORKDIR /go-plugins/

COPY go.mod go.sum client-auth.go ./

RUN go build -buildmode plugin -o /go-plugins/client-auth.so /go-plugins/client-auth.go

## Stage 2 - Bundle Kong with Plugings
FROM kong:2.3-alpine

COPY --from=builder /go/bin/go-pluginserver /usr/local/bin/go-pluginserver

RUN mkdir /tmp/go-plugins

COPY --from=builder /go-plugins/client-auth.so /tmp/go-plugins/client-auth.so

COPY kong.yaml /tmp/kong.yaml
