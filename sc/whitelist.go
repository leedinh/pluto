package sc

type SMWhitelist struct {
	Whitelist map[string]ContractInfo `json:"whitelist"`
}

func NewSMWhitelist(db *Database) *SMWhitelist {
	wl := &SMWhitelist{}
	wl.Whitelist, _ = db.LoadAllFromDB()
	return wl
}
