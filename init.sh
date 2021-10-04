#!/bin/sh

if [ -d "bin/" ]; then
	echo "direcory bin already exists"
	exit
fi

echo "creating bin directory"
mkdir bin/

if [ -d "data/" ]; then
	echo "direcory data already exists"
	exit
fi

echo "creating data directory"
mkdir data/

echo "downloading binaries"

curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/master/scripts/bootstrap.sh | bash -s -- 2.3.2 1.0.0 1.5.0 -s

echo "removing sample config folder"
rm -r config
