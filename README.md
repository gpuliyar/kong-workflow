# kong-workflow

```
docker network create kong-net
```

### Start Cassandra

```
docker run -d --name kong-database \
               --network=kong-net \
               -p 9042:9042 \
               cassandra:3
```

### After running `deck sync`

#### Create Consumer

```
http :8001/consumers username=consumer custom_id=consumer
```

#### Create API key

```
http :8001/consumers/consumer/key-auth key=b3637a1c3f4e
```

#### Check API Key

```
http :8000/mock/request apikey:b3637a1c3f4e
```
