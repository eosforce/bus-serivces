package cfg

import (
	"errors"

	"github.com/cihub/seelog"

	eos "github.com/eosforce/goforceio"
	"github.com/fanyang1988/force-go/config"
)

// Relayer transfer, watcher and checker
type Relayer struct {
	SideAccount  eos.PermissionLevel
	RelayAccount eos.PermissionLevel
}

// RelayCfg cfg for relay
type RelayCfg struct {
	Chain         string `json:"chain"`
	RelayContract string `json:"relaycontract"`
}

var relayCfg RelayCfg

// ChainCfgs cfg for each chain
var chainCfgs map[string]*config.Config

var transfers []Relayer
var watchers []Relayer

// GetChainCfg get chain cfg
func GetChainCfg(name string) *config.Config {
	c, ok := chainCfgs[name]
	if !ok || c == nil {
		panic(errors.New("no find chain cfg "))
	}
	return c
}

// GetTransfers get transfers
func GetTransfers() []Relayer {
	return transfers
}

// GetRelayCfg get cfg for relay
func GetRelayCfg() RelayCfg {
	return relayCfg
}

// GetWatchers get watchers
func GetWatchers() []Relayer {
	return watchers
}

// LoadCfgs load cfg for force-relay
func LoadCfgs(path string) error {
	cfgInFile := struct {
		Chains []struct {
			Name string            `json:"name"`
			Cfg  config.ConfigData `json:"cfg"`
		} `json:"chains"`
		Transfer []struct {
			SideAcc  string `json:"sideacc"`
			RelayAcc string `json:"relayacc"`
		} `json:"transfer"`
		Watcher []struct {
			SideAcc  string `json:"sideacc"`
			RelayAcc string `json:"relayacc"`
		} `json:"watcher"`
		Relay RelayCfg `json:"relay"`
	}{}

	err := config.LoadJSONFile(path, &cfgInFile)
	if err != nil {
		return err
	}

	chainCfgs = make(map[string]*config.Config)
	for _, c := range cfgInFile.Chains {
		cc := config.Config{}
		err := cc.Parse(&c.Cfg)
		seelog.Tracef("load cfg %v", cc)
		if err != nil {
			return err
		}
		chainCfgs[c.Name] = &cc
	}

	for _, t := range cfgInFile.Transfer {
		transfers = append(transfers, Relayer{
			SideAccount: eos.PermissionLevel{
				Actor:      eos.AN(t.SideAcc),
				Permission: eos.PN("owner"),
			},
			RelayAccount: eos.PermissionLevel{
				Actor:      eos.AN(t.RelayAcc),
				Permission: eos.PN("owner"),
			},
		})
	}

	for _, t := range cfgInFile.Watcher {
		watchers = append(watchers, Relayer{
			SideAccount: eos.PermissionLevel{
				Actor:      eos.AN(t.SideAcc),
				Permission: eos.PN("owner"),
			},
			RelayAccount: eos.PermissionLevel{
				Actor:      eos.AN(t.RelayAcc),
				Permission: eos.PN("owner"),
			},
		})
	}

	relayCfg = cfgInFile.Relay

	return nil
}
