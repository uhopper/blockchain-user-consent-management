#!/bin/sh

echo "Cleaning docker images ..."
docker rmi $(docker images --format {{.Repository}}:{{.Tag}} | grep dev-peer0.org.u-hopper.com-mortadella_contract)
docker rmi $(docker images --format {{.Repository}}:{{.Tag}} | grep dev-peer0.org2..com-mortadella_contract)
