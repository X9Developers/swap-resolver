[ [index](/README.md) | [<- previous](/LIGHTNING-02-connect.md) / [next ->](/LIGHTNING-04-swap.md) ]

# Lightning Payment Channels

## Bitcoin

Exchange A opens a bitcoin payment channel to Exchange B and pushes over 0.1 BTC at the same time. Exchange B finally got some bitcoin, yay!

```shell
$ xa-lnd-btc openchannel --node_key=$XB_BTC_PUBKEY --local_amt=16000000 --push_amt=1000000 --sat_per_byte=1000
{
	"funding_txid": "ea64f9bc10b9b81a3cdec7806b50fd0ba2212dee536cae4c3731b9dcdaa7320a"
}

```

The output gives you the `txid` of the funding transaction for the channel. The funding transaction must be confirmed for the channel to be opened. The default number of confirmations is 1.

Until confirmed (which could take a while; testnet...), the pending channels can be seen with the `pendingchannels` command

```shell
$ xa-lnd-btc  pendingchannels
{
    "total_limbo_balance": "0",
    "pending_open_channels": [
        {
            "channel": {
                "remote_node_pub": "022b74059a18bb77c6c906377e92023cc40ae9695920df0d056fc95f135221e69f",
                "channel_point": "246c05860d1589898f9a22317915c935b021a198ae93c1eacc6fc835ca96b5ac:0",
                "capacity": "16000000",
                "local_balance": "12977556",
                "remote_balance": "3000000"
            },
            "confirmation_height": 0,
            "commit_fee": "22444",
            "commit_weight": "724",
            "fee_per_kw": "31000"
        }
    ],
    "pending_closing_channels": [
    ],
    "pending_force_closing_channels": [
    ],
    "waiting_close_channels": [
    ]
}
```

Once the channel is opened, Exchange A lists the bitcoin payment channel as follows

```shell
$ xa-lnd-btc listchannels
{
    "channels": [
        {
            "active": true,
            "remote_pubkey": "022b74059a18bb77c6c906377e92023cc40ae9695920df0d056fc95f135221e69f",
            "channel_point": "246c05860d1589898f9a22317915c935b021a198ae93c1eacc6fc835ca96b5ac:0",
            "chan_id": "1513035751963164672",
            "capacity": "16000000",
            "local_balance": "12977556",
            "remote_balance": "3000000",
            "commit_fee": "22444",
            "commit_weight": "724",
            "fee_per_kw": "31000",
            "unsettled_balance": "0",
            "total_satoshis_sent": "0",
            "total_satoshis_received": "0",
            "num_updates": "0",
            "pending_htlcs": [
            ],
            "csv_delay": 1922,
            "private": false
        }
    ]
}
```



## Litecoin

Exchange A opens a litecoin payment channel to Exchange B and pushes over 0.05 LTC. Exchange B got some litecoin!

```shell
$ xa-lnd-ltc openchannel --node_key=$XB_LTC_PUBKEY --local_amt=10000000 --push_amt=5000000 --sat_per_byte=1000
{
        "funding_txid": "e217a814b2f360050b6548acf042c7375d78588aeeee91fc44096d24174dd724"
}
```

Until confirmed (which could take a while; testnet...), Exchange B lists the new channel as pending channel.
```shell
$ xb-lnd-ltc pendingchannels
{
    "total_limbo_balance": "0",
    "pending_open_channels": [
        {
            "channel": {
                "remote_node_pub": "030dc387fecfe056dcc5ed1043a4d7d1f6b1e712e1a5cc3cc57f2ff031df6ab33e",
                "channel_point": "e217a814b2f360050b6548acf042c7375d78588aeeee91fc44096d24174dd724:0",
                "capacity": "10000000",
                "local_balance": "5000000",
                "remote_balance": "4995475"
            },
            "confirmation_height": 0,
            "commit_fee": "4525",
            "commit_weight": "724",
            "fee_per_kw": "6250"
        }
    ],
    "pending_closing_channels": [
    ],
    "pending_force_closing_channels": [
    ],
    "waiting_close_channels": [
    ]
}
```

Once confirmed, Exchange B lists the litecoin payment channel as follows
```shell
$ xb-lnd-ltc listchannels
{
    "channels": [
        {
            "active": true,
            "remote_pubkey": "030dc387fecfe056dcc5ed1043a4d7d1f6b1e712e1a5cc3cc57f2ff031df6ab33e",
            "channel_point": "e217a814b2f360050b6548acf042c7375d78588aeeee91fc44096d24174dd724:0",
            "chan_id": "757450261840068608",
            "capacity": "10000000",
            "local_balance": "5000000",
            "remote_balance": "4995475",
            "commit_fee": "4525",
            "commit_weight": "724",
            "fee_per_kw": "6250",
            "unsettled_balance": "0",
            "total_satoshis_sent": "0",
            "total_satoshis_received": "0",
            "num_updates": "0",
            "pending_htlcs": [
            ],
            "csv_delay": 576,
            "private": false
        }
    ]
}
```

## Let's move on and swap!

[ [index](/README.md) | [<- previous](/LIGHTNING-02-connect.md) / [next ->](/LIGHTNING-04-swap.md) ]
