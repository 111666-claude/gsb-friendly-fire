# friendly-fire

友军伤害结算，拆成核心库与命令行两个 module，用 `go.work` 串起来。只用 Go 标准库。

```
(cd core && go test ./...)
(cd core && go vet ./...)
(cd cli && go vet ./...)
go run ./cli --sample=cross-team
go run ./cli --sample=exempt
go run ./cli --sample=cap
go run ./cli --sample=work
```

`go.work` 把两个 module 串起来，所以测试与 vet 要在各自的模块目录里跑，不能从仓库根目录 `./...`。

## 口径

- **友伤范围**：只有攻击者与目标的队伍相同才算友伤；不同队伍按全额伤害结算。
- **减免**：友伤伤害等于 `Amount / 2`（向下取整）。
- **豁免覆盖层**：`Exempt` 为真的目标完全免伤，伤害记 0。
- **阈值**：同一 `(attacker, target)` 累计友伤不超过 `Threshold`，超出部分截断。
- **幂等**：同一个 `id` 只结算一次。

## 不变量

- 单条命中结算出的伤害不超过 `Amount`。
- 同一对攻守的累计友伤不超过 `Threshold`。
- `Work` 记判定期间检查过的既有记录数，不得随记录总数放大。

## 输出契约

每条命中打印一行：

```
{"id": <原样>, "damage": <实际伤害>, "reason": <拒绝原因或空串>}
```

`--sample=work` 只打印 `{"items": N, "work": N}`。
