# Phase 6C / Phase 7 扩展：IM 实测指标与容量分析

## 1. 证据范围与报告口径

本报告保存已验证的 CI 结果、失败证据、研究结论和待测项，补充 [Phase 7 源码对比](phase7-im-source-ci-comparison.md)。全部使用独立临时 GitHub Ubuntu Runner、公开上游源码及合成账号；不修改 GlobalChat 正式产品。数据来自下列运行日志，而不是厂商宣传。各产品协议、数据库、拓扑、发送路径、客户端及负载不完全一致，不能据此排名。

| 产品 | 固定源码 SHA | 环境 |
|---|---|---|
| WuKongIM | 21d905a328bb7141167bc6dbe20fdb65588e0420 | 三节点；原生协议；2/3副本对比 |
| JuggleIM | 0a47433c1aed6f08b48a3bfcf6b80dfe0a55a6c4 | 社区版单节点 + MySQL；临时修正初始化 SQL；真实 WebSocket |
| OpenIM | a101d85ec16cc21bec60005169725adc1ae5649a | 源码单进程、MongoDB/Redis/MinIO；队列改 Redis，非默认 Kafka 生产配置 |

## 2. 已验证的性能指标

### JuggleIM 20 客户端，目标 40 publishes/s，测量 20s

| 指标 | 私聊 | 群聊 |
|---|---:|---:|
| 连接成功 | 20/20 | 20/20 |
| ACK 成功 / 失败 | 799 / 0 | 800 / 0 |
| ACK 吞吐 (/s) | 39.95 | 40 |
| ACK P50 / P95 / P99 (ms) | 0.71 / 1.23 / 1.69 | 20.09 / 78.25 / 116.41 |
| 接收 P50 / P95 / P99 (ms) | 2.82 / 4.04 / 6.49 | 2285.85 / 3799.97 / 4073.54 |
| 接收回调计数 | 799 | 46960 |
| 回调吞吐 (/s) | 39.95 | 2348 |

证据：[run 46019003789 / job 140591400888](https://github.com/trade3403/-globalchat-im-research/actions/runs/46019003789/job/140591400888)。群聊回调统计存在口径问题：20个客户端、800次发布却有46960次回调，必须按消息唯一 ID 与目标接收用户去重、独立构造期望集合，再判断丢失/重复；不能把回调量视为唯一有效投递。群聊延迟显著高于私聊是本轮观察，瓶颈原因尚未验证。

### JuggleIM 10000 客户端，目标 40 publishes/s，测量 60s

| 指标 | 实测 |
|---|---:|
| 成功建立连接 | 10000/10000 |
| 连接延迟 P95 | 9.95 ms |
| ACK 成功 / 失败 | 2399 / 0 |
| ACK 吞吐 | 39.98/s |
| ACK P95 / P99 | 2.55 / 3.54 ms |
| 接收回调 | 2399 |
| 接收 P95 / P99 | 4.72 / 7.96 ms |
| 容器重启 | 0 |

证据：[run 46026543674 / job 140599282727](https://github.com/trade3403/-globalchat-im-research/actions/runs/46026543674/job/140599282727)。通过的是该轮压测断言，未独立逐连接采样整个测量期的同时在线数，也未证明万人同时发消息或零丢包。

### WuKongIM：三节点双向消息与 2000 QPS 副本测试

三节点跨节点双向真实消息通过：[job 140548070409](https://github.com/trade3403/-globalchat-im-research/actions/runs/46011555775/job/140548070409)。副本测试每种配置各三轮、每轮约10秒。

| 指标 | 2副本 | 3副本 |
|---|---:|---:|
| 各轮成功发送/接收 | 20000, 20000, 19999 | 20000, 20000, 19996 |
| 实际吞吐范围 (/s) | 1999.9–2000 | 1999.6–2000 |
| 发送/接收错误 | 各轮0 | 各轮0 |
| ACK P50（各轮中位数，ms） | 18.20 | 38.74 |
| ACK P99（各轮中位数，ms） | 76.89 | 100.29 |
| ACK P99（各轮，ms） | 61.23, 76.89, 217.99 | 96.25, 130.32, 100.29 |
| 测试后平均 RSS (MiB) | 866.0 | 950.2 |
| 平均 RSS 增量 (MiB) | 232.6 | 279.7 |

证据：[run 46019003789 / job 140591349913](https://github.com/trade3403/-globalchat-im-research/actions/runs/46019003789/job/140591349913)。三副本 P99 的各轮中位数约高30%，但双副本存在217.99ms长尾异常；三轮不足以证明稳定优势。RSS 为原测试输出口径，非单连接内存成本；该测试是消息吞吐而非连接数压测。

### OpenIM：两用户 REST 发送与历史读取

| 指标 | 实测 |
|---|---:|
| 发送 | 50 |
| 接收方历史中唯一消息 | 50 |
| 缺失 | 0 |
| 最后序号 | 50 |
| REST 请求 P50 / P95 / P99 | 3.948 / 8.069 / 10.195 ms |
| 越权历史读取 | 被拒绝 |
| 原始重复计数 | 47 |

证据：[run 46019003789 / job 140591352978](https://github.com/trade3403/-globalchat-im-research/actions/runs/46019003789/job/140591352978)；[脚本](../scripts/im-source-runtime.py)。脚本多次拉取相同序号区间并累计 occurrence_count，所以47不是已证实的存储或投递重复缺陷。应分别核对单次历史响应中的唯一 ID、跨轮重复读取和实时推送。REST 延迟不等于 WebSocket 端到端延迟；越权只覆盖一个负向场景。

## 3. 容量测试失败与根因排查

| 产品 | 目标连接 | 结果 | 证据 |
|---|---:|---|---|
| JuggleIM | 10000 | 通过，10000 次成功连接，60秒低负载测量 | [job 140599282727](https://github.com/trade3403/-globalchat-im-research/actions/runs/46026543674/job/140599282727) |
| JuggleIM | 20000 | 失败，10903 次成功连接；IM 容器重启5次 | [job 140599282987](https://github.com/trade3403/-globalchat-im-research/actions/runs/46026543674/job/140599282987) |
| JuggleIM | 100000 | 失败，10904 次成功连接；连接重置及拒绝，初次 nofile 调整失败 | [job 140595100695](https://github.com/trade3403/-globalchat-im-research/actions/runs/46025317632/job/140595100695) |
| WuKongIM | 100000 | 失败，模拟器处于初始化/重试，采样 active 峰值0；不代表实际服务器只能支持0连接 | [job 140595101032](https://github.com/trade3403/-globalchat-im-research/actions/runs/46025317632/job/140595101032) |
| OpenIM | 100000 | 尚未运行；配置允许100000不等于验证达到100000 | [固定版本 gateway 配置](https://github.com/openimsdk/open-im-server/blob/a101d85ec16cc21bec60005169725adc1ae5649a/config/openim-msggateway.yml) |

JuggleIM 2万复测已成功提高客户端 nofile 至262144，仍在约10900次成功连接附近失败，MySQL健康且未获得明确OOM证据。最终容器状态 OOMKilled=false 不能排除历史OOM。容器重启掩盖退出原因，后续关闭自动重启、保留 pipeline 退出码、Docker events、完整日志及内核信息。成功连接总数不是同时在线峰值，不得称10900为架构极限。

原始 artifacts：[100k](https://github.com/trade3403/-globalchat-im-research/actions/runs/46025317632)；[10k/20k](https://github.com/trade3403/-globalchat-im-research/actions/runs/46026543674)。原始产物可能过期，永久报告应记录必要摘要及证据链接，不写入账号/token。

## 4. 当前阶梯测试

用户要求先测试1万、3万、5万，再挑战10万。以下是已触发的测试，不是完成结果。快照取自本次 GitHub 查询。

| 任务 | 目标 | 快照状态 |
|---|---|---|
| 根因诊断 | WuKongIM 1000/10000、JuggleIM 12000 | [46029886214](https://github.com/trade3403/-globalchat-im-research/actions/runs/46029886214) 运行中；Juggle 子任务 PID 参数错误，已在新梯度脚本修正 |
| WuKongIM / JuggleIM | 各10000、30000、50000 | [46030019543](https://github.com/trade3403/-globalchat-im-research/actions/runs/46030019543) 排队中 |
| OpenIM | 20客户端校验及10000、30000、50000 | [46030318969](https://github.com/trade3403/-globalchat-im-research/actions/runs/46030318969) 排队中 |

工作流：[两家阶梯](../.github/workflows/connection-ladder.yml)；[OpenIM 阶梯](../.github/workflows/openim-connection-ladder.yml)；[OpenIM 客户端](../scripts/openim-connection-ladder.go)。

OpenIM 新客户端按真实账号、Token、WebSocket 鉴权响应、心跳和在线连接进行测试，REST 发布后解析真实收件端 push 并按 marker 去重。代码尚未通过运行验证；它不是原生客户端发送 ACK 测试。三家当前拓扑、认证模式（WuKongIM token mode none）、消息发送路径不同，仍属探索性测试。测试端和服务端共享 Runner，文件描述符、端口池、CPU、内存均可能成为瓶颈；记录为共享环境实测，不外推生产容量。

## 5. 必须持续记录的核心指标

1. 最大实际同时在线；新建连接/s；连接成功率；登录延迟。
2. 私聊/群聊发布 QPS；唯一目标设备投递/s；ACK P50/P95/P99；端到端 P50/P95/P99。
3. 丢失率、重复率、顺序正确率；重连恢复；离线补拉。
4. CPU、RSS/内存、单连接增量内存、磁盘、网络。
5. 节点故障恢复、重启后数据完整性、24h/168h稳定性、每万在线用户成本。

统一比较须记录源 SHA、依赖版本、拓扑、Runner规格、TLS、存储及ACK语义、消息大小、连接数、发送速率、群成员、测量时间、错误及所有重复轮次。每轮保留失败证据。未测项写未测，不能以 CI 绿灯替代完整语义校验。

## 6. 对 GlobalChat Rust 的研究建议（目标，不是实测）

| 指标 | 初期建议目标 | 长期建议目标 |
|---|---:|---:|
| 单节点连接 | 10000 | 100000 |
| 私聊发布 | 2000 QPS | 20000 QPS |
| 私聊端到端 P99 | ≤50ms | ≤20ms |
| 群聊端到端 P99（需固定群规模） | ≤200ms | ≤100ms |
| 已确认持久化消息丢失 | 0 | 0 |
| 用户可见重复 | 0 | 0 |
| 集群故障恢复 | ≤10s | ≤3s |
| 持续稳定性 | 24h | 168h |

这些指标需结合硬件、持久化与业务 SLA 审核。优先学习 WuKongIM 的提交/投递分离与副本架构、JuggleIM 的业务分层与租户边界、OpenIM 的多端游标同步。当前未做 Rust 同条件运行，不作性能领先声明。

## 7. 报告维护约定

- 每轮 CI 结束后读取 job 日志、断言和可用 artifacts，更新本报告状态与数据；保留历史失败记录及对应链接。
- 将可复现的原始摘要、环境配置、临时补丁和分析结论纳入 Git，较大原始日志留 Actions artifacts；不得提交真实凭据、个人信息或第三方私有源码。
- 更新 RESEARCH.md 目录入口；新增测试必须同时说明验收标准、未覆盖范围及后续行动。
- 新报告以版本化提交保存，不在旧报告中覆盖历史实测值。未执行、排队、运行中、失败、通过分开；CI 通过与产品验收分开。
- 本约定不是自动写回任务；后续分析需实际提交并验证 Git commit。
