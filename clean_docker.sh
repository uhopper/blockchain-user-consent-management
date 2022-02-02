#!/bin/sh

echo "Cleaning docker images ..."
docker rmi $(docker images --format {{.Repository}}:{{.Tag}} | grep dev-peer0.org1.u-hopper.com-consentmanagementcontract)
docker rmi $(docker images --format {{.Repository}}:{{.Tag}} | grep dev-peer0.org.u-hopper.com-consentmanagementcontract)
