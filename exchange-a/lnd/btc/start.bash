~/go/bin/lnd --noseedbackup --rpclisten=localhost:10002 --listen=localhost:10012 --restlisten=8002 --datadir=data --logdir=logs  --nobootstrap --no-macaroons --xsncoin.active --xsncoin.testnet --xsncoin.node=xsnd --xsnd.rpcuser=ross --xsnd.rpcpass=ross --debuglevel=debug --xsnd.zmqpubrawtx=tcp://127.0.0.1:28445 --xsnd.zmqpubrawblock=tcp://127.0.0.1:28444 --alias="Exchange A BTC on 10002/10012"
