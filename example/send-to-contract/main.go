package main

import (
	"context"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"log"
	"math/rand"
	"strings"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

func main() {
	client := liteclient.NewConnectionPool()

	// connect to mainnet lite server
	err := client.AddConnection(context.Background(), "135.181.140.212:13206", "K0t3+IWLOXHYMvMcrGZDPs+pn58a17LFbnXoQkKc2xw=")
	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return
	}

	api := ton.NewAPIClient(client)

	// seed words of account, you can generate them with any wallet or using wallet.NewSeed() method
	words := strings.Split("birth pattern then forest walnut then phrase walnut fan pumpkin pattern then cluster blossom verify then forest velvet pond fiction pattern collect then then", " ")

	w, err := wallet.FromSeed(api, words, wallet.V3)
	if err != nil {
		log.Fatalln("FromSeed err:", err.Error())
		return
	}

	log.Println("wallet address:", w.Address())

	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		log.Fatalln("CurrentMasterchainInfo err:", err.Error())
		return
	}

	balance, err := w.GetBalance(context.Background(), block)
	if err != nil {
		log.Fatalln("GetBalance err:", err.Error())
		return
	}

	if balance.NanoTON().Uint64() >= 3000000 {
		// create transaction body cell, depends on what contract needs, just random example here
		body := cell.BeginCell().
			MustStoreUInt(0x123abc55, 32). // op code
			MustStoreUInt(rand.Uint64(), 64). // query id
			// payload:
			MustStoreAddr(address.MustParseAddr("EQAbMQzuuGiCne0R7QEj9nrXsjM7gNjeVmrlBZouyC-SCLlO")).
			MustStoreRef(
				cell.BeginCell().
					MustStoreBigCoins(tlb.MustFromTON("1.521").NanoTON()).
					EndCell(),
			).EndCell()

		log.Println("sending transaction and waiting for confirmation...")

		err = w.Send(context.Background(), &wallet.Message{
			Mode: 1, // pay fees separately (from balance, not from amount)
			InternalMessage: &tlb.InternalMessage{
				Bounce:  true, // return amount in case of processing error
				DstAddr: address.MustParseAddr("EQCD39VS5jcptHL8vMjEXrzGaRcCVYto7HUn4bpAOg8xqB2N"),
				Amount:  tlb.MustFromTON("0.03"),
				Body:    body,
			},
		}, true)
		if err != nil {
			log.Fatalln("Send err:", err.Error())
			return
		}

		return
	}

	log.Println("not enough balance:", balance.TON())
}
