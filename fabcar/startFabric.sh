#!/bin/bash

# Exit on first error
set -e

starttime=$(date +%s)

# a global environment variable overwrites the .env setting if available
#seems the domain is not changeable atm
export domain=example.com
export chaincode1=etrade
# chaincode gets cached, increase the version if you cant see any changes
export cc_ver=1.0.8


# Grab the current directoryinitLedger
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ ! -d ~/.hfc-key-store/ ]; then
	mkdir ~/.hfc-key-store/
fi
cp $PWD/network/creds/* ~/.hfc-key-store/

#
cd "${DIR}"/network

docker-compose -f "${DIR}"/network/docker-compose.yml up -d

# wait for Hyperledger Fabric to start
sleep 10

# Create the channel
# to change the channel name you have to provide another config file. (with configtxgen tool)
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer0.org1.$domain peer channel create -o orderer.$domain:7050 -c mychannel -f /etc/hyperledger/configtx/mychannel.tx
# Join peer0.org1.example.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.example.com/msp" peer0.org1.$domain peer channel join -b mychannel.block
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode install -n $chaincode1 -v $cc_ver -p github.com/$chaincode1
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode instantiate -o orderer.$domain:7050 -C mychannel -n $chaincode1 -v $cc_ver -c '{"Args":[""]}' -P "OR ('Org1MSP.member','Org2MSP.member')"
sleep 10
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.$domain:7050 -C mychannel -n $chaincode1 -c '{"function":"initLedger","Args":[""]}'

cd ../..

printf "\nTotal execution time : $(($(date +%s) - starttime)) secs ...\n\n"
