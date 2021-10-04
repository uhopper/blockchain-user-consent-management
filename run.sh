#!/bin/sh
docker network create tsundoku
docker-compose -f docker-compose/docker-compose-solo.yaml up -d
