package basic

import (
	//"fmt"
	force_relay_commit "github.com/eosforce/bus-service/force_relay_commit"
	eos "github.com/eosforce/goeosforce"
)

type block struct {
	Producer         eos.AccountName `json:"producer"`
	ID               eos.Checksum256 `json:"id"`
	Previous         eos.Checksum256 `json:"previous"`
	Confirmed        uint16          `json:"confirmed"`
	TransactionMRoot eos.Checksum256 `json:"transaction_mroot"`
	ActionMRoot      eos.Checksum256 `json:"action_mroot"`
	MRoot            eos.Checksum256 `json:"mroot"`
}

//链的名称是可配置的		transfer是发送者的名称		会给发送者奖励	目前应该没有用  可配置
type CommitParam struct {
	Name     eos.Name        `json:"chain"`
	Transfer eos.AccountName `json:"transfer"`
	Block    block           `json:"block"`
	Actions  []action        `json:"actions"`
}

type action struct {
	Account       eos.AccountName   `json:"account"`
	Name          eos.ActionName    `json:"name"`
	Authorization []permissionLevel `json:"authorization"`
	Data          []byte            `json:"data"`
}

type permissionLevel struct {
	Actor      eos.AccountName    `json:"actor"`
	Permission eos.PermissionName `json:"permission"`
}

var chain eos.Name
var transfer eos.AccountName

func SetChain(chainname string) {
	chain = eos.Name(chainname)
}

func SetTransfer(transfername string) {
	transfer = eos.AccountName(transfername)
}

func newCommitAction(relayblock *force_relay_commit.RelayBlock, Action []*force_relay_commit.RelayAction) *eos.Action {
	b := block{
		Producer:         eos.AN(relayblock.Producer),
		ID:               relayblock.Id,
		Previous:         relayblock.Previous,
		Confirmed:        uint16(relayblock.Confirmed),
		TransactionMRoot: relayblock.TransactionMroot,
		ActionMRoot:      relayblock.ActionMroot,
		MRoot:            relayblock.Mroot,
	}

	acts := make([]action, 0, 8)[:]
	for _, ActionValue := range Action {
		auth := make([]permissionLevel, 0, 8)[:]
		for _, authori := range ActionValue.Authorization {
			auth = append(auth, permissionLevel{
				Actor:      eos.AN(authori.Actor),
				Permission: eos.PN(authori.Permission),
			})
		}
		acts = append(acts, action{
			Account:       eos.AN(ActionValue.Account),
			Name:          eos.ActN(ActionValue.ActionName),
			Authorization: auth[:],
			Data:          ActionValue.Data[:],
		})
	}

	//chain := eos.Name("eosforce")
	return &eos.Action{
		Account: eos.AN("force.relay"),
		Name:    eos.ActN("commit"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AccountName(chain), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(CommitParam{
			Name:     chain,
			Transfer: transfer,
			Block:    b,
			Actions:  acts,
		}),
	}
}
