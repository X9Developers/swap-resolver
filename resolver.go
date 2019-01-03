package main

// TODO: seperate resolver RPC into a standalone package
// TODO: properly handle the time locks - Important!!!!

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"sync"

	"os"
	"path/filepath"
	"time"

	"github.com/ExchangeUnion/lnd/lnrpc"
	pb "github.com/ExchangeUnion/lnd/lnrpc"
	pbp2p "github.com/ExchangeUnion/swap-resolver/swapp2p"
	"github.com/btcsuite/btcutil"
	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defaultTLSCertFilename  = "tls.cert"
	defaultMacaroonFilename = "admin.macaroon"
	defaultRpcPort          = "10009"
	defaultRpcHostPort      = "localhost:" + defaultRpcPort

	Maker role = iota + 1
	Taker
)

var (
	//tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	//certFile   = flag.String("cert_file", "", "The TLS cert file")
	//keyFile    = flag.String("key_file", "", "The TLS key file")
	//port       = flag.Int("port", 10000, "The server port")

	//Commit stores the current commit hash of this build. This should be
	//set using -ldflags during compilation.
	Commit string

	deals []*deal

	defaultLndDir       = btcutil.AppDataDir("lnd", false)
	defaultTLSCertPath  = filepath.Join(defaultLndDir, defaultTLSCertFilename)
	defaultMacaroonPath = filepath.Join(defaultLndDir, defaultMacaroonFilename)
)

type role int

type deal struct {
	// Maker or Taker ?
	role role

	// global order it in XU network
	orderId string

	// takerCoin is the name of the coin the taker is expecting to get
	takerCoin   pbp2p.CoinType
	takerDealId string
	takerAmount int64
	takerPubKey string

	// makerCoin is the name of the coin the maker is expecting to get
	makerCoin   pbp2p.CoinType
	makerDealId string
	makerAmount int64
	makerPubKey string

	hash     [32]byte
	preImage [32]byte

	createTime  time.Time
	executeTime time.Time
}

func (d *deal) isTaker() bool {
	return d.role == Taker
}

func (d *deal) isValid() bool {
	//@TODO: implement amount validation ("check that I got the right amount before sending out the agreed amount")
	return true
}

type hashResolverServer struct {
	p2pServer *P2PServer
	mu        sync.Mutex // protects data structure
}

// ResolveHash retrieves a deal by the given request hash and handles the payment with Maker or Taker logic
func (s *hashResolverServer) ResolveHash(ctx context.Context, req *pb.ResolveRequest) (*pb.ResolveResponse, error) {
	log.Printf(" ResolveHash starting with [hash: %s] [amount: %d] \n", req.Hash, req.Amount)

	var deal *deal
	for _, d := range deals {
		if hex.EncodeToString(d.hash[:]) == req.Hash {
			deal = d
			break
		}
	}

	if deal == nil {
		log.Printf("unable to find deal [request hash: %s] \n", req.Hash)
		return nil, ErrorResolverNoDealFound
	}

	if !deal.isValid() {
		log.Printf("validation error [request-hash: %s] \n", req.Hash)
		return nil, ErrorResolverDealValidation
	}

	if deal.isTaker() {
		log.Printf("executing taker code")
		return s.resolveTaker(deal)
	}

	log.Printf("executing maker code")
	return s.resolveMaker(deal)
}

//resolveTaker will forward the payment to other chains
func (s *hashResolverServer) resolveTaker(deal *deal) (*pb.ResolveResponse, error) {
	cmdLnd := s.p2pServer.lnBTC

	switch deal.makerCoin {
	case pbp2p.CoinType_BTC:
	case pbp2p.CoinType_XSN:
		cmdLnd = s.p2pServer.lnXSN
	case pbp2p.CoinType_LTC:
		cmdLnd = s.p2pServer.lnLTC
	}

	resp, err := cmdLnd.SendPaymentSync(context.Background(), &lnrpc.SendRequest{
		DestString:  deal.makerPubKey,
		Amt:         deal.makerAmount,
		PaymentHash: deal.hash[:],
	})
	if err != nil {
		err = NewErrorSendPayment(deal.makerAmount, deal.makerCoin.String(), err)
		log.Printf(err.Error())
		return nil, err
	}
	if resp.PaymentError != "" {
		err = NewErrorPayment(deal.makerAmount, deal.makerCoin.String(), fmt.Errorf(resp.PaymentError))
		log.Printf(err.Error())
		return nil, err
	}

	log.Printf("dumping response from maker to taker: %+v \n", spew.Sdump(resp))

	return &pb.ResolveResponse{
		Preimage: hex.EncodeToString(resp.PaymentPreimage[:]),
	}, nil
}

func (s *hashResolverServer) resolveMaker(d *deal) (*pb.ResolveResponse, error) {
	return &pb.ResolveResponse{
		Preimage: hex.EncodeToString(d.preImage[:]),
	}, nil
}

func newHashResolverServer(p2pServer *P2PServer) *hashResolverServer {
	s := &hashResolverServer{
		p2pServer: p2pServer,
	}
	return s
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "[lncli] %v\n", err)
	os.Exit(1)
}

func getClient(ctx *cli.Context) (pbp2p.P2PClient, func()) {
	conn := getPeerConn(ctx, false)
	cleanUp := func() {
		conn.Close()
	}
	return pbp2p.NewP2PClient(conn), cleanUp
}

func getPeerConn(ctx *cli.Context, skipMacaroons bool) *grpc.ClientConn {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(ctx.GlobalString("peer"), opts...)
	if err != nil {
		fatal(err)
	}
	return conn
}

func main() {
	app := cli.NewApp()
	app.Name = "resolver"
	app.Version = fmt.Sprintf("%s commit=%s", "0.0.1", Commit)
	app.Usage = "Use me to simulate order taking by the taker"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen",
			Value: "localhost:10000",
			Usage: "An ip:port to listen for peer connections",
		},
		cli.StringFlag{
			Name:  "peer",
			Value: "a.b.com:10000",
			Usage: "A host:port of peer xu daemon",
		},
		cli.StringFlag{
			Name:  "lnd-rpc-ltc",
			Value: "localhost:10001",
			Usage: "RPC host:port of LND connected to LTC chain",
		},
		cli.StringFlag{
			Name:  "lnd-rpc-btc",
			Value: "localhost:10002",
			Usage: "RPC host:port of LND connected to LTC chain",
		},
		cli.StringFlag{
			Name:  "lnd-rpc-xsn",
			Value: "localhost:10003",
			Usage: "RPC host:port of LND connected to XSN chain",
		},
	}

	app.Action = func(c *cli.Context) error {
		log.Printf("Server starting")

		lnLTC, close := getLNDClient(c, "lnd-rpc-ltc")
		defer close()

		lnBTC, close := getLNDClient(c, "lnd-rpc-btc")
		defer close()

		lnXSN, close := getLNDClient(c, "lnd-rpc-xsn")
		defer close()

		xuPeer, close := getClient(c)
		defer close()

		log.Printf("Got peer connection")

		listen := c.GlobalString("listen")
		lis, err := net.Listen("tcp", listen)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("listening on %s", listen)

		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)
		p2pServer := newP2PServer(xuPeer, lnLTC, lnBTC, lnXSN)
		pb.RegisterHashResolverServer(grpcServer, newHashResolverServer(p2pServer))
		pbp2p.RegisterP2PServer(grpcServer, p2pServer)

		log.Printf("Server ready")
		grpcServer.Serve(lis)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fatal(err)
	}
}

func getLNDClient(ctx *cli.Context, name string) (lnrpc.LightningClient, func()) {
	conn := getLNDClientConn(ctx, false, name)

	cleanUp := func() {
		conn.Close()
	}
	return lnrpc.NewLightningClient(conn), cleanUp
}

func getLNDClientConn(ctx *cli.Context, skipMacaroons bool, name string) *grpc.ClientConn {
	creds, err := credentials.NewClientTLSFromFile(defaultTLSCertPath, "")
	if err != nil {
		fatal(err)
	}

	// Create a dial options array.
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(ctx.String(name), opts...)
	if err != nil {
		fatal(err)
	}
	return conn
}
