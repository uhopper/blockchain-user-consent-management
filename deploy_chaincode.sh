#!/bin/sh

VERSION="0.0.1"
TSL_ENABLE=false
CONTRACT_NAME=consentmanagementcontract
CHANNEL_NAME="consentmanagementchannel"
SEQUENCE="1"

#loading .env file
if [ -f ./docker-compose/.env ]
then
  export $(cat ./docker-compose/.env | sed 's/#.*//g' | xargs)
fi


# args=`getopt tv: $*`
# set -- $args
OPTIND=1
while getopts "htv:s:" opt;
do
	case "$opt" in
		h|\?)
        echo "use -v <version> to set the version of the code"
				echo "use -t to enable tsl (for etcdraft)"
        exit 0
        ;;
    v)  VERSION=$OPTARG
        ;;
    t)  TSL_ENABLE=true
        ;;
    s)  SEQUENCE=$OPTARG
        ;;
esac

done

echo "version: ${VERSION}"
echo "tsl enable: ${TSL_ENABLE}"


if $TSL_ENABLE; then
	echo "tsl enabled"
    echo "NOT IMPLEMENTED"
    exit 1
else
	
    echo "Packaging the chaincode"
    docker exec -w /opt/gopath/src/github.com/hyperledger/fabric/peer/chaincode_pkg cliUHopperOrg-${ENV}  peer lifecycle chaincode package ${CONTRACT_NAME}.tar.gz --path /opt/gopath/src/consent_management_chaincode --lang golang --label ${CONTRACT_NAME}_${VERSION}

    echo "Installing on peer0 org1"
    docker exec -w /opt/gopath/src/github.com/hyperledger/fabric/peer/chaincode_pkg cliUHopperOrg-${ENV}  peer lifecycle chaincode install ${CONTRACT_NAME}.tar.gz

    echo "Installing on peer0 org2"
    docker exec -w /opt/gopath/src/github.com/hyperledger/fabric/peer/chaincode_pkg cliUHopperOrg1-${ENV}  peer lifecycle chaincode install ${CONTRACT_NAME}.tar.gz

    echo "Check if the chaincode is installed"
    docker exec cliUHopperOrg-${ENV}  peer lifecycle chaincode queryinstalled
    docker exec cliUHopperOrg1-${ENV}  peer lifecycle chaincode queryinstalled

    docker exec cliUHopperOrg1-${ENV}  peer lifecycle chaincode queryinstalled >log.txt


    PACKAGE_ID=$(cat log.txt| tail -n 1 |sed -n "/${CC_NAME}_${CC_VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}")
    echo "PACKAGE_ID: ${PACKAGE_ID}"

    echo "Approve for org1"
    docker exec cliUHopperOrg-${ENV}  peer lifecycle chaincode approveformyorg -o orderer.orderer.u-hopper.com:7050 --ordererTLSHostnameOverride orderer.orderer.u-hopper.com  --channelID $CHANNEL_NAME --name ${CONTRACT_NAME} --version ${VERSION} --package-id ${PACKAGE_ID} --sequence ${SEQUENCE} #--init-required 

    echo "Approve for org2"
    docker exec cliUHopperOrg1-${ENV}  peer lifecycle chaincode approveformyorg -o orderer.orderer.u-hopper.com:7050 --ordererTLSHostnameOverride orderer.orderer.u-hopper.com  --channelID $CHANNEL_NAME --name ${CONTRACT_NAME} --version ${VERSION} --package-id ${PACKAGE_ID} --sequence ${SEQUENCE} #--init-required 

    echo "Check if the chaincode is ready to be committed"
    docker exec cliUHopperOrg-${ENV}  peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name ${CONTRACT_NAME} --version ${VERSION} --sequence ${SEQUENCE} --output json #--init-required

    echo "Commit the chaincode"
    docker exec cliUHopperOrg-${ENV}  peer lifecycle chaincode commit -o orderer.orderer.u-hopper.com:7050  --ordererTLSHostnameOverride orderer.orderer.u-hopper.com --channelID $CHANNEL_NAME --name ${CONTRACT_NAME} --peerAddresses peer0.org.u-hopper.com:7051 --peerAddresses peer0.org1.u-hopper.com:7051  --version ${VERSION} --sequence ${SEQUENCE} #--init-required

fi
