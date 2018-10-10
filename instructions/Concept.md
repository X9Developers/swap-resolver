# A new approach to cross-chain swaps on lightning

Let's assume there are 2 players A and B which want to swap BTC for LTC instantly and without involvement of any third party. A agrees to exchange x BTC for y LTC from B. Since there is no trust between the parties, they can't do this in two separate transactions; it must be done in a single atomic swap operation on lightning.

## Current Approach

* A, the originator (maker) creates a pre-image and hashes it (via invoice)
* A creates a two legs route: A to B, B to A. On the first leg A commits to sending x BTC to B providing that B knows the pre-image. On the second leg, B commits to sending y LTC to A providing A knows the preimage.
* Since A has the pre-image it can get the y LTC from B.
* B can now use the pre-image to get the x BTC from A

### LND - the current approach
The current approach is to enhance LND to support the swap natively. This was done as a PoC by Conner as part of the [swapz branch](https://github.com/ExchangeUnion/lnd-swap/tree/swapz).
LND was enhanced to support, in parallel both LTC and BTC and be able to route the payment request between them. 
LND was also enhanced to support a ticker command option, a multi-currency wallet, etc (not published)

### Weaknesses of the current approach
While doable, there are some weaknesses in the current approach:
* Implementing multi-currency swap with LND implies a lot of changes across all functional area of the project. This can be done in phases, but it implies slow progress.
* The current solution assumes that A and B have direct BTC and LTC payment channels. This is not a good restriction to have. The solution can be enhanced to remove the restriction but this will bring extra complication.
* Swap rate - this is currently hardcoded into LND code. Clearly, this can't be used in real life where A and B should agree on the conversion rate via an external system for every swap.
* Swap and multi-chain support are on LND's roadmap but are not considered high on the priority list. This is why understandably any pull requests related to multi-chain won't be prioritized. 

## New Approach
To overcome these shortcomings, we propose an improved solution which can mostly be implemented outside of LND. The new solution only requires small changes to LND which potentially make sense for non-multi-chain/swap use cases, too.

### LND - the NEW approach
The only change required to LND is the following:

When a payment request arrives at an LND (payee) and LND can't find the hash in its database, it can query an external system, like [XUD](https://github.com/ExchangeUnion/xud/), for the pre-image via RPC call.

### The setup
* LTCD - connected to the litecoin chain
* BTCD - connected to the bitcoin chain
* An instance of LND configured for bitcoin and using the BTCD instance (`a-lnd-btc`, `b-lnd-btc`)
* An instance of LND configured for litecoin and using the LTCD instance  (`a-lnd-ltc`, `b-lnd-ltc`)
* An instance of XUD which is aware of the LNDs processes (`a-XUD`, `b-XUD`)

### How the new swap works
1. `a-XUD` and `b-XUD` agree to execute a swap of x BTC for y LTC via an external order exchange & matching mechanism.
2. `a-XUD` creates a preimage and hashes it. The hash is shared with `b-XUD` and kept in both XUDs together with the swap information.
3. `a-XUD` requests its `a-lnd-btc` to create a route and transfer x BTC to `b-lnd-btc`. "I will pay you x BTC if you know the preimage of the hash".
4. `b-lnd-btc` does not know the preimage for the incoming transfer, so it informs `b-XUD` via RPC call.
5. `b-XUD` is aware of the negotiated trade and the hash so it requests its `b-lnd-ltc` to create a route and initiate a transfer of y LTC from its `b-lnd-ltc` to `a-lnd-ltc` using `a-XUD`'s provided hash: " I will pay you y LTC if you know the preimage of the hash".
6. `a-lnd-ltc` does not know the preimage, so it informs `a-XUD`.
7. `a-XUD` is the creator of the preimage and provides it to `a-lnd-ltc` which completes the LTC transaction by sending it to `b-lnd-ltc`.
8. `b-lnd-ltc` forwards the preimage to `b-XUD` which in turn provides it to `b-lnd-btc`, which completes the transaction by sending it to `a-lnd-btc`.
9. `a-XUD` receives the preimage and the swap is completed.

### Notes
- A and B do not need to have direct payment channels between each other. They can use multi-hop routes. The routes can be different for the LTC and BTC payment and are created separately.
- Discovery, order management & matching are done in a separate software layer and are not relevant for the order execution part which is described in this document.
- For step 2., we are exploring sharing the hash directly within the HTLC instead of via a separate communication channel, which would also ommit a 'received hash' response.
- The described swap approach makes partial order filling simpler
- The described swap approach abstracts swap logic from the lightning implementation and thus makes integration of other lightning implementations and also e.g. raiden feasible

[Here](/README.md) a complete guide along with the necessary tools to do the described swap yourself.
