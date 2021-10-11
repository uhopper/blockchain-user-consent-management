#!/bin/sh


#loading .env file
if [ -f ./docker-compose/.env ]
then
  export $(cat ./docker-compose/.env | sed 's/#.*//g' | xargs)
fi

docker network create tsundoku-blockchain-${ENV}
docker-compose -f docker-compose/docker-compose-solo.yaml --project-name=tsundoku-blockchain-${ENV} up -d
