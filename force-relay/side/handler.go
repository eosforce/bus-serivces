package side

import (
	"github.com/eosforce/bus-service/force-relay/chainhandler"
	"github.com/eosforce/bus-service/force-relay/logger"
)

var lastCommittedBlockNum uint32

// HandSideBlock handle block from side chain
func HandSideBlock(block *chainhandler.Block, actions []chainhandler.Action) {
	if handSideBlockImp(block, actions) {
		//time.Sleep(5 * time.Millisecond)
	}
}

// HandSideBlock handle block from side chain
func handSideBlockImp(block *chainhandler.Block, actions []chainhandler.Action) bool {
	const retryTimes int = 32

	num := block.GetNum()

	var blockCommitLast *chainhandler.Block
	var err error

	if lastCommittedBlockNum > 0 && num != 0 && lastCommittedBlockNum >= num {
		logger.Debugf("no need commit %v to %v", lastCommittedBlockNum, num)
		return false
	}

	if lastCommittedBlockNum == 0 {
		blockCommitLast, err = GetLastCommittedBlock()
		if err != nil {
			logger.LogError("get last commit block err", err)
		}
		lastCommittedBlockNum = blockCommitLast.GetNum()
	}

	if blockCommitLast != nil && num != 0 && lastCommittedBlockNum >= num {
		logger.Debugf("no need commit %v to %v", lastCommittedBlockNum, num)
		return false
	}

	commitWorkerMng.OnBlock(block, actions)

	return true
}
