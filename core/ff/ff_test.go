package ff

import "testing"

func TestFriendlyHitIsHalved(t *testing.T) {
	l := &Ledger{Threshold: 100}
	got, reason := l.Apply("e1", Hit{Attacker: "a", Target: "b", Team: "red", TargetTeam: "red", Amount: 10})
	if reason != "" || got != 5 {
		t.Fatalf("友伤应减半：%d %q", got, reason)
	}
}

func TestZeroAmountIsZero(t *testing.T) {
	l := &Ledger{Threshold: 100}
	if got, _ := l.Apply("e1", Hit{Attacker: "a", Target: "b", Team: "red", TargetTeam: "red", Amount: 1}); got != 0 {
		t.Fatalf("1 点伤害减半后向下取整为 0：%d", got)
	}
}

func TestDifferentAttackerKeepsOwnRecord(t *testing.T) {
	l := &Ledger{Threshold: 100}
	l.Apply("e1", Hit{Attacker: "a", Target: "b", Team: "red", TargetTeam: "red", Amount: 10})
	got, _ := l.Apply("e2", Hit{Attacker: "c", Target: "b", Team: "red", TargetTeam: "red", Amount: 10})
	if got != 5 {
		t.Fatalf("换一个攻击者应各记各的：%d", got)
	}
}

func TestEmptyLedgerHasNoPairs(t *testing.T) {
	l := &Ledger{Threshold: 100}
	if len(l.pairs) != 0 {
		t.Fatalf("初始应无记录：%v", l.pairs)
	}
}
