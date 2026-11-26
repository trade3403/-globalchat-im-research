# Phase 5 — GlobalChat source gap audit (preliminary)

Product baseline: trade3403/wechat @ c360758470edae514d237256c2b8c5cb270f39d1 (phases 1–4 integration). Phase 5 account branch diverges and is not included. Read-only inspection of ten Rust source files and Cargo.toml via GitHub. No product source copied or modified.

| Domain | Observed implementation (source location) | Gap / next verification |
|---|---|---|
| Transport | main.rs:615 live_upgrade, :636 live_socket; Axum WebSocket | Production gateway, reconnect/load tests and TLS deployment not established |
| Reliability | main.rs:279 store_message, :349 sync_messages, :429 acknowledge; tests :1125–1367 | Trace transaction and ACK ordering under crash, retry and concurrent load; compare WuKongIM |
| Storage | main.rs:916 initialize_schema, :990 open_database, :1012 open_database_pool; rusqlite | SQLite baseline; evaluate partitioning, replication and recovery, not assume replacement necessary |
| Authentication | auth.rs:249 register, :287 login, :363 resolve, :400 logout | Review session revocation, rate limits, security and tenant isolation |
| Private and group chat | auth.rs:428 send, :572 history; groups.rs:334 send, :523 sync, :592 ack | Review shared channel abstraction and compatibility of existing API |
| Business | social.rs, group_admin.rs, group_invites.rs, moments.rs, profile.rs | Existing module files; verify separation of transport, policy and persistence before extraction |
| Media | attachments.rs:82 upload, :176 download | Verify authorization, quotas, durable object storage and lifecycle |
| RTC | calls.rs:278 transition, :366 send_signal | Signaling exists; actual media transport not established |
| SDK | existing uni-app product | Audit client code and platform bridges separately; compare OpenIM |
| OEM | not established in inspected source | Tenant identity, policy interfaces and server-enforced isolation require explicit design |

The test server is described by its Cargo manifest as loopback-only experimental, not production-ready. Existing tests and symbols are evidence of implementation, not current test results. Cluster, OEM and SDK capability gaps are unverified, not claims of absence. Previous CI artifacts are static indexes; compare actual call paths before adopting any architecture. No benchmark claims.

Next gate: retrieve prior competitor evidence, confirm report integrity, then design equivalent reliability and performance test scenarios. Cross-repository CI checkout requires separate read-only credentials for the private product repository; do not expose tokens in logs or copy private source into artifacts.
