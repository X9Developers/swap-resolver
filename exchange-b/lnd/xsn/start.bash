lnd --noseedbackup --rpclisten=localhost:20003 --listen=localhost:20013 --restlisten=9003 --datadir=data --logdir=logs  --nobootstrap --no-macaroons --xsncoin.active --xsncoin.testnet --xsncoin.node=xsnd --xsnd.rpcuser=xu --xsnd.rpcpass=xu --xsnd.zmqpubrawblock=tcp://127.0.0.1:28333 --xsnd.zmqpubrawtx=tcp://127.0.0.1:28332 --debuglevel=debug --alias="Exchange B XSN on 20003/20013"
