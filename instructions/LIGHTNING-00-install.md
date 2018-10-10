[ [<- index](/README.md) / [next ->](/LIGHTNING-01-peers.md) ]

# Lightning
This guide can be used to perform an instant BTC<->LTC atomic swap between two exchanges A & B on the lightning network. The exchanges are running the lightning network daemon (`lnd`) and are connected to the bitcoin and litecoin networks using `btcd` and `ltcd`. The guide assumes that the reader has some lightning knowledge and is capable of setting up a simple network, fund wallets and open payment channels.

# Describing the setup

## Exchange setup
Each exchange (A, B) is running several components which together provide full support for an atomic swap between BTC and LTC:
1. `BTCD` - full node, connected to the bitcoin chain
2. `LTCD` - full node, connected to the litecoin chain
3. `LND-BTC` - lightning network daemon for the BTC network
4. `LND-LTC` - lightning network daemon for the LTC network
5. `swap-resolver` - a simulator for [XUD](https://github.com/exchangeunion/xud), Exchange Union's decentralized exchange layer. This component is passing payment hashes and pre-images between the `LND-BTC` and `LND-LTC` instances.

## Multiple Exchanges
We are going to setup two exchanges on a single machine. For that we would need to run 2x5 processes. In this guide we will share the `BTCD` and `LTCD` instances between the two exchanges, so we'll only need 8 processes.

Since everything is running on a single machine, we need to assign a different working directory and ports for each process. More about this later.

### Terminal setting
It is recommended to open 9 terminals to be used with this PoC, allowing each component to run in its own terminal with an additional terminal being used for commands (CLI).

# Installation 
## Install Lightning Dependencies
Download latest `Go` package from [official Go repository](https://golang.org/dl/) and decompress to `/usr/local` .
Alternatively you can use apt-get
```shell
sudo apt-get install golang-1.10-go
```

Add the following lines to the end of `$HOME/.bashrc` and source the file (or optionally reboot)
```shell
# add Go paths
export GOPATH=$HOME/go
export PATH=/usr/local/go/bin:$GOPATH/bin:$PATH
```

Install `Glide`
```shell
go get -u github.com/Masterminds/glide
```

## Installing `swap-resolver`
You will need the swap-resolver to pass on payment hashes and pre-images between `LND-BTC` and `LND-LTC` to allow atomic for swaps between BTC and LTC on lightning.

Install `swap-resolver`
```shell
git clone https://github.com/ExchangeUnion/swap-resolver.git $GOPATH/src/github.com/ExchangeUnion/swap-resolver
cd $GOPATH/src/github.com/ExchangeUnion/swap-resolver
dep ensure
```

## Build `lnd`

Since cross-chain swaps support was not yet merged into the official `lnd` master branch, we will use a slightly modified swap enabled `lnd` instead. 

#### Swap enabled `lnd`

Install the swap enabled `lnd`

```shell
git clone -b resolver https://github.com/ExchangeUnion/lnd.git $GOPATH/src/github.com/lightningnetwork/lnd
cd $GOPATH/src/github.com/lightningnetwork/lnd
make && make install
```

#### Bitcoin full node implementation `btcd`
If you don't have `btcd` set up yet, you will need to build and install the `btcd` bitcoin full node implementation
```shell
git clone https://github.com/roasbeef/btcd $GOPATH/src/github.com/roasbeef/btcd
cd $GOPATH/src/github.com/roasbeef/btcd
glide install
go install . ./cmd/...
```

#### Litecoin full node implementation `ltcd`
If you don't have `ltcd` set up yet, you will need to build and install the `ltcd` litecoin full node implementation
```shell
git clone https://github.com/ltcsuite/ltcd.git $GOPATH/src/github.com/ltcsuite/ltcd
cd $GOPATH/src/github.com/ltcsuite/ltcd
glide install
go install . ./cmd/...
```

## Setup Blockchain Clients

#### Bitcoin

Start `btcd`
```shell
btcd --testnet --txindex --rpcuser=xu --rpcpass=xu
```
The bitcoin testnet blockchain is downloaded to `$HOME/.btcd` by default.

Sync progress may be tracked as follows:
```shell
btcctl --testnet getinfo --rpcuser=xu --rpcpass=xu
```

#### Litecoin
 
Start `ltcd`

```shell
ltcd --testnet --txindex --rpcuser=xu --rpcpass=xu
```
litecoin testnet blockchain is downloaded to `$HOME/.ltcd` by default.

Sync progress may be tracked as follows
```shell
ltcctl --testnet  getinfo --rpcuser=xu --rpcpass=xu
```

### Coffee time
It will take some time for `btcd` and `ltcd` to sync with the chains. Time to get a coffee or two.


## Running Lightning Daemon(s)
Once testnet sync is done for the bitcoin & litecoin daemons, we continue to the next section which explains how to setup two `lnd` processes for Exchange-A and Exchange-B

[ [<- index](/README.md) / [next ->](/LIGHTNING-01-peers.md) ]
