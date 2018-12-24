package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"sync"

	"crypto/sha256"
	"time"

	"github.com/ExchangeUnion/lnd/lnrpc"
	pbp2p "github.com/ExchangeUnion/swap-resolver/swapp2p"
	"github.com/davecgh/go-spew/spew"
	"github.com/dchest/uniuri"
	"golang.org/x/net/context"
)

type P2PServer struct {
	xuPeer pbp2p.P2PClient
	lnLTC  lnrpc.LightningClient
	lnBTC  lnrpc.LightningClient
	lnXSN  lnrpc.LightningClient
	mu     sync.Mutex // protects data structure
}

// TakeOrder is called to initiate a swap between maker and taker
// it is a temporary service needed until the integration with XUD
// intended to be called from CLI to simulate order taking by taker
func (s *P2PServer) TakeOrder(ctx context.Context, req *pbp2p.TakeOrderReq) (*pbp2p.TakeOrderResp, error) {

	log.Printf("TakeOrder (maker) starting with [request: %+v] \n ", spew.Sdump(req))

	ctxt, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var info *lnrpc.GetInfoResponse
	var err error

	switch req.TakerCoin {
	case pbp2p.CoinType_BTC:
		info, err = s.lnBTC.GetInfo(ctxt, &lnrpc.GetInfoRequest{})
		if err != nil {
			return nil, NewErrorP2PServerGetInfoError("BTC LND", err)
		}

	case pbp2p.CoinType_LTC:
		info, err = s.lnLTC.GetInfo(ctxt, &lnrpc.GetInfoRequest{})
		if err != nil {
			return nil, NewErrorP2PServerGetInfoError("LTC LND", err)
		}

	case pbp2p.CoinType_XSN:
		info, err = s.lnXSN.GetInfo(ctxt, &lnrpc.GetInfoRequest{})
		if err != nil {
			return nil, NewErrorP2PServerGetInfoError("XSN LND", err)
		}

	}
	log.Printf("dump ln-info [coin: %s] getinfo: %v \n", req.TakerCoin.String(), spew.Sdump(info))
	spew.Sdump(info)

	newDeal := deal{
		role:        Taker,
		orderId:     req.Orderid,
		takerDealId: uniuri.New(),
		takerAmount: req.TakerAmount,
		takerCoin:   req.TakerCoin,
		takerPubKey: info.IdentityPubkey,
		makerAmount: req.MakerAmount,
		makerCoin:   req.MakerCoin,
		createTime:  time.Now(),
	}

	log.Printf("suggesting deal to peer. deal data: %+v \n", spew.Sdump(newDeal))

	suggestDealResp, err := s.xuPeer.SuggestDeal(ctx, &pbp2p.SuggestDealReq{
		Orderid:     newDeal.orderId,
		TakerCoin:   newDeal.takerCoin,
		TakerAmount: newDeal.takerAmount,
		TakerDealId: newDeal.takerDealId,
		TakerPubkey: newDeal.takerPubKey,
		MakerCoin:   newDeal.makerCoin,
		MakerAmount: newDeal.makerAmount,
	})
	if err != nil {
		return nil, NewErrP2PServerSuggestDealError(err)
	}

	newDeal.makerDealId = suggestDealResp.MakerDealId
	newDeal.makerPubKey = suggestDealResp.MakerPubkey
	copy(newDeal.hash[:], suggestDealResp.RHash[:32])

	log.Printf("deal agreed with maker. deal data %+v: \n", spew.Sdump(newDeal))

	// @TODO: enable mutex to protect data
	deals = append(deals, &newDeal)

	newDeal.executeTime = time.Now()

	swapResp, err := s.xuPeer.Swap(ctx, &pbp2p.SwapReq{
		MakerDealId: suggestDealResp.MakerDealId,
	})
	if err != nil {
		return nil, NewErrorP2PServerSwapError(suggestDealResp.MakerDealId, err)
	}

	ret := &pbp2p.TakeOrderResp{
		RPreimage: swapResp.RPreimage,
	}
	return ret, nil
}

// SuggestDeal is called by the taker to inform the maker that he
// would like to execute a swap. The maker may reject the request
// for now, the maker can only accept/reject and can't rediscuss the
// deal or suggest partial amount. If accepted the maker should respond
// with a hash that would be used for teh swap.
func (s *P2PServer) SuggestDeal(ctx context.Context, req *pbp2p.SuggestDealReq) (*pbp2p.SuggestDealResp, error) {

	log.Printf("SuggestDeal (taker) stating with %v: ", spew.Sdump(req))

	ctxt, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var info *lnrpc.GetInfoResponse
	var err error

	switch req.MakerCoin {
	case pbp2p.CoinType_BTC:
		info, err = s.lnBTC.GetInfo(ctxt, &lnrpc.GetInfoRequest{})
		if err != nil {
			return nil, NewErrorP2PServerGetInfoError("BTC LND", err)
		}

	case pbp2p.CoinType_LTC:
		info, err = s.lnLTC.GetInfo(ctxt, &lnrpc.GetInfoRequest{})
		if err != nil {
			return nil, NewErrorP2PServerGetInfoError("LTC LND", err)
		}

	case pbp2p.CoinType_XSN:
		info, err = s.lnXSN.GetInfo(ctxt, &lnrpc.GetInfoRequest{})
		if err != nil {
			return nil, NewErrorP2PServerGetInfoError("XSN LND", err)
		}

	}
	log.Printf("dump ln-info [coin: %s] getinfo: %v \n", req.TakerCoin.String(), spew.Sdump(info))
	spew.Sdump(info)

	newDeal := deal{
		role:        Maker,
		orderId:     req.Orderid,
		takerDealId: req.TakerDealId,
		takerAmount: req.TakerAmount,
		takerCoin:   req.TakerCoin,
		takerPubKey: req.TakerPubkey,
		makerPubKey: info.IdentityPubkey,
		makerAmount: req.MakerAmount,
		makerCoin:   req.MakerCoin,
		makerDealId: uniuri.New(),
		hash: [32]byte{
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
		},
		createTime: time.Now(),
	}

	// create preImage and a hash for the deal
	if _, err := rand.Read(newDeal.preImage[:]); err != nil {
		return nil, NewErrorP2PServerCreatePreImageError(err)
	}

	newDeal.hash = sha256.Sum256(newDeal.preImage[:])

	// @TODO: enable mutex to protect data
	deals = append(deals, &newDeal)

	ret := pbp2p.SuggestDealResp{
		Orderid:     newDeal.orderId,
		RHash:       newDeal.hash[:],
		MakerDealId: newDeal.makerDealId,
		MakerPubkey: newDeal.makerPubKey,
	}
	return &ret, nil
}

// Swap initiates the swap. It is called by the taker to confirm that
// he has the hash and confirm the deal.
func (s *P2PServer) Swap(ctx context.Context, req *pbp2p.SwapReq) (*pbp2p.SwapResp, error) {
	var deal *deal

	log.Printf("Swap (maker) starting with: %v ", spew.Sdump(req))

	for _, d := range deals {
		if d.makerDealId == req.MakerDealId {
			deal = d
			break
		}
	}

	if deal == nil {
		return nil, ErrorP2PServerNoDealFound
	}

	cmdLnd := s.lnLTC

	switch deal.makerCoin {
	case pbp2p.CoinType_BTC:
	case pbp2p.CoinType_LTC:
		cmdLnd = s.lnXSN
	case pbp2p.CoinType_XSN:
		cmdLnd = s.lnLTC
	}

	resp, err := cmdLnd.SendPaymentSync(context.Background(), &lnrpc.SendRequest{
		DestString:  deal.takerPubKey,
		Amt:         deal.takerAmount,
		PaymentHash: deal.hash[:],
	})

	if err != nil {
		return nil, NewErrorP2PServerSendPayment(deal.takerAmount, deal.takerCoin.String(), err)
	}

	if resp.PaymentError != "" {
		return nil, NewErrorP2PServerPayment(deal.takerAmount, deal.takerCoin.String(), fmt.Errorf(resp.PaymentError))
	}

	log.Printf("dumping sendPayment response : %+v \n", spew.Sdump(resp))

	ret := &pbp2p.SwapResp{
		RPreimage: deal.preImage[:],
	}
	return ret, nil
}

func newP2PServer(xuPeer pbp2p.P2PClient,
	lnLTC lnrpc.LightningClient,
	lnBTC lnrpc.LightningClient,
	lnXSN lnrpc.LightningClient,
) *P2PServer {
	s := &P2PServer{
		xuPeer: xuPeer,
		lnLTC:  lnLTC,
		lnBTC:  lnBTC,
		lnXSN:  lnXSN,
	}
	return s
}
