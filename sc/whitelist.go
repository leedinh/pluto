package sc

import (
	"github.com/leedinh/pluto/db"
	"github.com/leedinh/pluto/model"
)

type SMWhitelist struct {
	Whitelist map[string]model.ContractInfo `json:"whitelist"`
}

func NewSMWhitelist(db *db.Database) *SMWhitelist {
	wl := &SMWhitelist{}
	wl.Whitelist, _ = model.LoadAllContractInfos(db.Db)
	return wl
}

func (wl *SMWhitelist) RefreshWhiteList(db *db.Database) error {
	wl.Whitelist, _ = model.LoadAllContractInfos(db.Db)
	return nil
}
