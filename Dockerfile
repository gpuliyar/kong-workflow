## Stage 1 - Build Plugins
FROM golang:1.16-alpine as builder

RUN mkdir /plugins

WORKDIR /plugins

RUN apk add --no-cache git gcc libc-dev

RUN git clone https://github.com/Kong/go-pluginserver.git; \
    cd go-pluginserver; \
    go build; \
    cp go-pluginserver /usr/local/bin/

RUN cd /plugins

COPY go.mod go.sum client-auth.go ./

RUN go build -buildmode plugin -o client-auth.so client-auth.go

## Stage 2 - Bundle Kong with Plugings
FROM kong:latest

COPY --from=builder /usr/local/bin/go-pluginserver /usr/local/bin/go-pluginserver

RUN mkdir /tmp/go-plugins

COPY --from=builder /plugins/client-auth.so /tmp/go-plugins/client-auth.so

COPY kong.yaml /tmp/kong.yaml
