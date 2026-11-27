# Phase 7：三套 IM 的源码对比、CI 矩阵与任务计划

本次分析范围为公开上游源码。CI 在 GitHub 托管的 Ubuntu Runner 上使用临时环境和合成账号。测试场景覆盖各系统可独立核验的消息路径。各场景的协议、拓扑和负载不同，结果不可直接用于吞吐排名。

## 源码对比

| 项目 | 固定版本 | 源码链路 | 功能优势与边界 |
|---|---|---|---|
| 悟空 IM | `21d905a328bb7141167bc6dbe20fdb65588e0420` | `internal/contracts/channelappend/FLOW.md` 定义消息追加边界，`internal/runtime/channelappend` 实现频道追加运行时，`internal/contracts/onlinedelivery/FLOW.md` 定义在线投递边界。跨节点场景入口在 `test/e2e/message/cross_node_delivery`。 | 重点是按频道组织消息与跨节点投递。内置持久化和集群减少了独立依赖，账号、好友和业务权限由上层应用负责。 |
| JuggleIM | `0a47433c1aed6f08b48a3bfcf6b80dfe0a55a6c4` | `services/connectmanager/server/imwebsocketmsghandler.go` 处理长连接消息，`docs/architecture_zh.md` 描述 Actor/RPC、`app_key` 和消息存储路径，`commons/dbcommons/dbmigration.go` 处理启动迁移，`simulator/benchmark/runner.go` 记录 ACK 和接收回调。 | 消息核心、JuggleChat 社交服务和运营控制台分开。社区版使用单节点部署，集群能力按专业版范围理解。 |
| OpenIM | `a101d85ec16cc21bec60005169725adc1ae5649a` | `cmd/main.go` 装配单进程服务；`internal/api/router.go` 和 `internal/api/msg.go` 接入发送和拉取接口；`internal/rpc/msg/sync_msg.go` 校验访问并按序号查历史。 | 用户、好友、群组、会话和消息能力可由应用和 SDK 使用。单进程启动与多 RPC 服务部署是两种运行方式，本次只运行单进程。 |
| OpenIM SDK | `5c9670c7fe9a2d5993f5ef6964852d391470f40f` | `internal/interaction/msg_sync.go` 和 `internal/conversation_msg/incremental_sync.go` 提供消息及会话同步入口。 | 多端共享客户端同步核心。这次只读分析 SDK，运行矩阵只验证服务端 API。 |

源码：[悟空 IM](https://github.com/WuKongIM/WuKongIM/tree/21d905a328bb7141167bc6dbe20fdb65588e0420) · [JuggleIM](https://github.com/juggleim/im-server/tree/0a47433c1aed6f08b48a3bfcf6b80dfe0a55a6c4) · [OpenIM Server](https://github.com/openimsdk/open-im-server/tree/a101d85ec16cc21bec60005169725adc1ae5649a) · [OpenIM SDK](https://github.com/openimsdk/openim-sdk-core/tree/5c9670c7fe9a2d5993f5ef6964852d391470f40f)

### 消息链路中的差异

- **悟空 IM**：入口将发送命令交给频道追加路由；提交后才生成供收件人规划使用的消息封装。在线投递又把接受、可重试和丢弃分别记录，因此发送端确认和接收端投递可分开观察。
- **JuggleIM**：WebSocket 帧进入连接管理器，内部调用保留 `app_key`。私聊路径在分发前应用黑名单、好友规则和拦截器，再记录发送箱和历史。上游压测分别记录服务端 ACK 和收件端回调，回调统计无法直接推导未收到的消息数。
- **OpenIM**：服务端 REST 入口把消息交给消息 RPC，接收者可按会话序号查询历史。客户端 SDK 的 `MsgSyncer` 维护本地序号，在建立连接或处理推送时比较本地与目标序号，并拉取缺失区间。本次 REST 测试覆盖服务端序号拉取，SDK 重连和本地库仍待单独实测。

## CI 分析表

| 项目 | 历史实测 | Phase 7 构建与测试 | 通过条件 | 产物 |
|---|---|---|---|---|
| 悟空 IM | 两次跨节点收发通过 | run 46011555775 的官方三节点 E2E 实际通过 | 跨节点收发测试通过 | Actions 摘要含 Go JSON pass 事件。GitHub artifact 配额耗尽导致首轮上传失败，最终通过轮产物上传成功 |
| JuggleIM | 冒烟通过；前两轮收发测试未就绪 | run 46011555775 的临时 Compose 栈和私聊、群聊 JSON 断言均通过 | 本地私聊、群聊的连接与 ACK 场景通过，接收统计未逐条核对每个预期收件人 | 私聊/群聊报告、容器状态、日志和临时 SQL 修复补丁 |
| OpenIM | 首轮仅确认 Compose 启动步骤 | run 46011555775 从源码编译，管理员鉴权、注册测试用户、发送私聊及接收方历史拉取均通过 | 接收方能从按会话序号查询的历史接口读回刚发送的唯一消息 | 私聊历史 JSON、服务日志、容器状态和运行配置补丁 |

JuggleIM 的固定版本中，`sql/imserver.sql` 有一个多余的 `+CREATE TABLE`，使 MySQL 在写入最终初始化版本标记前中断。后续自动迁移会碰到已存在的 `is_portion`。Phase 7 只在 Runner 的临时源码副本去掉该字符并保存补丁，未改上游源码。该轮结果应视为“修复测试环境后通过”。

Phase 7 成功运行：[run 46011555775](https://github.com/trade3403/-globalchat-im-research/actions/runs/46011555775)。首轮 OpenIM 测试使用了协议不接受的用户 ID，修正后通过。另有首轮 artifact upload 报错，原因是 GitHub Actions artifact storage quota exceeded，即 Actions 产物存储配额耗尽，未能上传报告；不是 Runner 磁盘满。最终轮上传成功。\n\n历史证据：[冒烟](https://github.com/trade3403/-globalchat-im-research/actions/runs/45772497250) · [跨节点和 JuggleIM 诊断](https://github.com/trade3403/-globalchat-im-research/actions/runs/45782768181)。新版运行：[Phase 7 Actions](https://github.com/trade3403/-globalchat-im-research/actions/workflows/source-runtime-comparison.yml)。

## 任务计划表

| 编号 | 任务 | 验收产物 | 状态 |
|---|---|---|---|
| IM-CI-01 | 固定四个源码版本，追踪消息追加、投递和同步路径 | 本文源码表与上游链接 | 完成 |
| IM-CI-02 | 在 Runner 复测悟空 IM 三节点收发 | 官方 E2E 的通过事件 | 完成，run 46011555775 通过 |
| IM-CI-03 | 复测 JuggleIM 数据库新装和私聊、群聊 | 可审阅的临时修补、ACK 和回调报告 | 完成，临时修补后通过 run 46011555775 |
| IM-CI-04 | 从源码构建 OpenIM，并验证私聊历史 | 接收方读回私聊的 JSON 证据 | 完成，run 46011555775 通过 |
| IM-CI-05 | 建立三家统一的消息正确性场景 | 相同账号数、消息数与正文长度；按消息 ID 和设备核对丢失、重复、顺序、断线补齐 | 待实现 |
| IM-CI-06 | 比较三家性能 | 固定 Runner 规格、拓扑与负载，记录发送 ACK 和接收延迟 P50/P95/P99、吞吐、CPU、内存 | 依赖 IM-CI-05 |
