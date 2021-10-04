#!/bin/sh

CHANNEL_NAME=tsundokchannel
# CHANNEL_NAME=mychannel

echo "creating channel"
docker exec -e "CORE_PEER_LOCALMSPID=UHopperOrgMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org.u-hopper.com/msp" peer0.org.u-hopper.com peer channel create -o orderer.orderer.u-hopper.com:7050 -c "${CHANNEL_NAME}" -f /etc/hyperledger/configtx/channel.tx

if [ "$?" -ne 0 ]; then
  echo "Failed to create the channel"
  exit 1
fi

sleep 2

echo "joinin peer1 to  channel"
docker exec -e "CORE_PEER_LOCALMSPID=UHopperOrgMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org.u-hopper.com/msp" peer0.org.u-hopper.com peer channel join  -b "${CHANNEL_NAME}.block"


if [ "$?" -ne 0 ]; then
  echo "Failed to join the channle"
  exit 1
fi

sleep 2

echo "fetch channel block for peer2"
docker exec -e "CORE_PEER_LOCALMSPID=UHopperOrg1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.u-hopper.com/msp" peer0.org1.u-hopper.com peer channel fetch 0 "${CHANNEL_NAME}.block" -c "${CHANNEL_NAME}" -o orderer.orderer.u-hopper.com:7050


if [ "$?" -ne 0 ]; then
  echo "Failed to fetch channel block for peer2"
  exit 1
fi

sleep 2

echo "join peer2 to the channel"

docker exec -e "CORE_PEER_LOCALMSPID=UHopperOrg1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.u-hopper.com/msp" peer0.org1.u-hopper.com peer channel join -b "tsundokchannel.block"


if [ "$?" -ne 0 ]; then
  echo "Failed to join peer2 to the channel"
  exit 1
fi

sleep 2

echo "update anchor peer configuration for Org1"

docker exec -e "CORE_PEER_LOCALMSPID=UHopperOrgMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org.u-hopper.com/msp" peer0.org.u-hopper.com peer channel update -o orderer.orderer.u-hopper.com:7050 -c $CHANNEL_NAME -f /etc/hyperledger/configtx/UHopperOrgMSPanchors.tx

if [ "$?" -ne 0 ]; then
  echo "Failed to update anchor peer configuration for Org1"
  exit 1
fi

sleep 2

echo "update anchor peer configuration for Org2"

docker exec -e "CORE_PEER_LOCALMSPID=UHopperOrg1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.u-hopper.com/msp" peer0.org1.u-hopper.com  peer channel update -o orderer.orderer.u-hopper.com:7050 -c $CHANNEL_NAME -f /etc/hyperledger/configtx/UHopperOrg1MSPanchors.tx

if [ "$?" -ne 0 ]; then
  echo "Failed to update anchor peer configuration for Org2"
  exit 1
fi
