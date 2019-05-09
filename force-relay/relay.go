package main

import (
	"errors"

	"github.com/eosforce/bus-service/force-relay/logger"

	"github.com/eosforce/bus-service/force-relay/cfg"
	"github.com/eosforce/bus-service/force-relay/chainhandler"
	"github.com/eosforce/bus-service/force-relay/relay"
	"github.com/eosforce/bus-service/force-relay/side"
	"github.com/fanyang1988/force-block-ev/blockdb"
	"github.com/fanyang1988/force-go/p2p"
)

func startRelayService() {
	// from relay to side, so create side client
	_, p2ps := cfg.GetChainCfg("relay")
	chainTyp := cfg.GetChainTyp("relay")

	// for chain id
	info, err := side.Client().GetInfoData()
	if err != nil {
		panic(errors.New("get info err"))
	}

	logger.Debugf("get info %v", *info)

	p2pPeers := p2p.NewP2PClient(chainTyp, p2p.P2PInitParams{
		Name:          "watcher",
		ClientID:      info.ChainID.String(),
		StartBlockNum: 0,
		Peers:         p2ps,
		Logger:        logger.Logger(),
	})

	p2pPeers.RegHandler(&handlerImp{
		verifier: blockdb.NewFastBlockVerifier(p2ps, 0, chainhandler.NewChainHandler(
			func(block *chainhandler.Block, actions []chainhandler.Action) {
				relay.HandRelayBlock(block, actions)
			}, chainTyp)),
	})
	p2pPeers.Start()
}
