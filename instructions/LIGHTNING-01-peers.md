[ [index](/README.md) | [<- previous](/LIGHTNING-00-install.md) / [next ->](/LIGHTNING-02-connect.md) ]

# Lightning nodes Setup

Once we installed all components and the XSN & Litecoin testnet blockchains are synced, we can set up the `lnd` and `swap-resolver` processes. We are going to setup two `swap-resolver` processes, simulating two exchanges A & B running `xud` and four `lnd` processes, one for XSN and one for LTC for each exchange. It takes some time until `lnd` synced with the `xsnd` and `litecoin` daemons. It is OK to setup the `lnd`s and `swap-resolver`s in parallel.

## Aliases
To make our life easier with `lncli`, we recommend setting up aliases. Add the following to to `~/.bash_profile` or `~/.profile` and source it:

```bash
#Adding lncli aliases
alias xa-lnd-xsn='lncli --network testnet --rpcserver=localhost:10003 --no-macaroons'
alias xa-lnd-ltc='lncli --network testnet --rpcserver=localhost:10001 --no-macaroons'
alias xb-lnd-xsn='lncli --network testnet --rpcserver=localhost:20003 --no-macaroons'
alias xb-lnd-ltc='lncli --network testnet --rpcserver=localhost:20001 --no-macaroons'
```

Now we can use these aliases to communicate with the 4 `lnd` processes without the need to type long CLI arguments.

## Startup Scripts
To make life even easier, we find the following directory structure in `$GOPATH/src/github.com/ExchangeUnion/swap-resolver`:

*	exchange-a
	+	lnd (resolve.conf)
		*	xsn (start.bash)
		*	ltc (start.bash)
	+	resolver (start.bash)
*	exchange-b
	+	lnd (resolve.conf)
		*	xsn (start.bash)
		*	ltc (start.bash)
	+	resolver (start.bash)

The `start.bash` script invokes the LND process using the right parameters (ports, etc). The `resolve.conf` is needed for the swap-resolver to function. Just FYI, no need to do anything for now.

## Exchange A
### Launch `lnd-xsn`
Open a terminal to set Exchange A's `lnd-xsn` daemon
```shell
cd $GOPATH/src/github.com/ExchangeUnion/swap-resolver/exchange-a/lnd/xsn/
./start.bash
```

check progress with
```shell
xa-lnd-xsn getinfo
```
### Launch `lnd-ltc`
Open a terminal to set Exchange A's `lnd-ltc` daemon
```shell
cd $GOPATH/src/github.com/ExchangeUnion/swap-resolver/exchange-a/lnd/ltc/
./start.bash
```

check progress with
```shell
xa-lnd-ltc getinfo
```

### Launch `swap-resolver`
Open a terminal to set Exchange A's `xud` daemon
```shell
cd $GOPATH/src/github.com/ExchangeUnion/swap-resolver/exchange-a/resolver/
./start.bash
```


## Exchange B
### Launch `lnd-xsn`
Open a terminal to set Exchange B's `lnd-xsn` daemon
```shell
cd $GOPATH/src/github.com/ExchangeUnion/swap-resolver/exchange-b/lnd/xsn/
./start.bash
```

check progress with
```shell
xb-lnd-xsn getinfo
```
### Launch `lnd-ltc`
Open a terminal to set Exchange B's `lnd-ltc` daemon
```shell
cd $GOPATH/src/github.com/ExchangeUnion/swap-resolver/exchange-b/lnd/ltc/
./start.bash
```

check progress with
```shell
xb-lnd-ltc getinfo
```
### Launch `swap-resolver`
Open a terminal to set Exchange B's `resolver` daemon
```shell
cd $GOPATH/src/github.com/ExchangeUnion/swap-resolver/exchange-b/resolver/
./start.bash
```


## Coffee time v2

Give the four `lnd`s some time to sync with `xsnd` and `litecoind`. You can check the status by using the `getinfo` command (use the cli terminal for this). You would want to see `"synced_to_chain": true,` for all four `lnd`s.

### Check status 

Example - Check the status of Exchange A's `lnd-xsn`:

```shell
xa-lnd-xsn getinfo
{
    "identity_pubkey": "02b07e2983eb179a7c3172c927eee88303f68c4e6ca0f21a971e73ec465da77149",
    "alias": "Exchange A XSN on 10003/10013",
    "num_pending_channels": 0,
    "num_active_channels": 0,
    "num_peers": 0,
    "block_height": 21402,
    "block_hash": "3eb5c50a752798edadc1df1593cf94db54e18af3fccd66fbaf4a7c8a768c305f",
    "synced_to_chain": true,
    "testnet": true,
    "chains": [
        "xsncoin"
    ],
    "uris": [
    ],
    "best_header_timestamp": "1539258342",
    "version": "0.5.0-beta commit=fec8a2f221c42fde67f63a5ba43b2a80b118e1cc"
```


# Fund Exchange A

## Balance after creating

Query Exchange A wallet balances for both `XSN` and `Litecoin` after creation (we expect to see zeros!)
```shell
$ xa-lnd-xsn walletbalance
{
    "total_balance": "0",
    "confirmed_balance": "0",
    "unconfirmed_balance": "0"
}
$ xa-lnd-ltc walletbalance
{
    "total_balance": "0",
    "confirmed_balance": "0",
    "unconfirmed_balance": "0"
}
```

## Create BTC and XSN addresses for deposit

Create Segwit addresses for both, `xsn` and `litecoin`
```shell
$ xa-lnd-xsn newaddress np2wkh 
{
        "address": "8pyfqUssgZYXEXT7eqQHCJFrpa8mHGw82v"
}
$ xa-lnd-ltc newaddress np2wkh 
{
        "address": "2N2yE6ZbdxePtco8DuTSQVEJNsNg74KvGd3"
}
```

## Send some money

Send some XSNt (0.2 or more is great) and some LTCt (10 is great) to Exchange A's addresses. Balances should appear in the wallet once the transactions are confirmed.

## Balance after funding the wallets

Query Exchange A wallet balances for both `xsm` and `litecoin` after funding and make sure you see the amount as confirmed balance

```shell
$ xa-lnd-xsn walletbalance
{
    "total_balance": "130000000",
    "confirmed_balance": "130000000",
    "unconfirmed_balance": "0"
}
$ xa-lnd-ltc walletbalance
{
    "total_balance": "1000000000",
    "confirmed_balance": "1000000000",
    "unconfirmed_balance": "0"
}
```

We are now ready with Exchange A's wallets. Exchange B is left with zero balance wallets and this is fine for now. There is no need to separately fund Exchange B's wallets for our PoC, it will get funds via the channel setup [later](/LIGHTNING-03-channels.md). 

We are now ready to move on to the next step and connect our `lnd` instances.

[ [index](/README.md) | [<- previous](/LIGHTNING-00-install.md) / [next ->](/LIGHTNING-02-connect.md) ]
