if [ $# -eq 0 ]
  then
    echo "No arguments supplied, supported modes to swap are: btc_ltc, btc_xsn, ltc_xsn"
    exit 1
fi

cd $GOPATH/src/github.com/ExchangeUnion/swap-resolver

case "$1" in
  "btc_ltc")
    ./resolver --listen localhost:7002 --peer localhost:7001 --lnd-rpc-ltc localhost:20001 --lnd-rpc-btc localhost:20002
    ;;
  "btc_xsn")
    ./resolver --listen localhost:7002 --peer localhost:7001 --lnd-rpc-btc localhost:20002 --lnd-rpc-xsn localhost:20003
    ;;
  "ltc_xsn")
    ./resolver --listen localhost:7002 --peer localhost:7001 --lnd-rpc-ltc localhost:20001 --lnd-rpc-xsn localhost:20003
    ;;
  *)
    echo "Wrong argument passed, supported modes to swap are: btc_ltc, btc_xsn, ltc_xsn"
    exit 1
    ;;
esac
