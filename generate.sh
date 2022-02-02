#!/bin/sh
export PATH=$GOPATH/src/github.com/hyperledger/fabric/build/bin:${PWD}/bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}
CHANNEL_NAME=consentmanagementchannel

if [ -d "config/" ]; then
  # Control will enter here if config/* exists.
	echo "direcory config already exists"
	exit
fi

echo "creating config directory"
mkdir config/

if [ -d "crypto-config/" ]; then
  # Control will enter here if config/* exists.
	echo "direcory crypto-config already exists"
	exit
fi

echo "creating crypto-config directory"
mkdir crypto-config/

# generate crypto material
cryptogen generate --config=./crypto-config.yaml
if [ "$?" -ne 0 ]; then
  echo "Failed to generate crypto material..."
  exit 1
fi

# generate genesis block for orderer
configtxgen -profile OrdererGenesis -outputBlock ./config/genesis.block -channelID "${CHANNEL_NAME}-config"
if [ "$?" -ne 0 ]; then
  echo "Failed to generate orderer genesis block..."
  exit 1
fi

# generate channel configuration transaction
configtxgen -profile consentmanagementchannel -outputCreateChannelTx ./config/channel.tx -channelID $CHANNEL_NAME
if [ "$?" -ne 0 ]; then
  echo "Failed to generate channel configuration transaction..."
  exit 1
fi

# generate anchor peer transaction
configtxgen -profile consentmanagementchannel -outputAnchorPeersUpdate ./config/UHopperOrgMSPanchors.tx -channelID $CHANNEL_NAME -asOrg UHopperOrgMSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for Org1MSP..."
  exit 1
fi

# generate anchor peer transaction
configtxgen -profile consentmanagementchannel -outputAnchorPeersUpdate ./config/UHopperOrg1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg UHopperOrg1MSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for Org2MSP..."
  exit 1
fi
