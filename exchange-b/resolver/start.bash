if [ $# -eq 0 ]
  then
    echo "No arguments supplied"
fi

cd ~/go/src/github.com/ExchangeUnion/swap-resolver

case "$1" in
  "btc_ltc")
    go run resolver.go --listen localhost:7002 --peer localhost:7001 --lnd-rpc-ltc localhost:20001 --lnd-rpc-btc localhost:20002
    ;;
  "btc_xsn")
    go run resolver.go --listen localhost:7003 --peer localhost:7002 --lnd-rpc-btc localhost:20002 --lnd-rpc-xsn localhost:20003
    ;;
  "ltc_xsn")
    go run resolver.go --listen localhost:7003 --peer localhost:7001 --lnd-rpc-ltc localhost:20001 --lnd-rpc-xsn localhost:20003
    ;;
  *)
    echo "Wrong argument passed, supported modes to swap are: btc_ltc, btc_xsn, ltc_xsn"
    exit 1
    ;;
esac
