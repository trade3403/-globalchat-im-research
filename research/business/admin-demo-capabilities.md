# IM Admin Demo — 功能与逻辑数据模型观察

- 研究记录：BIZ-13
- 观察日期：2026-11-26
- 范围：用户提供的 66demo 管理后台演示站，仅记录管理端可见功能、页面字段和操作入口
- 证据类型：已登录演示后台后的 UI 观察
- 状态：初步页面盘点，逻辑数据模型待 API/数据库验证

## 1. 结论摘要

该演示管理后台覆盖 IM 社交产品常见的运营、安全、商业化、客户端配置和运维功能。页面能力包括用户/会话管理、群组和频道管理、内容审核、钱包及交易、AI 服务接入、客户端 UI 装修和多平台发布、账号安全审计、运行监控及数据治理。

文中对页面名称、可见表单、筛选项、列表列及按钮的描述来自界面观察。数据实体和字段是为后续产品/技术比较整理的逻辑模型草案，**不代表演示系统的真实数据库结构**。研究未访问源码、数据库、API 文档或供应商合同，亦未验证产品归属、授权范围、价格、生产可用性或商业 SLA。

## 2. 观察方法和限制

- 从用户提供的后台入口登录演示账号，按左侧菜单逐页浏览主要路由，并核对页面标题、表单、筛选控件、表格列、标签页和操作按钮。
- 检查了工作台、用户、私聊、群组、频道、消息查询、通话、客服、聊天规则、内容运营、会员资金、AI、客户端发布、系统设置、安全审计和运维等页面。
- 演示管理员权限使部分写入按钮或字段只读/禁用。没有执行新增、编辑、审核、发送、资金调整、封禁、删除、发布或正式清理。
- 没有把演示环境中出现的真实用户姓名、手机号、IP、UUID、登录记录或私密配置值写入本文。
- UI 树中有些页面延迟加载。个别子页的详情弹窗、全部页签和字段校验规则未逐项验证。页面盘点可作为产品表面证据，不应当作完整验收报告。

## 3. 页面功能盘点

### 3.1 工作台

**业务概览**  
显示总用户、新增用户、在线用户、群/频道、待审核及系统状态；提供 7/30 天新增用户趋势、最新注册/最近活跃记录、在线会话、依赖服务状态、系统告警与常用入口。

**用户管理**  
支持账号精确/模糊搜索和状态筛选。列表表头包含用户信息、性别、注册来源、登录凭证、状态、设备/IP、最后活跃、客服/邀请码归属、注册时间和操作。当前演示角色主要可查看。

### 3.2 聊天与社群

| 页面 | 已观察的页面能力和主要字段 |
|---|---|
| 私聊管理 | 按关键词和状态筛选私聊；列表包含会话双方、会话类型、状态、成员数、创建时间和详情入口。 |
| 群组管理 | 汇总群总数、今日新增、筛选结果及封禁数。列表包含群信息、群主、成员数/成员上限、公开性、状态、创建时间、操作。可见详情、聊天记录、封禁、解散、删除等入口，部分按钮受权限禁用。 |
| 频道管理 | 搜索、状态筛选、频道总数及官方频道标记入口；列表含频道信息、频道主、状态和操作。 |
| 消息查询 | 多条件筛选消息记录；表格可见发送时间等消息审计信息。完整查询项需 API/UI 进一步核对。 |
| 通话记录 | 通话记录及拨号诊断两个标签，含筛选控件、诊断和删除入口。 |
| 官方客服 | 管理客服账号、专属邀请码、欢迎语、备注、旗下用户和启用状态。页面说明邀请码注册可自动加好友并发送欢迎语。 |
| 聊天与关系规则 | 客户端精确/模糊搜索、加好友方式、非好友拉群限制、群成员上限、频道订阅者上限、普通用户撤回时限、消息删除总开关、用户撤回和管理员撤回开关。 |

### 3.3 内容运营

| 页面 | 已观察的页面能力和主要字段 |
|---|---|
| 动态与广场 | 按搜索词、场景、可见范围过滤；有待审核、有原因、待审封面、广场广告等筛选/管理入口。 |
| 话题管理 | 搜索、新建、启停话题；列表含话题名称、描述、动态数、官方属性、排序和状态。 |
| 发现页内容 | 管理轮播图和入口，提供刷新及新增操作。 |
| 表情包管理 | 管理客户端表情商店目录；列表含 Pack ID、名称、描述、预览、贴纸数、排序、内置包、状态和更新时间，可筛选上下架状态。 |
| 公告管理 | 表单包含公告标题、公告类型、公告内容，并可发送公告。 |
| 举报处理 | 多条件筛选举报记录并进入处置流程。举报原因、对象种类、处置结论等具体列需再核。 |
| 敏感词管理 | 按关键词和分类查询，支持单条添加及批量添加违禁词。 |

### 3.4 会员与资金

| 页面 | 已观察的页面能力和主要字段 |
|---|---|
| 会员与权益 | 普通用户额度、套餐设置、会员用户、订单记录等标签；套餐支持新增/编辑入口，部分操作禁用。 |
| 用户钱包 | 搜索用户钱包，可进入钱包管理或余额调整。 |
| 充值管理 | 充值方式、充值审核两个标签；支持新增充值方式。 |
| 提现管理 | 按状态检索提现记录/申请。 |
| 红包记录 | 按条件查询红包并打开详情。 |
| 转账记录 | 搜索、重置、分页查看转账记录。 |
| 资金规则 | 基础设置、提现方式、实名要求；基础设置包含货币符号/名称、红包过期时间、转账过期时间、钱包公告和提现公告。 |

### 3.5 机器人与 AI

| 页面 | 已观察的页面能力和主要字段 |
|---|---|
| AI bots | 机器人列表和配置；群内自动启动开关、群摘要/治理/高级自动化配置入口、待人工复核列表。 |
| 知识库与文档 | 知识库/RAG 页面，支持新建知识库。文档导入/切分字段需进一步核实。 |
| 模型、语音与翻译 | 模型供应商、语音转文字与翻译分区。可见供应商名、Base URL、模型名、API Key、超时、最大输出 Token、输入/输出价格、月预算、默认供应商和内容审核开关，带连接测试/保存入口。 |
| Mini App 与应用市场 | Mini App 配置和机器人应用市场管理；可新建 Mini App、新增市场应用、刷新列表。 |
| 用量与成本 | AI 运行概览入口；各费用维度和计量口径尚未确认。 |

### 3.6 客户端与发布

| 页面 | 已观察的页面能力和主要字段 |
|---|---|
| 装修工作台 | 品牌与颜色、底部导航、联系人/发现/设置/个人中心/聊天工具栏、公告活动、服务链接、协议文档、启动页与栏目；草稿/已发布预览、深色模式、套用方案、查看改动、保存草稿和发布客户端。 |
| 素材管理 | 全部/默认素材/新 UI/自己上传分类筛选；搜索、选择图片、草稿引用和刷新。上传、保存受权限控制。 |
| 装修发布 | 发布管理、操作记录；Android/iOS/Windows/macOS/H5 发布范围及各平台最低版本；提供导出草稿、发布草稿入口。 |
| App 版本与下载 | 系统版本、各平台最新版本和最低支持版本、iOS 商店/Android 下载地址、强制更新策略。 |
| 热更新补丁 | 补丁搜索、过滤、刷新、新建入口和记录列表。补丁 payload 的精确字段未核验。 |

### 3.7 系统配置

| 页面 | 已观察的页面能力和主要字段 |
|---|---|
| 配置概览 | 配置概览/筛选入口，卡片和字段未全部展开记录。 |
| 基本信息与客服 | 系统名称、H5 分享/邀请地址、在线客服链接、QQ 客服号。 |
| 注册登录与身份规则 | Apple 登录、注册与身份规则；开放注册、一键注册、设备/IP 日限额、邀请码、性别必填、手机号绑定。 |
| 客户端功能开关 | 按功能组控制客户端开关，页面可见聊天扩展菜单、钱包/VIP、消息待办等配置分组。 |
| 短信服务 | 短信开关、渠道、正文模板、用户名和密码配置。 |
| 推送服务 | 推送服务配置、Android 保活策略入口。 |
| 第三方服务 | 音视频、实名认证、支付接入；观察到 Agora/LiveKit 选择及服务参数表单。 |
| 文件与存储 | 本地、阿里云 OSS、七牛云、Amazon S3 选项；访问域名、图片直传、平台适用范围等。 |
| 协议与文案 | 用户协议、隐私政策文本和预览入口。 |

### 3.8 安全与审计

| 页面 | 已观察的页面能力和主要字段 |
|---|---|
| 管理员账号 | 账号、昵称、角色、状态、Google 验证状态、最后登录 IP/时间、创建时间；可刷新，并对未绑定验证的账号提供绑定验证入口。 |
| 角色与权限 | 角色权限矩阵，提供全选、清空、恢复默认和保存。 |
| 登录与在线活动 | 全球登录态势和登录日志。 |
| 访问与加密策略 | IP/用户消息频率限制；明文、兼容加密、严格加密模式。 |
| 企业加密对话 | 留存策略、密钥环境摘要、审计员、加密会话查询及审计授权。 |
| 操作审计 | 管理员操作记录入口，完整表头/过滤字段需后续逐项确认。 |

### 3.9 运维中心

| 页面 | 已观察的页面能力和主要字段 |
|---|---|
| 运行监控 | 实时监控、实例健康、健康趋势、告警中心；趋势范围 15 分钟/1 小时/6 小时/24 小时/7 天，告警有当前/已恢复标签。 |
| 服务器与容量 | 容量概览、服务器节点、架构与诊断；集群 CPU、内存、本地磁盘和节点状态。 |
| 客户端连接线路 | 动态入口开关、配置版本、客户端缓存、访问线路、节点分配模式、节点 ID、API 地址、WebSocket 地址、备注、优先级。 |
| 推送投递报告 | 通道状态、监控概览、测试推送；可选用户 ID/设备 ID、推送场景。 |
| 数据治理 | 治理总览、待清理数据、文件管理与预览、容量/健康评分/增长趋势、保留周期、候选清单、保护上限、备份门禁和执行审计。页面显示正式清理需要备份门禁，正式清理可能永久删除数据。研究中没有运行清理。 |

## 4. 逻辑数据库模型草案

下列 snake_case 字段仅为统一表达的候选逻辑结构。建议大部分实体统一加 id、created_at、updated_at。凡页面计数可能由关系表计算的字段，建议不要在主表重复保存，除非缓存/性能需求明确。

### 4.1 用户、认证、设备

| 实体 | 候选字段 |
|---|---|
| users | id, uuid, account, short_id, nickname, avatar_url, gender, registration_source, credential_status, phone, status, customer_service_user_id, invite_code, last_active_at, registered_at |
| user_auth_identities | id, user_id, provider, provider_subject, phone_verified_at, status |
| user_devices | id, user_id, platform, device_model, device_id_hash, last_ip, push_status, bound_status, last_seen_at |

### 4.2 聊天、社群、通话

| 实体 | 候选字段 |
|---|---|
| conversations | id, conversation_type, status, member_count, created_at, last_message_at |
| conversation_members | conversation_id, user_id, role, joined_at, left_at, member_status |
| groups | id, uuid, name, owner_user_id, description, member_limit, visibility, status, official_flag, created_at |
| channels | id, uuid, name, owner_user_id, subscriber_limit, status, official_flag, created_at |
| messages | id, conversation_id, sender_user_id, message_type, content_ref, sent_at, recalled_at, recall_actor_id, visibility |
| calls | id, conversation_id, caller_user_id, callee_user_id, call_type, status, started_at, ended_at, duration_seconds, diagnostic_ref |
| chat_rules | search_mode, friend_add_mode, non_friend_invite_group, group_member_limit, channel_subscriber_limit, user_recall_window_seconds, message_delete_enabled, user_recall_enabled, admin_recall_enabled |

### 4.3 内容、审核、公告、素材

| 实体 | 候选字段 |
|---|---|
| moments | id, author_user_id, content, media_refs, scene, visibility, cover_status, review_status, review_reason, created_at, reviewed_by, reviewed_at |
| topics | id, name, description, official_flag, sort_order, enabled, created_at |
| reports | id, reporter_user_id, target_type, target_id, reason_code, description, status, assigned_admin_id, resolution, created_at, resolved_at |
| banned_words | id, word, category, match_mode, enabled, created_by, created_at |
| announcements | id, title, type, body, status, published_by, published_at, created_at |
| emoji_packs | id, pack_id, name, description, preview_asset_id, sort_order, builtin_flag, status, updated_at |
| emoji_stickers | id, pack_id, asset_id, name, sort_order, status |
| discover_entries | id, entry_type, title, icon_asset_id, target_url, sort_order, enabled |
| assets | id, asset_type, storage_provider, object_key, url, mime_type, size_bytes, width, height, uploaded_by, created_at |

### 4.4 客服、会员、钱包与交易

| 实体 | 候选字段 |
|---|---|
| official_service_accounts | id, user_id, invite_code, welcome_message, note, enabled, created_at |
| vip_plans | id, name, description, price, currency, duration_days, quota_config, enabled, sort_order |
| vip_memberships | id, user_id, plan_id, status, started_at, expires_at, source_order_id |
| wallets | id, user_id, currency, balance, frozen_balance, updated_at |
| wallet_ledger | id, wallet_id, transaction_type, amount, balance_after, reference_type, reference_id, operator_id, note, created_at |
| recharge_methods | id, name, provider, config_ref, enabled, sort_order |
| recharge_orders | id, user_id, method_id, amount, currency, status, provider_trade_no, created_at, reviewed_by, reviewed_at |
| withdrawal_orders | id, user_id, amount, currency, method, real_name_status, status, reviewed_by, reviewed_at, created_at |
| red_packets | id, sender_user_id, group_id, total_amount, currency, packet_count, claimed_amount, status, expires_at, created_at |
| transfers | id, sender_user_id, recipient_user_id, amount, currency, status, expires_at, created_at |
| wallet_settings | currency_symbol, currency_name, red_packet_expiry_hours, transfer_expiry_hours, wallet_notice, withdraw_notice, withdrawal_methods, real_name_required |

### 4.5 AI 与客户端配置

| 实体 | 候选字段 |
|---|---|
| ai_providers | id, display_name, base_url, model_name, api_key_secret_ref, timeout_seconds, max_output_tokens, input_price_per_million, output_price_per_million, monthly_budget, is_default, moderation_enabled, enabled |
| ai_bots | id, name, provider_id, model_id, system_prompt_ref, enabled, auto_start_on_group_join, created_at |
| knowledge_bases | id, name, description, embedding_provider_id, status, created_at |
| knowledge_documents | id, knowledge_base_id, file_asset_id, parse_status, chunk_count, indexed_at, status |
| mini_apps | id, app_id, name, description, entry_url, icon_asset_id, status, sort_order |
| market_apps | id, app_id, name, category, version, listing_status, sort_order |
| client_brand_config | id, brand_name, primary_color, theme_config_json, logo_asset_id, version, status |
| client_nav_items | id, config_version_id, name, route, icon_asset_id, active_icon_asset_id, sort_order, enabled |
| client_content_config | id, config_version_id, section_type, payload_json, sort_order, enabled |
| client_config_releases | id, version, draft_payload, release_scope, minimum_versions_json, status, created_by, published_by, published_at |
| app_versions | platform, latest_version, minimum_supported_version, download_url, force_update_enabled, updated_at |
| hot_update_patches | id, platform, version, package_url, checksum, mandatory, status, published_at, created_by |

### 4.6 系统、安全、审计、运维和数据治理

| 实体 | 候选字段 |
|---|---|
| system_settings | setting_key, setting_value_json, scope, updated_by, updated_at |
| admin_accounts | id, username, display_name, role_id, status, google_auth_enabled, last_login_ip, last_login_at, created_at |
| admin_roles | id, name, description, is_system_role, created_at |
| admin_role_permissions | role_id, permission_key, allowed |
| admin_login_logs | id, admin_id, ip, user_agent, result, failure_reason, created_at |
| admin_audit_events | id, admin_id, action, resource_type, resource_id, before_json, after_json, ip, created_at |
| security_policies | ip_message_rate_limit, user_message_rate_limit, encryption_mode, retention_policy_json |
| secure_chat_audit_grants | id, auditor_admin_id, conversation_id, reason, granted_by, expires_at, created_at |
| service_health_samples | id, service_name, status, latency_ms, error_rate, sampled_at |
| servers | id, node_id, name, status, cpu_usage, memory_usage, disk_usage, last_seen_at |
| client_endpoints | id, node_id, api_url, websocket_url, note, priority, enabled, config_version |
| push_delivery_logs | id, provider, user_id, device_id, scenario, status, error_code, sent_at |
| data_retention_policies | id, dataset, retention_days, enabled, max_delete_per_run, updated_by |
| cleanup_manifests | id, generated_at, expires_at, candidate_count, protected_count, backup_gate_status, status |
| cleanup_audit_logs | id, manifest_id, execution_mode, deleted_count, skipped_count, error_count, operator_id, created_at |

## 5. 字段映射和实现注意事项

| UI 项 | 候选逻辑字段 | 备注 |
|---|---|---|
| 用户信息/账号 | users.account, users.nickname, users.uuid | 列表可能是多个字段拼接。 |
| 设备/IP | user_devices.device_model, user_devices.last_ip | 属于敏感信息，需最小权限和留存策略。 |
| 群成员 159/500 | conversation_members 聚合 + groups.member_limit | 成员数可计算，成员上限来自群设置。 |
| 官方/启用状态 | official_flag, enabled/status | UI 状态文案应映射稳定的后端枚举值。 |
| 模型 API Key | ai_providers.api_key_secret_ref | 建议由密钥管理服务保存，不将明文回传。 |
| 草稿/已发布 | client_config_releases.status/version | 发布配置应保存不可变版本快照。 |
| 线路顺序 | client_endpoints.priority | 拖动排序需映射稳定优先级。 |
| 保留周期 | data_retention_policies.retention_days | 清理执行应受限时清单、备份门禁和审计流程控制。 |

## 6. 通用数据设计建议

1. 金额使用最小币种单位整数或 DECIMAL，避免浮点误差。钱包应有不可变流水，不只直接改余额。
2. 时间统一存 UTC，并在页面层转换展示时区。
3. 枚举字段存稳定 code，中文文案由前端/字典映射。
4. 对账号/状态、会话/消息时间、订单状态/创建时间、举报状态/时间、审计操作者/时间建立合适索引。
5. 用户手机号、IP、设备标识、API 密钥和供应商证书按敏感等级加密或限制访问。
6. 版本发布采用草稿与不可变发布快照，并记录操作者、平台范围、时间及回滚基线。
7. 数据清理需要保留候选清单、保护数量、备份门禁结果、执行数量和操作者审计。
8. 若系统支持多租户，为核心业务实体增加 tenant_id，并设计跨租户隔离约束。
9. 删除流程应明确软删除、归档和永久清理的边界，资金和安全审计记录通常需要独立留存策略。

## 7. 后续核验清单

要将本草案变成准确的物理数据库字典，需要读取数据库迁移/DDL、API OpenAPI schema 和实际请求响应，并核实：

- 真实表/列名称、数据类型、可空性、默认值、唯一约束、索引及外键。
- UI 中各状态的后端枚举 code。
- 钱包金额的币种精度、订单状态机、退款/冲正和账本幂等约束。
- 用户与身份凭证、设备、群成员、频道订阅的关系基数。
- 消息内容的加密、附件引用、撤回及保留策略。
- 各配置项是否拆分成多表、JSON 配置或密钥系统引用。
- 每个页面未展开的详情抽屉、二级标签、批量操作和权限矩阵。

该演示站页面观察不足以证明其生产级能力、规模、SLA、功能授权或实际技术实现。用于竞品研究时应保留此证据边界。
