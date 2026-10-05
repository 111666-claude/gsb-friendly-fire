// Command ffcli 结算一批命中并打印每条的实际伤害。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"example.com/ff-core/ff"
)

type line struct {
	ID     string `json:"id"`
	Damage int    `json:"damage"`
	Reason string `json:"reason"`
}

func scenario(name string) (*ff.Ledger, []struct {
	ID  string
	Hit ff.Hit
}) {
	type item = struct {
		ID  string
		Hit ff.Hit
	}
	switch name {
	case "cross-team":
		return &ff.Ledger{Threshold: 100}, []item{
			{ID: "e1", Hit: ff.Hit{Attacker: "a", Target: "b", Team: "red", TargetTeam: "blue", Amount: 40}},
		}
	case "exempt":
		return &ff.Ledger{Threshold: 100}, []item{
			{ID: "e1", Hit: ff.Hit{Attacker: "a", Target: "b", Team: "red", TargetTeam: "red", Amount: 40, Exempt: true}},
		}
	case "cap":
		return &ff.Ledger{Threshold: 10}, []item{
			{ID: "e1", Hit: ff.Hit{Attacker: "a", Target: "b", Team: "red", TargetTeam: "red", Amount: 40}},
			{ID: "e2", Hit: ff.Hit{Attacker: "a", Target: "b", Team: "red", TargetTeam: "red", Amount: 40}},
		}
	case "work":
		l := &ff.Ledger{Threshold: 1000000}
		items := make([]item, 0, 3000)
		for i := 0; i < 3000; i++ {
			items = append(items, item{
				ID: fmt.Sprintf("e%d", i),
				Hit: ff.Hit{
					Attacker: fmt.Sprintf("a%d", i%500), Target: fmt.Sprintf("t%d", i%500),
					Team: "red", TargetTeam: "red", Amount: 10,
				},
			})
		}
		return l, items
	}
	fmt.Fprintf(os.Stderr, "未知场景 %q\n", name)
	os.Exit(2)
	return nil, nil
}

func main() {
	sample := flag.String("sample", "cross-team", "场景名")
	flag.Parse()

	l, items := scenario(*sample)
	for _, it := range items {
		damage, reason := l.Apply(it.ID, it.Hit)
		if *sample != "work" {
			out, _ := json.Marshal(line{ID: it.ID, Damage: damage, Reason: reason})
			fmt.Println(string(out))
		}
	}
	if *sample == "work" {
		out, _ := json.Marshal(map[string]int{"items": len(items), "work": l.Work})
		fmt.Println(string(out))
	}
}
