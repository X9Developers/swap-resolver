package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/ExchangeUnion/swap-resolver/swapp2p"
	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	//default rpc settings
	defaultRpcPort     = "10000"
	defaultRpcHostPort = "localhost:" + defaultRpcPort

	//coin options for resolver
	CoinBTC = "BTC"
	CoinLTC = "LTC"
	CoinXSN = "XSN"

	//context flags
	ContextOrderId     = "order_id"
	ContextTakerAmount = "taker_amount"
	ContextMakerAmount = "maker_amount"
	ContextTakerCoin   = "taker_coin"
	ContextMakerCoin   = "maker_coin"
)

var (
	//Commit stores the current commit hash of this build. This should be
	//set using -ldflags during compilation.
	Commit string
)

func main() {
	app := cli.NewApp()
	app.Name = "resolver-cli"
	app.Version = fmt.Sprintf("%s commit=%s", "0.0.1", Commit)
	app.Usage = "Use me to simulate order taking by the taker"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "rpcserver",
			Value: defaultRpcHostPort,
			Usage: "host:port of resolver daemon",
		},
	}
	app.Commands = []cli.Command{
		takeOrderCmd,
	}

	if err := app.Run(os.Args); err != nil {
		handleFatal(err)
	}
}

var takeOrderCmd = cli.Command{
	Name:     "takeorder",
	Category: "Order",
	Usage:    "Instruct resolver to take an order.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  ContextOrderId,
			Usage: "the chain to generate an address for",
		},
		cli.Int64Flag{
			Name:  ContextMakerAmount,
			Usage: "the number of coins denominated in {l/s}atoshis the maker is expecting to get",
		},
		cli.StringFlag{
			Name:  ContextMakerCoin,
			Usage: "the coins which the maker is expecting to get",
		},
		cli.Int64Flag{
			Name:  ContextTakerAmount,
			Usage: "the number of coins denominated in {l/s}atoshis the taker is expecting to get",
		},
		cli.StringFlag{
			Name:  ContextTakerCoin,
			Usage: "the coins which the taker is expecting to get",
		},
	},
	Action: takeOrder,
}

//takeOrder enriches a TakeOrderRequest and sends request to the p2p swap client.
func takeOrder(ctx *cli.Context) (err error) {
	client, cleanUp := getClient(ctx)
	defer cleanUp()

	// Show command help if no arguments provided
	if ctx.NArg() == 0 && ctx.NumFlags() == 0 {
		cli.ShowCommandHelp(ctx, "takeOrder")
		return nil
	}

	err = validateFlags(ctx)
	if err != nil {
		return err
	}

	req := &pb.TakeOrderReq{}
	req.Orderid = ctx.String(ContextOrderId)
	req.MakerAmount = int64(ctx.Int(ContextMakerAmount))
	req.TakerAmount = int64(ctx.Int(ContextTakerAmount))
	req.MakerCoin, err = getCoinType(ContextMakerCoin, ctx.String(ContextMakerCoin))
	if err != nil {
		return err
	}

	req.TakerCoin, err = getCoinType(ContextTakerCoin, ctx.String(ContextTakerCoin))
	if err != nil {
		return err
	}

	if req.MakerCoin == req.TakerCoin {
		return fmt.Errorf("maker and taker coin must not be the same ([%s])", req.MakerCoin.String())
	}

	log.Printf("starting takeOrder command -  %s \n", spew.Sdump(req))
	ctxt, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.TakeOrder(ctxt, req)
	if err != nil {
		log.Fatalf("%v.ResolveHash(_) = _, %v: ", client, err)
		return err
	}
	log.Printf("swap completed successfully.\n  Swap preImage is %s \n", hex.EncodeToString(resp.RPreimage))
	return nil
}

//validateFlags checks that all mandatory fields are set
func validateFlags(ctx *cli.Context) (err error) {
	if !ctx.IsSet(ContextOrderId) {
		return fmt.Errorf("%s argument missing", ContextOrderId)
	}
	if !ctx.IsSet(ContextMakerCoin) {
		return fmt.Errorf("%s argument missing", ContextMakerCoin)
	}
	if !ctx.IsSet(ContextTakerAmount) {
		return fmt.Errorf("%s argument missing", ContextTakerAmount)
	}
	if !ctx.IsSet(ContextMakerCoin) {
		return fmt.Errorf("%s argument missing", ContextMakerCoin)
	}
	if !ctx.IsSet(ContextTakerCoin) {
		return fmt.Errorf("%s argument missing", ContextTakerCoin)
	}
	return nil
}

//getCoinType returns CoinType based on the context-value for key maker_coin or taker_coin.
func getCoinType(ctxKey string, ctxVal string) (pb.CoinType, error) {
	switch ctxKey {
	case CoinBTC:
		return pb.CoinType_BTC, nil
	case CoinLTC:
		return pb.CoinType_LTC, nil
	case CoinXSN:
		return pb.CoinType_XSN, nil
	}
	return -1, fmt.Errorf("[key %s, value: %s] not in supported list [BTC, LTC, XSN]", ctxKey, ctxVal)
}

//handleFatal prints error and exits application
func handleFatal(err error) {
	fmt.Fprintf(os.Stderr, "[lncli] %v\n", err)
	os.Exit(1)
}

//getClient returns a swap P2P client
func getClient(ctx *cli.Context) (pb.P2PClient, func()) {
	conn := getClientConn(ctx, false)
	cleanUp := func() {
		conn.Close()
	}
	return pb.NewP2PClient(conn), cleanUp
}

//getClientConn returns a grpc client connection
func getClientConn(ctx *cli.Context, skipMacaroons bool) *grpc.ClientConn {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(ctx.GlobalString("rpcserver"), opts...)
	if err != nil {
		handleFatal(err)
	}
	return conn
}