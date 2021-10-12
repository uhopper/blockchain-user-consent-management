# Tsundoku blockchan

## Intialize

Download the binaries and the docker image, type:

```
./init.sh
```

## Configure the environment

In the `docker-compose/template` folder there are three templates. Move these files in the `docker-compose` folder and fill all the values.


## Run the BlockChain


First generete the cryptografic material:

```
./generate.sh
```

Then start the blockchain:

```
./run.sh
```

Wait a couple of minuts in order to allow every node to complete the bootstrap procedure, then configure the network, by typing:

```
./configure.sh
```
 You need to configure the network only the first time you run the blockchain or after an execution of the clean procedure

Finaly deploy the chaincode

```
./deploy_chaincode.sh -v 0.0.1 -s 1
```

If you edit the chaicode and you want to run the new code use the deploy_chaincode scrip with a new version and a new sequence number:

```
./deploy_chaincode.sh -v 0.0.2 -s 2
```

Run the webservice

```
./run_ws
```

## Note

* If multiple instance of the blockchain are located on the same sarver is suggested to make this explicit in the chaincode deployment command. For example:

```
./deploy_chaincode.sh -v 0.0.1-prod -s 1
```

```
./deploy_chaincode.sh -v 0.0.1-staging -s 1
```

## Stop the blockchain

```
./stop.sh
```

## Clean all data

```
./clean_all.sh
```
