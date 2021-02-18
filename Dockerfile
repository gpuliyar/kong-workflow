## Stage 1 - Build Plugins
FROM kong/go-plugin-tool:latest-alpine-latest as builder

COPY go.mod go.sum client-auth.go ./

RUN go build -buildmode plugin -o client-auth.so client-auth.go

## Stage 2 - Bundle Kong with Plugings
FROM kong:latest

COPY --from=builder /usr/local/bin/go-pluginserver /usr/local/bin/go-pluginserver

RUN mkdir /tmp/go-plugins

COPY --from=builder /plugins/client-auth.so /tmp/go-plugins/client-auth.so

COPY kong.yaml /tmp/kong.yaml
