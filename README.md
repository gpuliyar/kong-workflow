# kong-workflow

### Compile plugin code

```shell
docker run --rm -v $(pwd):/plugins kong/go-plugin-tool:latest-alpine-latest build client-auth.go
```

docker run -ti --rm --name kong-go-plugins \
  -e "KONG_DATABASE=off" \
  -e "KONG_GO_PLUGINS_DIR=/tmp/go-plugins" \
  -e "KONG_DECLARATIVE_CONFIG=/tmp/config.yml" \
  -e "KONG_PLUGINS=client-auth" \
  -e "KONG_PROXY_LISTEN=0.0.0.0:8000" \
  -p 8000:8000 \
  kong-demo
