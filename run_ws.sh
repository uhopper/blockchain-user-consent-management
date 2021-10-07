#!/bin/sh


#loading .env file
if [ -f ./docker-compose/.env ]
then
  export $(cat ./docker-compose/.env | sed 's/#.*//g' | xargs)
fi

docker network create tsundoku-${ENV}
docker-compose -f docker-compose/docker-compose-ws.yaml --project-name=tsundoku-blockchain-${ENV}-ws up -d
