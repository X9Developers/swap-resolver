[ [index](/README.md) | [<- previous](/LIGHTNING-03-channels.md) ]

# Lightning Cross-Chain Swap

This is the exciting part.

## Balance before the swap

Let's check the XSN and LTC channel balances on Exchange A and B before we execute the swap. In the outputs below, `local_balance` is the balance on Exchange A's side of the channel and `remote_balance` on Exchange B's side.

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


$ xa-lnd-ltc listchannels
{
    "channels": [
        {
            "active": true,
            "remote_pubkey": "021be8d225008d415eaa9b64da60926dfc197e9139deb0f09bae09603352878b5c",
            "channel_point": "e217a814b2f360050b6548acf042c7375d78588aeeee91fc44096d24174dd724:0",
            "chan_id": "757450261840068608",
            "capacity": "10000000",
            "local_balance": "4995475",
            "remote_balance": "5000000",
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

## SWAP

In our example, Exchange A is willing to sell 200 satoshi (0.000002 XSN) for 10000 litoshi (0.0001 LTC).

The command is executed against the `swap-resolver` which controls the `lnd`'s.

```shell
$ cd $GOPATH/src/github.com/ExchangeUnion/swap-resolver

$ ./swap-resolver --rpcserver localhost:7001 takeorder --order_id=123 --maker_amount 200 --maker_coin XSN --taker_amount 10000 --taker_coin=LTC

2018/10/10 16:10:33 Starting takeOrder command -  (*swapresolver.TakeOrderReq)(0xc420099f50)(orderid:"123" taker_amount:10000 taker_coin:LTC maker_amount:200 maker_coin:XSN )
2018/10/10 16:10:34 Swap completed successfully.
  Swap preImage is  32969e4f1a07a77b9468669979006e804501ef4789c599030e7427d14247d322
```

## Balance after the swap

Let's see the impact in the channels. `local_balance` is the balance on Exchange A's side of the channel and `remote_balance` on Exchange B's side.

```shell
$ xa-lnd-xsn listchannels
{
{
    "channels": [
        {
            "active": true,
            "remote_pubkey": "02492ed0c9be232bd2ba82888a8977c3de2a47ff3d6f12428ca860473d37b5ccc8",
            "channel_point": "189148d065850eea0b192f0d8f71a9ec6a1ea2d45b0a3067540d1e4fc892e6bc:0",
            "chan_id": "20905014579036160",
            "capacity": "16000000",
            "local_balance": "14999617",
            "remote_balance": "1000200",
            "commit_fee": "183",
            "commit_weight": "724",
            "fee_per_kw": "253",
            "unsettled_balance": "0",
            "total_satoshis_sent": "200",
            "total_satoshis_received": "0",
            "num_updates": "6",
            "pending_htlcs": [
            ],
            "csv_delay": 1922,
            "private": false
        }
    ]
}
$ xa-lnd-ltc listchannels
{
    "channels": [
        {
            "active": true,
            "remote_pubkey": "021be8d225008d415eaa9b64da60926dfc197e9139deb0f09bae09603352878b5c",
            "channel_point": "e217a814b2f360050b6548acf042c7375d78588aeeee91fc44096d24174dd724:0",
            "chan_id": "757450261840068608",
            "capacity": "10000000",
            "local_balance": "5005475",
            "remote_balance": "4990000",
            "commit_fee": "4525",
            "commit_weight": "724",
            "fee_per_kw": "6250",
            "unsettled_balance": "0",
            "total_satoshis_sent": "0",
            "total_satoshis_received": "10000",
            "num_updates": "2",
            "pending_htlcs": [
            ],
            "csv_delay": 576,
            "private": false
        }
    ]
}
```

As we can see, Exchange A now owns 10000 litoshi more and 200 satoshi less.

Yay!

[ [index](/README.md) | [<- previous](/LIGHTNING-03-channels.md) ]
