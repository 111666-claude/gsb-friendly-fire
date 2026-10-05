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
		return 0, "duplicate id"
	}
	l.seen[id] = true

	// 不同队伍不是友伤，按全额伤害结算。
	if h.Team != h.TargetTeam {
		return h.Amount, ""
	}

	// 豁免目标完全免伤。
	if h.Exempt {
		return 0, ""
	}

	// 友伤减半（向下取整）。
	damage := h.Amount / 2

	// 同一对攻守的累计友伤不超过阈值，超出部分截断。
	key := pairKey(h.Attacker, h.Target)
	used := l.pairs[key]
	l.Work++
	if remaining := l.Threshold - used; damage > remaining {
		damage = remaining
	}
	if damage < 0 {
		damage = 0
	}

	l.pairs[key] += damage
	return damage, ""
}
