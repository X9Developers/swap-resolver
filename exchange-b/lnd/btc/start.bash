~/go/bin/lnd --noseedbackup --rpclisten=localhost:20002 --listen=localhost:20012 --restlisten=9002 --datadir=data --logdir=logs  --nobootstrap --no-macaroons --xsncoin.active --xsncoin.testnet  --xsncoin.node=xsnd --xsnd.rpcuser=ross --xsnd.rpcpass=ross --debuglevel=debug --xsnd.zmqpubrawtx=tcp://127.0.0.1:28445 --xsnd.zmqpubrawblock=tcp://127.0.0.1:28444 --alias="Exchange B BTC on 20002/20012"
