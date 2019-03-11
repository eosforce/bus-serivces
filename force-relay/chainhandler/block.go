package chainhandler

import (
	"encoding/binary"

	eos "github.com/eosforce/goeosforce"
)

type Block struct {
	Producer         eos.AccountName `json:"producer"`
	Num              uint32          `json:"num"`
	ID               eos.Checksum256 `json:"id"`
	Previous         eos.Checksum256 `json:"previous"`
	Confirmed        uint16          `json:"confirmed"`
	TransactionMRoot eos.Checksum256 `json:"transaction_mroot"`
	ActionMRoot      eos.Checksum256 `json:"action_mroot"`
	MRoot            eos.Checksum256 `json:"mroot"`
}

func (b *Block) GetNum() uint32 {
	return BlockID2Num(b.ID)
}

func BlockID2Num(blockID []byte) uint32 {
	if len(blockID) < 32 {
		return 0
	}
	return binary.BigEndian.Uint32(blockID[:32])
}

type Action struct {
	Account       eos.AccountName   `json:"account"`
	Name          eos.ActionName    `json:"name"`
	Authorization []PermissionLevel `json:"authorization"`
	Data          []byte            `json:"data"`
}

type PermissionLevel struct {
	Actor      eos.AccountName    `json:"actor"`
	Permission eos.PermissionName `json:"permission"`
}