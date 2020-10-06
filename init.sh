#!/bin/bash


KEY="mykey"
KEY1="wade"
CHAINID=200812
MONIKER="mymoniker"

# remove existing chain environment, data and
rm -rf  ~/.desmos*

make install

desmoscli config keyring-backend test

# if mykey exists it should be deleted
desmoscli keys add $KEY
desmoscli keys add $KEY1

# Set moniker and chain-id for Ethermint (Moniker can be anything, chain-id must be an integer)
desmosd init $MONIKER --chain-id $CHAINID


os=`uname -a`
mac='Darwin'
if [[ $os =~ $mac ]];then

  sed -i ''  's/"max_gas": "-1"/"max_gas": "1000000000"/g'   ~/.desmosd/config/genesis.json
else
  sed -i 's/"max_gas": "-1"/"max_gas": "1000000000"/g'   ~/.desmosd/config/genesis.json
fi

# Set up config for CLI
desmoscli config chain-id $CHAINID
desmoscli config output json
desmoscli config indent true
desmoscli config trust-node true

# if $KEY exists it should be deleted


# Set moniker and chain-id for Ethermint (Moniker can be anything, chain-id must be an integer)

# Allocate genesis accounts (cosmos formatted addresses)
echo "add-genesis-account "
echo $(desmoscli keys show $KEY -a)
desmosd add-genesis-account $(desmoscli keys show $KEY -a)  100000000000000000000000000desmos
desmosd add-genesis-account $(desmoscli keys show $KEY1 -a) 100000000000000000000000000desmos   # --vesting-amount 600hale  --vesting-start-time  1591781100   --vesting-end-time  1591781700

# Sign genesis transaction
desmosd gentx --amount 1000000000desmos  --name $KEY --keyring-backend test

# desmosd gentx --amount 1000000000udaric --name jack

# Collect genesis tx
desmosd collect-gentxs
desmosd collect-gentxs




# Run this to ensure everything worked and that the genesis file is setup correctly
desmosd validate-genesis

# Command to run the rest server in a different terminal/window
echo -e '\n\nRun this rest-server command in a different terminal/window:'
echo -e "desmoscli rest-server --laddr \"tcp://localhost:8545\" --unlock-key $KEY --chain-id $CHAINID\n\n"

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
desmosd start   --pruning=nothing --rpc.unsafe --log_level "main:info,state:info,mempool:info"   --rpc.laddr "tcp://0.0.0.0:26657"
