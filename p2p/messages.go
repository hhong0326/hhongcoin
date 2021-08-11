package p2p

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hhong0326/hhongcoin/blockchain"
	"github.com/hhong0326/hhongcoin/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota // auto const type number setting
	MessageAllBlocksRequest
	MessageAllBlocksResponse
	MessageNewBlockNotify
	MessageNewTxNotify
	MessageNewPeerNotify
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func makeMessage(kind MessageKind, payload interface{}) []byte {

	m := Message{
		Kind:    kind,
		Payload: utils.ToJSON(payload),
	}

	return utils.ToJSON(m)
}

func sendNewestBlock(p *peer) {

	b, err := blockchain.FindBlock(blockchain.BlockChain().NewestHash)
	utils.HandleErr(err)

	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m
}

func requestAllBlocks(p *peer) {

	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	blocks := blockchain.Blocks(blockchain.BlockChain())

	m := makeMessage(MessageAllBlocksResponse, blocks)
	p.inbox <- m
}

func notifyNewBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageNewBlockNotify, b)
	p.inbox <- m
}

func notifyNewTx(tx *blockchain.Tx, p *peer) {
	m := makeMessage(MessageNewTxNotify, tx)
	p.inbox <- m
}

func notifyNewPeer(address string, p *peer) {
	m := makeMessage(MessageNewPeerNotify, address)
	p.inbox <- m
}

func handleMsg(m *Message, p *peer) { // who send

	switch m.Kind {

	case MessageNewestBlock:
		fmt.Printf("Received the newest block from %s\n", p.key)

		var payload blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		b, err := blockchain.FindBlock(blockchain.BlockChain().NewestHash)
		utils.HandleErr(err)
		if payload.Height >= b.Height {
			// request all blocks from 4000
			fmt.Printf("Requesting all blocks from %s\n", p.key)
			requestAllBlocks(p)
		} else {
			// send 4000 our blocks
			fmt.Printf("Sending newest block to %s\n", p.key)
			sendNewestBlock(p) // to
		}
	case MessageAllBlocksRequest:
		fmt.Printf("%s wants all the blocks\n", p.key)
		sendAllBlocks(p)

	case MessageAllBlocksResponse:
		fmt.Printf("Received all the blocks from %s\n", p.key)

		var payload []*blockchain.Block // big size for pointer
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))

		blockchain.BlockChain().Replace(payload)

	case MessageNewBlockNotify:
		fmt.Printf("Notify a new block from %s\n", p.key)

		var payload *blockchain.Block // big size for pointer
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))

		blockchain.BlockChain().AddPeerBlock(payload)

	case MessageNewTxNotify:
		var payload *blockchain.Tx // big size for pointer
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))

		blockchain.Mempool().AddPeerTx(payload)
	case MessageNewPeerNotify:
		var payload string
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))

		parts := strings.Split(payload, ":")
		AddPeer(parts[0], parts[1], parts[2], false)
	}
}
