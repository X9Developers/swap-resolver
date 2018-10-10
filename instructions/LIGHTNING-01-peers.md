[ [index](/README.md) | [<- previous](/LIGHTNING-00-install.md) / [next ->](/LIGHTNING-02-connect.md) ]

# Lightning nodes Setup

Once we installed all components and the Bitcoin & Litecoin testnet blockchains are synced, we can set up the `lnd` and `swap-resolver` processes. We are going to setup two `swap-resolver` processes, simulating two exchanges A & B running `xud` and four `lnd` processes, one for BTC and one for LTC for each exchange. It takes some time until `lnd` synced with the `btcd` and `ltcd` daemons. It is OK to setup the `lnd`s and `swap-resolver`s in parallel.

## Aliases
To make our life easier with `lncli`, we recommend setting up aliases. Add the following to to `~/.bash_profile` or `~/.profile` and source it:

```bash
#Adding lncli aliases
alias xa-lnd-btc='lncli --network testnet --rpcserver=localhost:10002 --no-macaroons'
alias xa-lnd-ltc='lncli --network testnet --rpcserver=localhost:10001 --no-macaroons'
alias xb-lnd-btc='lncli --network testnet --rpcserver=localhost:20002 --no-macaroons'
alias xb-lnd-ltc='lncli --network testnet --rpcserver=localhost:20001 --no-macaroons'
```

Now we can use these aliases to communicate with the 4 `lnd` processes without the need to type long CLI arguments.

## Startup Scripts
To make life even easier, we find the following directory structure in `$GOPATH/src/github.com/ExchangeUnion/swap-resolver`:

*	exchange-a
	+	lnd (resolve.conf)
		*	btc (start.bash)
		*	ltc (start.bash)
	+	xud (start.bash)
*	exchange-b
	+	lnd (resolve.conf)
		*	btc (start.bash)
		*	ltc (start.bash)
	+	xud (start.bash)

The `start.bash` script invokes the LND process using the right parameters (ports, etc). The `resolve.conf` is needed for the swap-resolver to function. Just FYI, no need to do anything for now.

## Exchange A
### Launch `lnd-btc`
Open a terminal to set Exchange A's `lnd-btc` daemon
```shell
cd $GOPATH/src/github.com/ExchangeUnion/swap-resolver/exchange-a/lnd/btc/
./start.bash
```

check progress with
```shell
xa-lnd-btc getinfo
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
cd $GOPATH/src/github.com/ExchangeUnion/swap-resolver/exchange-a/xud/
./start.bash
```


## Exchange B
### Launch `lnd-btc`
Open a terminal to set Exchange B's `lnd-btc` daemon
```shell
cd $GOPATH/src/github.com/ExchangeUnion/swap-resolver/exchange-b/lnd/btc/
./start.bash
```

check progress with
```shell
xb-lnd-btc getinfo
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
Open a terminal to set Exchange B's `xud` daemon
```shell
cd $GOPATH/src/github.com/ExchangeUnion/swap-resolver/exchange-b/xud/
./start.bash
```


## Coffee time v2

Give the four `lnd`s some time to sync with `btcd` and `ltcd`. You can check the status by using the `getinfo` command (use the cli terminal for this). You would want to see `"synced_to_chain": true,` for all four `lnd`s.

### Check status 

Example - Check the status of Exchange A's `lnd-btc`:

```shell
xa-lnd-btc getinfo
{
    "identity_pubkey": "035e3d3884e9a26dcd238b2b2d1ef608a31888365fed327b4af563671a2ee49bb6",
    "alias": "035e3d3884e9a26dcd23",
    "num_pending_channels": 0,
    "num_active_channels": 1,
    "num_peers": 1,
    "block_height": 1356729,
    "block_hash": "00000000007dd1ef9454e52b3986b939ed91b7f81587c95ff0921977be316114",
    "synced_to_chain": true,
    "testnet": true,
    "chains": [
        "bitcoin"
    ],
    "uris": [
    ],
    "best_header_timestamp": "1533285232",
    "version": "0.4.2-beta commit=6a8368a632d22762299686f7ea94eb777a5887c8"
}
```


# Fund Exchange A

## Balance after creating

Query Exchange A wallet balances for both `Bitcoin` and `Litecoin` after creation (we expect to see zeros!)
```shell
$ xa-lnd-btc walletbalance
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

## Create BTC and LTC addresses for deposit

Create Segwit addresses for both, `bitcoin` and `litecoin`
```shell
$ xa-lnd-btc newaddress np2wkh 
{
        "address": "2NBRavXuXmbd73tRvgCAHhTiSuPoa3LqKdd"
}
$ xa-lnd-ltc newaddress np2wkh 
{
        "address": "2N2yE6ZbdxePtco8DuTSQVEJNsNg74KvGd3"
}
```

## Send some money

Send some BTCt (0.2 or more is great) and some LTCt (10 is great) to Exchange A's addresses. You can get some for free via [bitcoin testnet faucets](README.bitcoin.md/#bitcoin-testnet-faucet) and [litecoin testnet faucets](README.litecoin.md/#litecoin-testnet-faucet). Balances should appear in the wallet once the transactions are confirmed.

## Balance after funding the wallets

Query Exchange A wallet balances for both `bitcoin` and `litecoin` after funding and make sure you see the amount as confirmed balance

```shell
$ xa-lnd-btc walletbalance
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
