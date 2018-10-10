package main

import (
	"log"
	"time"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/ExchangeUnion/swap-resolver/swapp2p"
	"os"
	"fmt"
	"encoding/hex"
	"github.com/davecgh/go-spew/spew"
)

const (
	defaultRpcPort          = "10000"
	defaultRpcHostPort      = "localhost:" + defaultRpcPort
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
		takeordercommand,
	}

	if err := app.Run(os.Args); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "[lncli] %v\n", err)
	os.Exit(1)
}

func getClient(ctx *cli.Context) (pb.P2PClient, func()) {
	conn := getClientConn(ctx, false)

	cleanUp := func() {
		conn.Close()
	}

	return pb.NewP2PClient(conn), cleanUp
}

func getClientConn(ctx *cli.Context, skipMacaroons bool) *grpc.ClientConn {
	var opts []grpc.DialOption
	//if *tls {
	//if *caFile == "" {
	//	*caFile = testdata.Path("ca.pem")
	//}
	//creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
	//if err != nil {
	//	log.Fatalf("Failed to create TLS credentials %v", err)
	//}
	//opts = append(opts, grpc.WithTransportCredentials(creds))
	//} else {
	opts = append(opts, grpc.WithInsecure())
	//}
	conn, err := grpc.Dial(ctx.GlobalString("rpcserver"), opts...)

	if err != nil {
		fatal(err)
	}

	return conn
}


var takeordercommand = cli.Command{
	Name:     "takeorder",
	Category: "Order",
	Usage:    "Instruct resolver to take an order.",
	Description: `
	`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "order_id",
			Usage: "the chain to generate an address for",
		},
		cli.Int64Flag{
			Name:  "maker_amount",
			Usage: "the number of coins denominated in {l/s}atoshis the maker is expecting to get",
		},
		cli.StringFlag{
			Name:  "maker_coin",
			Usage: "the coins which the maker is expecting to get",
		},
		cli.Int64Flag{
			Name:  "taker_amount",
			Usage: "the number of coins denominated in {l/s}atoshis the taker is expecting to get",
		},
		cli.StringFlag{
			Name:  "taker_coin",
			Usage: "the coins which the taker is expecting to get",
		},

	},

	Action: takeOrder,
}

func takeOrder(ctx *cli.Context) error{

	client, cleanUp := getClient(ctx)
	defer cleanUp()

	// Show command help if no arguments provided
	if ctx.NArg() == 0 && ctx.NumFlags() == 0 {
		cli.ShowCommandHelp(ctx, "takeOrder")
		return nil
	}

	req := &pb.TakeOrderReq{
	}

	if !ctx.IsSet("order_id"){
		return fmt.Errorf("order_id argument missing")
	}
	req.Orderid = ctx.String("order_id")

	if !ctx.IsSet("maker_amount"){
		return fmt.Errorf("maker_amount argument missing")
	}
	req.MakerAmount = int64(ctx.Int("maker_amount"))

	if !ctx.IsSet("maker_coin"){
		return fmt.Errorf("maker_coin argument missing")
	}

	switch ctx.String("maker_coin") {
	case "BTC":
		req.MakerCoin = pb.CoinType_BTC
	case "LTC":
		req.MakerCoin = pb.CoinType_LTC
	case "XSN":
                req.MakerCoin = pb.CoinType_XSN
	default:
		return fmt.Errorf("Invalid maker coin %v. Valid values are BTC, LTC and XSN only.", ctx.String("maker_coin"))
	}

	if !ctx.IsSet("taker_amount"){
		return fmt.Errorf("taker_amount argument missing")
	}
	req.TakerAmount = int64(ctx.Int("taker_amount"))

	if !ctx.IsSet("taker_coin"){
		return fmt.Errorf("taker_coin argument missing")
	}
	switch ctx.String("taker_coin") {
	case "BTC":
		req.TakerCoin = pb.CoinType_BTC
	case "LTC":
		req.TakerCoin = pb.CoinType_LTC
	case "XSN":
                req.TakerCoin = pb.CoinType_XSN
	default:
		return fmt.Errorf("Invalid taker coin %v. Valid values are BTC, LTC and XSN only.", ctx.String("taker_coin"))
	}

	if req.MakerCoin == req.TakerCoin{
		return fmt.Errorf("Maker and Taker coins can't be the same")
	}


	log.Printf("Starting takeOrder command -  %s",spew.Sdump(req))
	ctxt, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := client.TakeOrder(ctxt, req)
	if err != nil {
		log.Fatalf("%v.ResolveHash(_) = _, %v: ", client, err)
	}
	log.Printf("Swap completed successfully.\n  Swap preImage is  %v \n",hex.EncodeToString(resp.RPreimage))


	return nil;

}




