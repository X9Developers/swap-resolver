[ [index](/README.md) | [<- previous](/LIGHTNING-01-peers.md) / [next ->](/LIGHTNING-03-channels.md) ]

# Lightning Peer Connection
In this step we connect Exchange A and Exchange B on network level so they become peers of each other. We create two parallel P2P networks, one for XSN and one for LTC.    

## Establishing Connection
First, we extract Exchange B's pubKeys and bind them to `XB_XSN_PUBKEY` and `XB_LTC_PUBKEY`, so we can use them to set connections and channels

```shell
XB_XSN_PUBKEY=`xb-lnd-xsn getinfo|grep identity_pubkey|cut -d '"' -f 4`
XB_LTC_PUBKEY=`xb-lnd-ltc getinfo|grep identity_pubkey|cut -d '"' -f 4`
```


By using Exchange B's pubKey, host and port number, Exchange A establishes two connections, using the following commands (output is `{}`, that's fine)

```shell
xa-lnd-xsn connect $XB_XSN_PUBKEY@127.0.0.1:20013
xa-lnd-ltc connect $XB_LTC_PUBKEY@127.0.0.1:20011
```

### Exchange A post-connection
```shell
$ xa-lnd-xsn listpeers
{
    "peers": [
        {
            "pub_key": "02492ed0c9be232bd2ba82888a8977c3de2a47ff3d6f12428ca860473d37b5ccc8",
            "address": "127.0.0.1:20013",
            "bytes_sent": "5736",
            "bytes_recv": "5050",
            "sat_sent": "0",
            "sat_recv": "0",
            "inbound": false,
            "ping_time": "288"
        }
    ]
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
