# kong-workflow

### Compile plugin code

```shell
docker run --rm -v $(pwd):/plugins kong/go-plugin-tool:latest-alpine-latest build client-auth.go
```
