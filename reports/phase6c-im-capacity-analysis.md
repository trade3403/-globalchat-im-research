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


## 8. 阶梯压测完成后的增量分析（2026-11-28）

以下补充保留第4节历史状态快照，取代其作为当前运行状态的依据。三个工作流均已结束且总体 failure；通过的子任务单独记录。实际查阅各 job 日志，不把启动失败、注册失败、鉴权失败解释为容量极限。

| 目标连接 | WuKongIM | JuggleIM | OpenIM |
|---:|---|---|---|
| 1000 | 鉴权诊断未通过 | 本轮未测 | 本轮未测 |
| 10000 | 鉴权失败，无有效在线采样 | 通过，10000/10000 | 第一批注册失败 |
| 20000 | 未测 | 历史复测10903次连接成功，未通过 | 未测 |
| 30000 | 鉴权失败，无有效在线采样 | 服务未就绪，尚未进入连接阶段 | 第一批注册失败 |
| 50000 | 鉴权失败，无有效在线采样 | 服务未就绪，尚未进入连接阶段 | 第一批注册失败 |
| 100000 | 历史尝试未进入有效连接采样 | 历史尝试10904次连接成功，未通过 | 未测 |

### 8.1 JuggleIM 万人连接重复验证

| 指标 | run 46026543674 | run 46030019543 |
|---|---:|---:|
| 连接成功 | 10000/10000 | 10000/10000 |
| 测量时间 | 60s | 60s |
| 目标发布速率 | 40/s | 40/s |
| 连接 P95 (ms) | 9.95 | 4.47 |
| ACK 成功/失败 | 2399/0 | 2399/0 |
| ACK 吞吐 (/s) | 39.98 | 39.98 |
| ACK P95/P99 (ms) | 2.55/3.54 | 1.90/2.22 |
| 接收回调计数 | 2399 | 2399 |
| 接收 P95/P99 (ms) | 4.72/7.96 | 3.47/3.97 |
| IM 容器重启 | 0 | 0 |

新轮证据：[job 140611144676](https://github.com/trade3403/-globalchat-im-research/actions/runs/46030019543/job/140611144676)。两轮均通过脚本验收，但只测低发布速率；没有持续逐连接在线计数、按预期接收集合验收或统一硬件复测，不能据此保证万人同时发消息、绝对零丢失或认定新版本性能提升。

3万和5万均出现 `Local stack did not become ready: no response`，IM 容器 ExitCode=1、OOMKilled=false、restart=0，启动约0.25/0.29秒后退出，MySQL healthy。与历史2万/10万在约10900次连接成功后的异常不同，这两轮没有进入连接阶段。关闭自动重启揭示退出状态，但尚未确定启动失败的具体代码路径。不可填报3万/5万实测连接上限。
- [3万 job 140611144794](https://github.com/trade3403/-globalchat-im-research/actions/runs/46030019543/job/140611144794)
- [5万 job 140611144855](https://github.com/trade3403/-globalchat-im-research/actions/runs/46030019543/job/140611144855)

### 8.2 WuKongIM：鉴权阻断

三个阶梯任务反复报告 `client: connack reason=ReasonAuthFail`，active采样为0。启动初期也有 healthz connection refused，之后主要为鉴权失败。不是容量饱和证据。保留此前三节点收发和2/3副本约2000QPS测试的独立通过记录。先根据固定源码的鉴权及 token 准备路径修正测试账号，再以小规模真实连接验证；不得关闭生产鉴权以取得容量数据。
- [1万 job 140611144377](https://github.com/trade3403/-globalchat-im-research/actions/runs/46030019543/job/140611144377)
- [3万 job 140611144820](https://github.com/trade3403/-globalchat-im-research/actions/runs/46030019543/job/140611144820)
- [5万 job 140611144846](https://github.com/trade3403/-globalchat-im-research/actions/runs/46030019543/job/140611144846)

### 8.3 OpenIM：注册阻断

20、10000、30000、50000任务均在首批注册阶段输出 `registration failed 0`，结果记录 connected_total=0、register_batch=1。未进入WebSocket握手，当前无法测定容量。客户端错误输出未保存具体HTTP/业务错误码；用户ID格式、批量参数、服务端错误需要进一步验证，不将猜测作为根因。先修复小规模20客户端闭环，再扩大负载。历史50条REST私聊与接收方历史查询结果依然独立有效。
- [20客户端 job 140612158659](https://github.com/trade3403/-globalchat-im-research/actions/runs/46030318969/job/140612158659)
- [1万 job 140612158809](https://github.com/trade3403/-globalchat-im-research/actions/runs/46030318969/job/140612158809)
- [3万 job 140612158252](https://github.com/trade3403/-globalchat-im-research/actions/runs/46030318969/job/140612158252)
- [5万 job 140612158654](https://github.com/trade3403/-globalchat-im-research/actions/runs/46030318969/job/140612158654)

### 8.4 后续验收与技术借鉴

| 优先级 | 任务 | 验收 |
|---|---|---|
| P0 | WuKongIM 合法鉴权与真实连接 | 小规模成功握手、心跳、目标收件验证 |
| P0 | OpenIM 注册错误定位与20客户端闭环 | 记录脱敏错误码、合法注册、真实在线与唯一消息推送 |
| P0 | JuggleIM 启动退出及约10900连接异常 | 完整启动日志、退出证据、进程资源采样；分别验证启动和负载 |
| P1 | 1万→3万→5万梯度 | 小档成功后再扩档；独立记录同时在线峰值及稳定时间、连接失败原因 |
| P1 | 统一40/1000/2000发布QPS场景 | 同拓扑/硬件/消息大小/持久化语义，记录ACK与端到端分位数和唯一投递集合 |
| P1 | 可靠性及资源 | 丢失、重复、乱序、断线补拉、节点故障、CPU/RSS/FD/网络/磁盘 |
| P2 | 10万连接及GlobalChat Rust对照 | 独立负载机与服务端；重复试验；不得由低负载或配置上限外推 |

GlobalChat 可参考 WuKongIM 的提交/投递边界及副本成本，JuggleIM 的业务层与租户边界，OpenIM 的序号游标与多端同步。三个系统尚无公平的容量排名；GlobalChat 没有加入本轮实测，不能宣称性能领先。此次为报告更新，不修改测试代码、不启动新的测试，也未建立自动写回机制。


## 9. Phase 6D：协议预检与启动稳定性修复（已提交，待验收）

本阶段不宣称新的容量通过结果；先解决第8节的阻断，再恢复1万、3万、5万连接梯度。

| 对象 | 修正措施 | 本轮验收 |
|---|---|---|
| WuKongIM | 根据上游 cloud-small.yaml 与 bench token 测试，将模拟器 token_mode 从 none 改为 bench_api，使用真实 token 准备路径，不取消服务端鉴权 | 1000 个活跃连接、7次连续采样、有效消息流量且无新增收发错误 |
| JuggleIM | MySQL 健康检查改为 TCP 查询最终 jimdb_version=20261103，避免临时初始化实例被判为已就绪；保留关闭 IM 自动重启、pipefail、Docker events、退出状态及资源日志 | 12000 次真实认证连接和60秒私聊测量；仍需独立验证稳定在线数 |
| OpenIM | 合成用户 ID 改为纯字母数字；API 响应先以 RawMessage 解析，非 token 接口不强制 data 为对象；保留 HTTP/业务错误码 | 20个真实鉴权 WebSocket 连接，心跳及唯一收件推送；REST发布，不等同原生客户端发送ACK |

JuggleIM 的数据库启动竞态、OpenIM 的 ID/响应形态是基于源码的修正假设，具体旧故障根因仍需新日志证实。测试工具本身及配置尚未验收。所有变更限独立研究仓库和 Runner 临时源码副本。

- [悟空 IM / JuggleIM 预检工作流](../.github/workflows/im-protocol-preflight.yml)，提交 01e068a7756488efac6c1b713f632e82fb4235c2。
- [OpenIM 预检工作流](../.github/workflows/openim-protocol-preflight.yml)，提交 02f5ec178051a2cba8d0931bb89f2a3da787bd96；[运行 46083477773](https://github.com/trade3403/-globalchat-im-research/actions/runs/46083477773)。
- [OpenIM 压测客户端](../scripts/openim-connection-ladder.go)，提交 8c5f9f12dd3f1922619d23a97327b5626cb4297a。

各项预检通过后再扩大对应产品负载；失败需记录具体阶段及证据，不填造连接容量。后续再执行统一负载、资源效率、故障恢复和 Rust 同条件对比。本次尚未启动新的大规模梯度，未修改 GlobalChat 产品。


## 10. Phase 6D 预检结果（2026-11-28）

两个预检工作流均结束，总体 failure；悟空 IM 子任务 success。来源为本次读取的原始 GitHub job 日志。

| 产品 | 目标 | 实际结果 | 状态 |
|---|---:|---|---|
| WuKongIM | 1000 | 1000 active，连续7次采样，首末相隔约72秒；脚本的有效发送与无新增收发错误条件通过 | 通过 |
| JuggleIM | 12000 | 累计10918次成功连接，未进入完整测量阶段；仍出现 connection reset/refused | 失败 |
| OpenIM | 20 | 20次真实鉴权连接全部成功，测量与宽限期均保持20 active；2400 REST发送成功，2397唯一收件推送 | 连接通过、完整验收失败 |

### WuKongIM
切换 bench_api token 准备后本轮不再出现 ReasonAuthFail，证明此测试环境鉴权修正有效。启动初期 healthz 连接拒绝后恢复；七次连续 active=1000 样本从23:45:26至23:46:38 UTC。证据：[job 140787687634](https://github.com/trade3403/-globalchat-im-research/actions/runs/46083509930/job/140787687634)。仅验证1000连接及当前低负载采样条件，不等于1万或全量接收可靠性验收。

### JuggleIM
12k目标实际累计10918次连接成功，约90.98%；连接失败前出现 reset by peer，随后 connection refused。最终 IM 容器仍 running、restart=0、OOMKilled=false，MySQL healthy。不能据最终容器状态判定内部监听器是否短暂退出，也不能将累计成功连接视作同时在线峰值。与历史10903、10904形成重复异常区间；新数据库就绪检查在本轮允许进入连接阶段，但尚不能断定彻底解决所有历史启动失败。证据：[job 140787687501](https://github.com/trade3403/-globalchat-im-research/actions/runs/46083509930/job/140787687501)。后续应采集进程和监听器事件、连接失败分类与持续在线数。

### OpenIM
| 指标 | 实测 |
|---|---:|
| 注册 / 鉴权连接 | 20 / 20 |
| 最小采样在线 | 20 |
| 连接失败 / 意外断开 | 0 / 0 |
| 心跳Pong累计 | 140 |
| 连接P95 | 6.545ms |
| REST发送成功 / 失败 | 2400 / 0 |
| 唯一目标接收推送 | 2397 |
| 观察到重复 | 0 |
| 测量窗口内未观察到推送 | 3 (0.125%) |
| REST到推送P99 | 90.114ms |

60秒发送窗口及10秒接收宽限期。严格验收要求所有发送均有对应目标推送，因此失败。未观察到的3条不能直接判为永久丢失；后续需要逐消息ID、目标用户及历史补拉核验，区分实时投递、延迟到达、统计缺陷与持久化丢失。无原生客户端发送ACK测量。证据：[job 140787587896](https://github.com/trade3403/-globalchat-im-research/actions/runs/46083477773/job/140787587896)。账号注册与真实连接的阻断已在本轮解除。

下一步：悟空IM可进入1万连接梯度；JuggleIM继续定位约10900异常；OpenIM先核验缺少的3条实时推送，再扩量。当前报告不启动额外CI，不宣称Phase 6D整体通过。正式产品源码未变更。


## 11. Phase 6E：万人梯度与消息完整性诊断（2026-11-28，已提交，待 CI 验收）

本轮复核默认分支原报告提交 0c00fa9 及第10章两个运行的实际 job 状态和失败日志，未发现新的已完成容量结果。悟空 IM 1000 连接已通过，不再重复该档。JuggleIM 仍在 10918 次累计成功连接后出现 reset/refused，最终容器 running、restart=0、OOMKilled=false；未获得足以确定根因的证据。OpenIM 仍为 2400 REST 成功、2397 唯一实时推送、20 在线、0 重复、P99 90.114ms。

| 任务 | 本轮提交 | 验收目标 |
|---|---|---|
| WuKongIM / JuggleIM | [7376f14](https://github.com/trade3403/-globalchat-im-research/commit/7376f14c0899346439aa469580c9443be6a0ae99) | 悟空 IM 保留 bench_api 鉴权和原有7次连续目标在线及流量断言，目标改为10000，采样期限900秒；JuggleIM 保持12000目标与原负载，新增进程/监听器/FD/cgroup及TCP状态连续采样，失败时输出诊断摘要 |
| OpenIM 客户端 | [f654273](https://github.com/trade3403/-globalchat-im-research/commit/f654273cd105c87f135e0b87078827d205dfd944) | 每条成功发送记录合成 marker、预期发送者、接收者；核对真实推送对象；对10秒宽限后未观察到的消息使用接收者 token 查询历史，记录晚到推送、历史恢复数、查询错误和未查到数，不放宽原实时完整性断言 |
| OpenIM CI | [1b14fa4](https://github.com/trade3403/-globalchat-im-research/commit/1b14fa4a788a8396d2306e1e20344b7b78aa8a73) | 20客户端重新运行并产出 message-ledger.json 与 history-reconciliation.json，复核历史补拉接口和严格验收 |

两个工作流均为 main 分支对应 YAML 的 push 路径触发；上述提交已成功写入 Git。当前连接器的 commit workflow runs 接口仅列 PR 触发运行，combined status 返回空数组，未能确认本次 push 的新 Run ID、排队或执行状态。不可写成运行通过或已完成，需从 [Actions](https://github.com/trade3403/-globalchat-im-research/actions) 获取新运行与日志后更新结果。旧 OpenIM 客户端只保存汇总，且临时数据库在 CI 清理，无法从旧日志恢复原3条消息的具体 ID；本轮是可审计复现而非追溯旧数据。新客户端和监控脚本尚未经 CI 编译/执行验证；若历史查询失败，不得将其计为永久丢失。历史补拉使用当前已知的私聊接口与合成账号，非生产账号。JuggleIM 1秒采样可能漏掉短暂事件，且本轮未加入逐连接认证在线证明。

下一步：仅当悟空 IM 10000档通过再扩到30000、50000；JuggleIM 先读监听器与进程资源证据再决定修正；OpenIM 先完成实时推送与历史核对再扩容。保持原始证据、固定上游 SHA、研究仓库隔离；本轮未修改 GlobalChat 正式产品仓库，未获得任何新的容量通过数据。


## 12. Phase 6E 实测结果（2026-11-28）

GitHub Actions 两组运行均已完成，总体 failure，悟空 IM 子任务 success。以下数据来自实际 job 日志，不是配置目标。

| 产品 | 目标 | 实际结果 | 判定 |
|---|---:|---|---|
| WuKongIM | 10000 | 连续7次 active=10000，03:05:37–03:06:49 UTC，约72秒；原有真实鉴权及消息流量断言通过 | 通过 |
| JuggleIM | 12000 | 累计10916次成功连接，随后 reset/refused，未完成消息测量 | 失败 |
| OpenIM | 20 | 20在线，2400 REST发送成功，2397唯一实时推送，3条历史找回，0重复，0意外断开 | 实时验收失败，缺失项历史核验通过 |

- WuKongIM [job 140840847207](https://github.com/trade3403/-globalchat-im-research/actions/runs/46100257265/job/140840847207)：启动后等待约5分50秒才首次采到10000 active；未独立逐消息核对全量投递，也不是单节点容量。下一档为30000。
- JuggleIM [job 140840846959](https://github.com/trade3403/-globalchat-im-research/actions/runs/46100257265/job/140840846959)：首次 reset 02:57:44 UTC，随后宿主机19003端口 refused；02:57:52采样容器内 imserver PID 15777 仍监听9003等端口，进程 FD=29、nofile=262144、ESTAB=13（全容器网络命名空间），memory.events oom_kill=0、pids.events max=0。采样在大量连接失败后，不能证明失败前峰值或排除瞬时问题。重复出现约10900累计成功连接；宿主机端口转发/代理、连接寿命及客户端行为值得优先隔离，尚未确认根因，不是产品架构上限。
- OpenIM [job 140840924060](https://github.com/trade3403/-globalchat-im-research/actions/runs/46100281734/job/140840924060)：连接P95=10.848ms，REST到推送P99=111.233ms，实时缺失率0.125%，缺失 marker 恰为最初三条 capacity-msg-0/1/2，均通过接收者 token 历史查询找到，查询错误0，晚到推送0。仅对本轮缺失项证明历史可读；旧轮三条未追溯，未全量审计已推送消息的持久化。首批推送缺失提示连接注册/路由就绪时序可能有影响，仍需对照实验，不直接认定根因。保留严格实时验收。

原始运行：[悟空/Juggle](https://github.com/trade3403/-globalchat-im-research/actions/runs/46100257265)，[OpenIM](https://github.com/trade3403/-globalchat-im-research/actions/runs/46100281734)。正式产品仓库未修改。


### 12.1 Phase 6F 后续实验（待验收）

- WuKongIM 30000 连接梯度，提交 `1e9ef0f`，保留原鉴权和验收口径，准备/采样期限扩至1800秒，工作流60分钟；不重复已通过10000档。
- OpenIM 20客户端连接完成后等待15秒再发送，保留首条起全部消息及严格实时断言；提交 `b67c212`、`1a2253e`。这是针对首三条缺失的时序对照，并非已修复。
- JuggleIM 仍待隔离宿主机转发和客户端/服务端连接寿命，不重复相同失败条件。
- [Phase 6F OpenIM readiness-delay experiment run 46111585956](https://github.com/trade3403/-globalchat-im-research/actions/runs/46111585956)：查询时 queued。
- [Phase 6F WuKongIM 30k connection ladder run 46111575138](https://github.com/trade3403/-globalchat-im-research/actions/runs/46111575138)：查询时 in_progress。

新增试验结果尚未验收，正式产品仓库未修改。


## 13. Phase 6F 实测结果（2026-11-28）

两个工作流均 completed/success，实际 job 日志复核如下。

| 产品 | 目标 | 结果 |
|---|---:|---|
| WuKongIM | 30000 | 连续7次 active=30000，05:44:07–05:45:19 UTC（约72秒），真实鉴权及有效流量、无新增收发错误的原断言通过 |
| OpenIM | 20 | 20连接全程采样在线，2396 REST成功 / 0失败，2396唯一实时推送，0重复、0意外断开、0未观察到的推送；连接P95 13.577ms，REST到推送P99 101.058ms，Pong 160 |
| JuggleIM | — | 本轮未重复同条件失败测试，约10900累计成功连接异常仍未定位 |

- [WuKongIM job 140876062355](https://github.com/trade3403/-globalchat-im-research/actions/runs/46111575138/job/140876062355)，artifact 15250636763。三节点共享 Runner、抽样消息验证、短时在线观察，不等于单节点3万、全量投递或长时间稳定性。准备阶段长，首次观察到目标在线前状态为 waiting。
- [OpenIM job 140876097002](https://github.com/trade3403/-globalchat-im-research/actions/runs/46111585956/job/140876097002)，artifact 15249124403。15秒连接就绪等待后，本轮所有成功发送均有匹配目标实时推送，严格断言通过；未丢弃初始消息。2396 是实际成功发送数，非固定2400。历史查询只对实时缺项执行，本轮无缺项，不代表全量历史持久化已经核验。单次实验支持就绪时序假设，不能证明上游根因或永久修复，也不能外推万人连接。仍使用 REST 发布而非原生客户端发送ACK。
- 下一步悟空 IM 50000；OpenIM 按1000→10000扩展，保留就绪等待、心跳、真实目标推送和历史核验。JuggleIM 优先隔离宿主机端口代理与直接容器网络路径。正式产品仓库未修改。
