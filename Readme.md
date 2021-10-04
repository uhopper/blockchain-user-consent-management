
## Intialize

Download the binaries and the docker image, type:

```
./init.sh
```

## Run the BlockChain


First generete the cryptografic material:

```
./generate.sh
```

Then start the blockchain:

```
./run.sh
```

Wait a couple of minuts in order to allow every node to complite the bootstrap procedure, then configure the network, by typing:

```
./configure.sh
```
 You need to configure the network only the first time you run the blockchain or after an execution of the clean procedure

Finaly deploy the chaincode

```
./deploy_chaincode.sh -v 0.0.1
```

If you edit the chaicode and you want to run the new code use the deploy_chaincode scrip with a new version and a new sequence number:

```
./deploy_chaincode.sh -v 0.0.2 -s 2
```

## Stop the blockchain

```
./stop.sh
```

## Clean all data

```
./clean_all.sh
```
