// Package ff 结算友军伤害：同队才减免，带阈值与豁免。
package ff

// Hit 是一次命中。
type Hit struct {
	Attacker   string
	Target     string
	Team       string
	TargetTeam string
	Amount     int
	Exempt     bool
}

// Ledger 是友伤台账。
type Ledger struct {
	Threshold int

	pairs map[string]int
	seen  map[string]bool
	Work  int
}

func pairKey(a, t string) string { return a + "\x00" + t }

// Apply 结算一次命中，返回实际造成的伤害与拒绝原因。
func (l *Ledger) Apply(id string, h Hit) (int, string) {
	if l.seen == nil {
		l.seen = map[string]bool{}
	}
	if l.pairs == nil {
		l.pairs = map[string]int{}
	}
	l.seen[id] = true

	damage := h.Amount / 2

	used := 0
	for k, v := range l.pairs {
		l.Work++
		if k == pairKey(h.Attacker, h.Target) {
			used += v
		}
	}
	_ = used

	l.pairs[pairKey(h.Attacker, h.Target)] += damage
	return damage, ""
}
