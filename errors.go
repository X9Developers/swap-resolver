package main

import (
	"errors"
	"fmt"
)

/*
 * Resolver errors
 */
var (
	ErrorResolverNoDealFound    = errors.New("no deal found")
	ErrorResolverDealValidation = errors.New("deal validation error")
)

type ErrorResolverPayment struct {
	amt  int64
	coin string
	err  error
}

func (e *ErrorResolverPayment) Error() string {
	return fmt.Sprintf(
		"error on payment [amount: %d] [coin: %s] by taker, err: %+v \n", e.amt, e.coin, e.err)
}

func NewErrorPayment(amt int64, coin string, err error) error {
	return &ErrorResolverPayment{amt, coin, err}
}

type ErrorResolverSendPayment struct {
	ErrorResolverPayment
}

func (e *ErrorResolverSendPayment) Error() string {
	return fmt.Sprintf(
		"error sending payment [amount: %d] [coin: %s] by taker, err: %+v \n",
		e.amt, e.coin, e.err,
	)
}

func NewErrorSendPayment(amt int64, coin string, err error) error {
	return &ErrorResolverSendPayment{ErrorResolverPayment{amt, coin, err}}
}

/*
 * P2P Server errors
 */
var (
	ErrorP2PServerNoDealFound = errors.New("no deal found")
)

type ErrorP2PServerGetInfo struct {
	err   error
	chain string
}

func (e *ErrorP2PServerGetInfo) Error() string {
	return fmt.Sprintf(
		"unable to get information for [chain: %s], err: %+v \n ", e.chain, e.err,
	)
}

func NewErrorP2PServerGetInfoError(chain string, err error) error {
	return &ErrorP2PServerGetInfo{err, chain}
}

type ErrorP2PServerSuggestDeal struct {
	err error
}

func (e *ErrorP2PServerSuggestDeal) Error() string {
	return fmt.Sprintf("suggest deal failed, err: %+v \n ", e.err)
}

func NewErrP2PServerSuggestDealError(err error) error {
	return &ErrorP2PServerSuggestDeal{err}
}

type ErrorP2PServerSwap struct {
	id  string
	err error
}

func (e *ErrorP2PServerSwap) Error() string {
	return fmt.Sprintf("swap failed [makerDealId: %s], err: %+v \n ", e.id, e.err)
}

func NewErrorP2PServerSwapError(id string, err error) error {
	return &ErrorP2PServerSwap{id, err}
}

type ErrorP2PServerCreatePreImage struct {
	err error
}

func (e *ErrorP2PServerCreatePreImage) Error() string {
	return fmt.Sprintf("unable to create pre-image, err: %+v \n ", e.err)
}

func NewErrorP2PServerCreatePreImageError(err error) error {
	return &ErrorP2PServerCreatePreImage{err}
}

type ErrorP2PServerPayment struct {
	ErrorResolverPayment
}

func NewErrorP2PServerPayment(amt int64, coin string, err error) error {
	return &ErrorP2PServerPayment{ErrorResolverPayment{amt, coin, err}}
}

type ErrorP2PServerSendPayment struct {
	ErrorResolverPayment
}

func NewErrorP2PServerSendPayment(amt int64, coin string, err error) error {
	return &ErrorP2PServerSendPayment{ErrorResolverPayment{amt, coin, err}}
}
