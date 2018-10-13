[ [index](/README.md) | [<- previous](/instructions/LIGHTNING-02-connect.md) / [next ->](/instructions/LIGHTNING-04-swap.md) ]

# Lightning Payment Channels

## XSN

Exchange A opens a xsn payment channel to Exchange B and pushes over 0.1 XSN at the same time. Exchange B finally got some coins, yay!

```shell
$ xa-lnd-xsn openchannel --node_key=$XB_XSN_PUBKEY --local_amt=16000000 --push_amt=1000000 --sat_per_byte=1000
{
	"funding_txid": "ea64f9bc10b9b81a3cdec7806b50fd0ba2212dee536cae4c3731b9dcdaa7320a"
}

```

The output gives you the `txid` of the funding transaction for the channel. The funding transaction must be confirmed for the channel to be opened. The default number of confirmations is 1.

Until confirmed (which could take a while; testnet...), the pending channels can be seen with the `pendingchannels` command

```shell
$ xa-lnd-xsn  pendingchannels
{
    "total_limbo_balance": "0",
    "pending_open_channels": [
        {
            "channel": {
                "remote_node_pub": "02492ed0c9be232bd2ba82888a8977c3de2a47ff3d6f12428ca860473d37b5ccc8",
                "channel_point": "189148d065850eea0b192f0d8f71a9ec6a1ea2d45b0a3067540d1e4fc892e6bc:0",
                "capacity": "16000000",
                "local_balance": "14999817",
                "remote_balance": "1000000"
            },
            "confirmation_height": 0,
            "commit_fee": "183",
            "commit_weight": "724",
            "fee_per_kw": "253"
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
$ xa-lnd-xsn listchannels
{
    "channels": [
        {
            "active": true,
            "remote_pubkey": "03e634d7505a8c2840700bbf38a696409e4ec73e59175f2b24abf9df8014fd2016",
            "channel_point": "b8f58dd932368fe66e1dcbbdce7327a86bba4a3eadb37b0ada3206a3da8fe589:0",
            "chan_id": "20530081113964544",
            "capacity": "16000000",
            "local_balance": "14999817",
            "remote_balance": "1000000",
            "commit_fee": "183",
            "commit_weight": "724",
            "fee_per_kw": "253",
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

[ [index](/README.md) | [<- previous](/instructions/LIGHTNING-02-connect.md) / [next ->](/instructions/LIGHTNING-04-swap.md) ]
