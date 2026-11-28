# GlobalChat IM Research

独立研究仓库，不修改 trade3403/wechat。

## 研究顺序
1. WuKongIM：网关、协议、ACK、持久化、分片、Raft。
2. JuggleIM：通信与业务模块、运营后台。
3. OpenIM：SDK、同步、业务体系。
4. GlobalChat：Rust 差距分析，待产品源码授权接入。

CI 首轮只采集公开源码的版本、文件清单、许可证与关键字证据，不执行第三方源码。关键字命中不是功能证明；需人工追踪调用链与测试。性能压测另行设计。禁止将私有源码或凭据输出到报告。

## 源码构建与运行验证

[Phase 7 源码对比、CI 矩阵与任务计划](reports/phase7-im-source-ci-comparison.md) · [运行记录](https://github.com/trade3403/-globalchat-im-research/actions/workflows/source-runtime-comparison.yml)

## 并发与性能实测（Phase 6C）

[三家 IM 容量、延迟、失败根因及后续记录规范](reports/phase6c-im-capacity-analysis.md)。包含已实测指标、未验证范围、1万/3万/5万阶梯测试、GlobalChat Rust 建议验收目标。每轮完成后须更新 Git 报告，保留失败证据及原始 CI 链接。

[悟空 IM / JuggleIM 阶梯 CI](https://github.com/trade3403/-globalchat-im-research/actions/workflows/connection-ladder.yml) · [OpenIM 阶梯 CI](https://github.com/trade3403/-globalchat-im-research/actions/workflows/openim-connection-ladder.yml)


## 单节点与集群公平对照（Phase 8）

[单节点、三节点、五节点部署和统一验收矩阵](reports/phase8-im-topology-comparison.md) · [悟空 IM 单节点 CI](.github/workflows/wukong-single-preflight.yml) · [OpenIM 三实例 CI](.github/workflows/openim-cluster-preflight.yml) · [JuggleIM 桥接诊断 CI](.github/workflows/juggle-bridge-diagnosis.yml)。区分逻辑实例与独立主机、开放社区版与商业授权；按相同总资源和相同每节点资源分别比较。
