[ [<- index](/instructions/README.md) / [next ->](/instructions/LIGHTNING-01-install.md) ]

# Lightning
This guide can be used to perform an instant atomic swap between two exchanges A & B on the lightning network(BTC, LTC, XSN are supported). The exchanges are running the lightning network daemon (`lnd`) and can be connected to the bitcoin, litecoin and xsn networks using `btcd`,` ltcd`, `bitcoind`, `litecoind`, `xsnd` daemons. The guide assumes that the reader has some lightning knowledge and is capable of setting up a simple network, fund wallets and open payment channels.

# Describing the setup on XSN <-> LTC payment

## Exchange setup
Each exchange (A, B) is running several components which together provide full support for an atomic swap between XSN and LTC:
1. `xsnd` - full node, connected to the xsn chain
2. `litecoind` - full node, connected to the litecoin chain
3. `LND-XSN` - lightning network daemon for the XSN network
4. `LND-LTC` - lightning network daemon for the LTC network
5. `swap-resolver` - a simulator for [XUD](https://github.com/exchangeunion/xud), Exchange Union's decentralized exchange layer. This component is passing payment hashes and pre-images between the `LND-BTC`, `LND-LTC` and `LND-XSN` instances.

## Multiple Exchanges
We are going to setup two exchanges on a single machine. For that we would need to run 2x5 processes. In this guide we will share the `litecoind` and `xsnd` instances between the two exchanges, so we'll only need 8 processes.

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
You will need the swap-resolver to pass on payment hashes and pre-images between `LND-XSN` and `LND-LTC` to allow atomic for swaps between XSN and LTC on lightning.

Install `swap-resolver`
```shell
git clone https://github.com/X9Developers/swap-resolver.git $GOPATH/src/github.com/ExchangeUnion/swap-resolver
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

for XSN support also we need to download and apply xsn patch:
patch location: https://github.com/X9Developers/lnd/tree/master/patches/xsn_swap_integration.patch

git apply xsn_swap_integration.patch
```

## Setup Blockchain Clients

#### XSN

Start `xsnd` with zmq enabled
```shell
xsnd --testnet --rpcuser=user --rpcpassword=pass --zmqpubrawblock=tcp://127.0.0.1:28444 --zmqpubrawtx=tcp://127.0.0.1:28445 --daemon
```

#### Litecoin
 
Start `litecoind`

```shell
litecoind --testnet --rpcuser=user --rpcpassword=pass --zmqpubrawblock=tcp://127.0.0.1:29332 -zmqpubrawtx=ipc:///tmp/litecoind.tx.raw --daemon
```

### Coffee time
It will take some time for `xsnd` and `litecoind` to sync the chains. Time to get a coffee or two.


## Running Lightning Daemon(s)
Once testnet sync is done for the xsn & litecoin daemons, we continue to the next section which explains how to setup two `lnd` processes for Exchange-A and Exchange-B

[ [<- index](/instructions/README.md) / [next ->](/instructions/LIGHTNING-01-peers.md) ]
