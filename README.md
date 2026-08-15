<div align="center">
<img src="webs/src/assets/logo.png" width="150px" height="150px" />
</div>

<div align="center">
  <img src="https://img.shields.io/badge/Go-1.24.3-green.svg"/>
  <img src="https://img.shields.io/badge/Vue-3-brightgreen.svg"/>
  <img src="https://img.shields.io/badge/Element%20Plus-2.x-blue.svg"/>
  <img src="https://img.shields.io/badge/license-MIT-green.svg"/>
  <div align="center"> 中文 | <a href="README.en-US.md">English</div>
</div>

# 项目简介

`sublinkE` 是基于优秀开源项目 [sublinkX](https://github.com/gooaclok819/sublinkX) 进行二次开发的自托管**代理节点订阅管理/分发系统，在保留原项目全部能力的基础上，加入了多用户体系、安全加固与多种实用功能。

- 前端基于 [vue3-element-admin](https://github.com/youlaitech/vue3-element-admin)；
- 后端采用 Go + Gin + Gorm（SQLite 单文件数据库）；
- 默认账号：admin 密码：123456，**请安装后务必自行修改**；
- 当前版本：v1.1.8.3
- 本项目使用VibeCodeing

# 功能特性

## 节点与订阅

- 支持多种代理协议：**ss、ssr、trojan、vmess、vless、hy、hy2、tuic、AnyTLS、Socks5**
- 订阅输出格式：**v2ray（base64）、Clash（YAML）、Surge（conf）**，支持按 `User-Agent` 自动识别客户端
- 订阅导入与**定时更新**（cron 表达式），自动解析节点入库
- 节点排序（拖拽）、Clash `dialer-proxy` 属性、IPv6 地址包裹处理
- 模板管理：自定义 Clash / Surge 订阅模板

## 多用户体系

- 用户注册（**滑块验证码 + 拖动轨迹检测**），支持邀请码注册开关与**邀请码使用上限**
- 新注册用户自动分配默认订阅，每个用户持有独立的**订阅 Token**
- 订阅链接统一为用户 Token 体系（已移除旧版 `md5(订阅名)` 链接），管理员可在「用户管理」分配订阅、限制可用地区、禁用/删除用户、重置订阅 Token
- 用户中心：查看/复制自己的订阅链接（自动识别 / Clash / Surge / V2Ray）、修改密码（需验证旧密码）、查看拉取记录

## 访问控制与审计

- **地区限制**：可为用户配置允许地区（如 `中国,福建`），按 IP 归属地匹配，超范围拒绝拉取
- **User-Agent 检测**（后台可开关）：开启后订阅拉取仅允许代理软件（关键字列表**可在后台配置**），浏览器/脚本直接访问返回 403
- **拉取日志**：记录每次拉取的 IP、地区、客户端、**User-Agent** 与状态（成功 / 地区拦截 / UA 拦截），管理员与用户均可查看
- 禁用用户即时断流：无法登录、无法通过订阅 Token 拉取、已签发 JWT 立即失效

# 安全特性

- **API Key 归属校验**：普通用户只能为自己生成/删除/查看 API Key，防止越权提权
- **接口权限收敛**：用户列表、模板管理接口仅管理员可访问
- **登录/注册限流**：登录 10 次/15 分钟/IP，注册 5 次/小时/IP，防暴力破解与批量注册
- **Token 吊销机制**：改密、登出、被禁用后，已签发的 JWT 立即失效（token 版本号自增）
- **JWT 有效期可配置**（`expire_days`），替代原先硬编码 14 天
- **密码安全**：bcrypt 存储（禁止明文回退）、后端强制最小 6 位
- **密钥保护**：管理接口返回的 `jwt_secret` / `api_encryption_key` 已掩码
- **滑块验证码**：缺口位置仅存服务端内存，校验拖动位置（容差内）+ 拖动轨迹启发式检测，拦截脚本直提
- **订阅链接安全**：用户 Token 为随机串（`sub_时间戳_24位随机`），不可推导；禁用/重置即时失效

# 安装说明


## 一键安装

```bash
wget https://raw.githubusercontent.com/Yhorizen/sublinkE/main/install.sh && sh install.sh
```

> ⚠ **注意**
> 在 **Alpine Linux** 上运行时，由于 Alpine 使用 `musl` 而非 `glibc`，插件模块无法正常工作。
> 推荐优先使用 **Docker 部署**，或选择 **Debian / Ubuntu** 等发行版。

## 直接运行

- Windows：运行 `sublinke-windows-amd64.exe`（首次运行自动创建 `db/`、`template/`、`logs/`，默认端口 8000）
- 修改端口/账号：
  ```bash
  sublinke setting --username admin --password 新密码 --port 8000
  sublinke run --port 8000
  ```

## 使用提示

- 首次登录后请在「系统管理 → 注册配置」设置：默认订阅 ID、邀请码注册开关、**UA 检测开关与关键字**
- 在「用户管理」为用户分配订阅、配置允许地区；邀请码在「邀请码管理」创建（可设使用上限、可删除）
- 修改密码需要验证旧密码，修改后所有旧登录会失效（需重新登录）

# 开发与构建

## 本地构建（前端 + 后端）

```bash
# 前端
cd webs && pnpm install && pnpm build
# 后端(生产模式, 内嵌前端产物)
cd ..
cp -r webs/dist static
go build -tags=prod -ldflags "-s -w" -o sublinke .
rm -rf static
```

## 自动发布（GitHub Actions）

推送 `v*` 标签自动编译 **linux-amd64 / linux-arm64 / windows-amd64** 三个平台的二进制并上传到 GitHub Release：

```bash
git tag v1.2.0
git push origin v1.2.0
```

> 前端依赖已通过 `pnpm-lock.yaml` 锁定（typescript 5.4.2 / vue-tsc 2.0.6），CI 与本地构建保持一致。

# 插件说明

`sublinkE` 提供了灵活的插件系统（实验性），允许开发者扩展系统功能而无需修改核心代码。

## 插件开发指南

1. **创建插件文件**：参照 `plugins_examples/email_plugin.go` 编写自定义插件
2. **编译插件**：使用 `plugins_examples/build_plugin.sh email_plugin.go` 编译成 `.so` 文件
3. **部署插件**：将编译好的 `.so` 文件放入 `plugins` 目录

所有插件必须实现 `plugins.Plugin` 接口：

```go
Name() string                           // 插件名称
Version() string                        // 插件版本
Description() string                    // 插件描述
DefaultConfig() map[string]interface{}  // 默认配置
SetConfig(map[string]interface{})       // 设置配置
Init() error                            // 初始化
Close() error                           // 关闭清理

// API 事件监听
OnAPIEvent(ctx *gin.Context, event plugins.EventType, path string,
           statusCode int, requestBody interface{}, responseBody interface{}) error
InterestedAPIs() []string
InterestedEvents() []plugins.EventType
```

内置示例插件（版本更新可能失效，建议自己编译）：

| 插件名称 | 功能描述 | 源代码 |
|---------|--------|-------|
| **邮件通知插件** | 监控登录事件并发送邮件通知 | [email_plugin.go](https://github.com/eun1e/sublinkE/blob/main/plugins_examples/email_plugin.go) |

可通过 Web 界面的「插件管理」启用/禁用插件、配置参数、查看状态。

> ⚠ GitHub Actions 交叉编译产物（CGO 禁用）不包含插件加载支持；需要插件请在目标平台原生构建或使用 Docker 部署。

# 技术栈

| 层 | 技术 |
|---|---|
| 后端 | Go 1.24 + Gin + Gorm + SQLite（纯 Go 驱动） |
| 前端 | Vue 3 + TypeScript + Element Plus + Vite + Pinia |
| 认证 | JWT（HS256）+ bcrypt 密码哈希 |
| 验证码 | 自研滑块验证码（服务端校验位置 + 轨迹检测） |
| 部署 | Docker / 一键脚本 / 预编译二进制 |

# 项目预览

![预览1](webs/src/assets/1.png)
![预览2](webs/src/assets/2.png)
![预览3](webs/src/assets/3.png)
![预览4](webs/src/assets/4.png)
![预览5](webs/src/assets/5.png)
![预览6](webs/src/assets/6.png)

# 免责声明

本项目仅用于技术学习与个人自用。使用者应遵守所在国家/地区的法律法规，请勿用于任何非法用途。
