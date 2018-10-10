[ [index](/README.md) | [<- previous](/LIGHTNING-01-peers.md) / [next ->](/LIGHTNING-03-channels.md) ]

# Lightning Peer Connection
In this step we connect Exchange A and Exchange B on network level so they become peers of each other. We create two parallel P2P networks, one for BTC and one for LTC.    

## Establishing Connection
First, we extract Exchange B's pubKeys and bind them to `XB_BTC_PUBKEY` and `XB_LTC_PUBKEY`, so we can use them to set connections and channels

```shell
XB_BTC_PUBKEY=`xb-lnd-btc getinfo|grep identity_pubkey|cut -d '"' -f 4`
XB_LTC_PUBKEY=`xb-lnd-ltc getinfo|grep identity_pubkey|cut -d '"' -f 4`
```


By using Exchange B's pubKey, host and port number, Exchange A establishes two connections, using the following commands (output is `{}`, that's fine)

```shell
xa-lnd-btc connect $XB_BTC_PUBKEY@127.0.0.1:20012
xa-lnd-ltc connect $XB_LTC_PUBKEY@127.0.0.1:20011
```

### Exchange A post-connection
```shell
$ xa-lnd-btc listpeers
{
    "peers": [
        {
            "pub_key": "022b74059a18bb77c6c906377e92023cc40ae9695920df0d056fc95f135221e69f",
            "address": "127.0.0.1:20012",
            "bytes_sent": "279",
            "bytes_recv": "279",
            "sat_sent": "0",
            "sat_recv": "0",
            "inbound": false,
            "ping_time": "0"
        }
    ]
}

$ xa-lnd-ltc listpeers
{
    "peers": [
        {
            "pub_key": "021be8d225008d415eaa9b64da60926dfc197e9139deb0f09bae09603352878b5c",
            "address": "127.0.0.1:20011",
            "bytes_sent": "137",
            "bytes_recv": "137",
            "sat_sent": "0",
            "sat_recv": "0",
            "inbound": false,
            "ping_time": "0"
        }
    ]
}

```

We are no ready to set up payment channels. 

[ [index](/README.md) | [<- previous](/LIGHTNING-01-peers.md) / [next ->](/LIGHTNING-03-channels.md) ]