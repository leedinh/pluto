package parser

import (
	"github.com/leedinh/pluto/sc"
)

type Rule struct {
	Name      string
	Condition func(*Transaction) bool
	Action    func()
	Whitelist sc.SMWhitelist
}

func (r *Rule) Check(t *Transaction) bool {
	return r.Condition(t)
}
func FilterAddliquidity(t *Transaction) bool {
	return t.To == "0x7d0556d55ca1a92708681e2e231733ebd922597d"
}

func Init_Addliquidity() Rule {
	return Rule{
		Name:      "Addliquidity",
		Condition: FilterAddliquidity,
		Action:    func() {},
	}
}

func InitRules() []Rule {
	return []Rule{
		Init_Addliquidity(),
	}
}
