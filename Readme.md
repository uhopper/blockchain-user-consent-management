# Tsundoku blockchan

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [Intoduction](#intoduction)
- [Configure the environment](#configure-the-environment)
- [Configure the blockchain](#configure-the-blockchain)
  - [crypto-config file](#crypto-config-file)
  - [configtx file](#configtx-file)
- [Run the BlockChain](#run-the-blockchain)
- [Note](#note)
- [Stop the blockchain](#stop-the-blockchain)
- [Clean all data](#clean-all-data)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Intoduction

This repository contains the sample deployment of the blockchain used for maintaining the user consent of the personal data processing. This repository also contains a simple web service that allows to easily access the data.

This README will guide you in the configuration and the deployment of the blockain. 

The deployment will contains:
- One orderer
- Two peers
- Two CouchDB (used by the peers for mainening the word state)
- One WS for an easy access to the blockchain


## Configure the environment

- First create a file called  `.env` in the `docker-compose` folder and add the following variables:

```bash
ENV=<enviroment-name> # the  name of you enviroment, this will allow you to deploy multiple blockchains on the same server
AUTHORIZED_APIKEY=<ws-apikey> # choose an apikey for protect your ws
```

The `.env.template` file in the`docker-compose/template` can be used as example.

- Then create `couch_db_peer0.org0.env` file in the `docker-compose` directory with the following content:

```bash
CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=<username of the coutch db user of peer 1>
CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=<password of the coutch db user of peer 1>
COUCHDB_USER=<username of the coutch db user of peer 1>
COUCHDB_PASSWORD=<password of the coutch db user of peer 1>
```

The `couch_db_peer0.org0.env.template` file in the`docker-compose/template` can be used as example.

- Finally create the file `couch_db_peer0.org1.env` as the previous case.

## Configure the blockchain

### crypto-config file

By editing the `crypto-config.yaml` file is possible to change tge organization involved in the blockchain and the initial number of peers and orderer.

### configtx file

The `configtx.yaml` file containes the initial configuration of the blockchain. By editing this file is possible to change the permission and roles of the involved organizations and the intial channel and peers configuration.

## Run the BlockChain


First, you have to download the hyperledger binaries and the required docker images, simply type:

```
./init.sh
```

To run the blockchain you have to first generate all the cryptographic material required by the peers, orderers, and client applications. You can use the `generate.sh` script. Simply type:

```
./generate.sh
```

Then you can start the blockchain containers:

```
./run.sh
```

Wait a couple of minutes to allow every node to complete the bootstrap procedure, then configure the network, by typing:

```
./configure.sh
```
This script will create the `tsundochannel` and will add your peers to the channel. You need to configure the network only the first time you run the blockchain.

Finally deploy the chaincode

```
./deploy_chaincode.sh -v 0.0.1 -s 1
```

If you edit the chaincode and you want to run the new code, you can use the deploy_chaincode scrip again with a new version and a new sequence number:

```
./deploy_chaincode.sh -v 0.0.2 -s 2
```

Run the webservice

```
./run_ws
```

## Note

* If multiple instances of the blockchain are located on the same server, is suggested to make this explicit in the chaincode deployment command. For example:
```
./deploy_chaincode.sh -v 0.0.1-prod -s 1
```

```
./deploy_chaincode.sh -v 0.0.1-staging -s 1
```

## Stop the blockchain

It is possible to shut down the blockchain using the `stop.sh` script. Type:

```
./stop.sh
```

## Clean all data

```
./clean_all.sh
```
