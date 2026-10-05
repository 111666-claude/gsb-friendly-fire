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
	if l.seen[id] {
		return 0, "duplicate"
	}
	l.seen[id] = true

	if h.Exempt {
		return 0, "exempt"
	}

	if h.Team != h.TargetTeam {
		return h.Amount, ""
	}

	damage := h.Amount / 2

	key := pairKey(h.Attacker, h.Target)
	used := l.pairs[key]
	l.Work++
	if used+damage > l.Threshold {
		damage = l.Threshold - used
		if damage < 0 {
			damage = 0
		}
	}
	l.pairs[key] = used + damage
	return damage, ""
}
