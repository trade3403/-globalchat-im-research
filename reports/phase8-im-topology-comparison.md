# Phase 8：单节点与集群公平对照研究

## 范围与部署分类
仅研究仓库；固定上游 SHA，与 Phase 6C–6G 一致。禁止改动正式产品仓库。逻辑实例不等于独立物理主机，多个网关不等于完整高可用集群。

| 产品 | 单节点 | 三节点/实例 | 五节点/实例 | 约束 |
|---|---|---|---|---|
| WuKongIM | 新增预检 | 已有实测 | 计划 | 内置集群；单节点 slot replica=1，三/五节点分别记录副本数及 quorum |
| JuggleIM | 社区版已测 | 专业版待授权 | 专业版待授权 | 不用独立实例冒充集群，不绕过授权 |
| OpenIM | 单进程已测 | 新增多独立进程预检 | 三实例通过后扩展 | 共享 Redis/MongoDB/MinIO；Redis 队列实验环境与默认 Kafka 微服务部署分开 |
| GlobalChat Rust | 待接入 | 规划 | 规划 | 正式产品仓库此次不修改，不宣称已有集群能力 |

上游证据：[JuggleIM README](https://github.com/juggleim/im-server/blob/0a47433c1aed6f08b48a3bfcf6b80dfe0a55a6c4/README.md) 明确专业版支持集群；[OpenIM cmd/main.go](https://github.com/openimsdk/open-im-server/blob/a101d85ec16cc21bec60005169725adc1ae5649a/cmd/main.go) 有 standalone Redis 地址注册及跨实例广播路径，[网关配置](https://github.com/openimsdk/open-im-server/blob/a101d85ec16cc21bec60005169725adc1ae5649a/config/openim-msggateway.yml) 支持端口配置。完整 etcd/RPC/Kafka 多服务集群另列任务，不将 standalone 多实例称为完整微服务 HA。

## 阶梯与验收
1. 功能预检：各节点实际监听、合法 token、每节点在线数、跨实例双向消息及唯一 ID、历史读取。OpenIM 三实例预检使用30个合成账号，10个/实例，按用户ID轮转分配，REST发送入口按发送者所在实例选取，接收者由另一实例真实 WebSocket 接收；不等同原生客户端发送。五实例须三实例通过后再建立。
2. 容量：1000→10000→30000→50000→100000，每档记录在线采样、连接成功/失败、连接速率、真实消息流量；逐档通过才扩量。
3. 公平性能：同一硬件、相同消息大小、发送频率（40/1000/2000 QPS，再扩档）、连接分布、TLS、持久化 ACK 语义及副本数；不把不同业务语义强行混为同一指标。私聊、群聊分开。比较单节点对单节点；3对3、5对5。分别记录固定总资源与固定每节点资源两组实验，每档至少三次。
4. 资源：全系统 CPU、RSS、内存、磁盘、网络、FD、依赖组件、压测端开销；同机 Runner 与独立物理主机分开报告。
5. 可靠性：按期望目标集合核对 ACK、唯一实时推送、历史持久化、重复与顺序；断开一个节点后重连、跨节点补拉、恢复时间。依赖单点仍在时不得宣称整套 HA。
6. 持续负载：1h、24h、168h；长时采用适合的独立 Runner/服务器，不假设免费额度或单次 GitHub job 可覆盖所有时长。
7. GlobalChat Rust：未来按单/集群一键部署、服务发现、跨节点路由、多副本持久化、恢复、观测与扩容验收；源码授权后单独接入。

## 已提交的执行入口
- [WuKongIM 单节点预检](../.github/workflows/wukong-single-preflight.yml)：1000连接，临时上游配置仅保留 node1 与单副本，模拟器不自动启动依赖的其他节点。提交 d3b4526；[run 46147935753](https://github.com/trade3403/-globalchat-im-research/actions/runs/46147935753) 查询时 queued。
- [OpenIM 三实例预检](../.github/workflows/openim-cluster-preflight.yml)：同一 Runner 三个独立 standalone 进程，API/WS 端口各自隔离，共享依赖，真实鉴权与跨实例接收。客户端提交 b52a67e；工作流 b869d69；[run 46147921354](https://github.com/trade3403/-globalchat-im-research/actions/runs/46147921354) 查询时 in_progress。未验证进程启动或跨实例推送成功。
- [JuggleIM 桥接诊断](../.github/workflows/juggle-bridge-diagnosis.yml)：仅在隔离 Runner 内使用 Docker inspect 返回的桥接地址，显式启用上游允许非 loopback 的测试参数；其他业务逻辑与鉴权不变。提交 8da9196，待验收。
- WuKongIM 五节点、OpenIM 五实例及真正分布式 RPC 集群尚未实现工作流；JuggleIM 专业版待授权。禁止将预检提交当作容量通过。

商业研究：集群能力与社区/专业版边界、OEM 授权、私有化支持、客户部署成本，分别记录实测与厂商声称；不做未验证性能排名。
