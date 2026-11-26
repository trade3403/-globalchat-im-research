# GlobalChat IM Research

独立研究仓库，不修改 trade3403/wechat。

## 研究顺序
1. WuKongIM：网关、协议、ACK、持久化、分片、Raft。
2. JuggleIM：通信与业务模块、运营后台。
3. OpenIM：SDK、同步、业务体系。
4. GlobalChat：Rust 差距分析，待产品源码授权接入。

CI 首轮只采集公开源码的版本、文件清单、许可证与关键字证据，不执行第三方源码。关键字命中不是功能证明；需人工追踪调用链与测试。性能压测另行设计。禁止将私有源码或凭据输出到报告。
