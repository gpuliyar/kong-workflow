#!/bin/bash

docker-compose stop
docker-compose rm -f
docker system prune -f
docker volume rm $(docker volume ls)
docker network rm $(docker network ls)
