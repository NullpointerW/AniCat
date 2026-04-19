# Telegram Bot 推送配置指南

AniCat 支持通过 Telegram Bot 在剧集下载完成后发送通知。本文档介绍从创建 Bot 到接收通知的完整配置流程。

---

## 前置条件

- 一个 Telegram 账号
- 能够访问 Telegram（如需代理，见[代理配置](#代理配置)）

---

## 第一步：创建 Bot 并获取 Token

1. 在 Telegram 中搜索并打开 [@BotFather](https://t.me/BotFather)
2. 发送 `/newbot`
3. 按提示输入 Bot 的**显示名称**（例如 `AniCat 推送`）
4. 再输入 Bot 的**用户名**，必须以 `bot` 结尾（例如 `anicat_notify_bot`）
5. 创建成功后，BotFather 会回复一段消息，其中包含 Token，格式如下：

   ```
   Use this token to access the HTTP API:
   1234567890:ABCDEFabcdefABCDEFabcdefABCDEF12345
   ```

   复制冒号后面的完整字符串，即为 `token`。

---

## 第二步：获取 Chat ID

### 个人私聊

1. 在 Telegram 中搜索并打开刚创建的 Bot，点击 **Start**
2. 向 Bot 发送任意一条消息（例如 `hello`）
3. 在浏览器访问以下地址（将 `<token>` 替换为你的 Token）：

   ```
   https://api.telegram.org/bot<token>/getUpdates
   ```

4. 返回的 JSON 中找到 `message.chat.id` 字段，该数字即为 `chat_id`：

   ```json
   {
     "result": [{
       "message": {
         "chat": {
           "id": 123456789,
           ...
         }
       }
     }]
   }
   ```

### 群组推送

1. 将 Bot 添加到目标群组，并赋予发送消息权限
2. 在群组中发送一条消息 `@你的bot用户名 test`
3. 同样访问上方 `getUpdates` 地址，群组的 `chat.id` 为**负数**，例如 `-1001234567890`

---

## 第三步：配置 env.yaml

在 `env.yaml` 的 `push` 节点下添加 `telegram` 配置：

```yaml
push:
  telegram:
    token: "1234567890:ABCDEFabcdefABCDEFabcdefABCDEF12345"
    chat_id: "123456789"
  # 邮件推送可与 Telegram 同时启用，互不影响
  email:
    host: "smtp.example.com"
    port: 465
    username: "your@email.com"
    password: "yourpassword"
```

> `token` 或 `chat_id` 为空时，Telegram 推送自动禁用，不影响其他推送方式。

---

## 第四步：启动验证

启动 AniCat 后，日志中出现以下内容说明 Bot 初始化成功：

```
telegram bot init completed  chat_id=123456789
```

若显示：

```
telegram push disable
```

说明 `token` 未配置，请检查 `env.yaml`。

---

## 通知格式

下载完成后，Bot 将发送如下格式的消息：

```
[AniCat] 剧集更新提醒
名称: 我心里危险的东西
剧集: S01E04
文件: 我心里危险的东西 S01E04.mp4
大小: 342 MB
```

---

## 代理配置

若服务器无法直接访问 `api.telegram.org`，可在 `env.yaml` 中配置 HTTP 代理：

```yaml
proxy:
  http: "http://127.0.0.1:7897"
  https: "http://127.0.0.1:7897"
```

> 当前 Telegram pusher 使用的是 Go 默认 HTTP Client，会自动读取环境变量 `HTTP_PROXY` / `HTTPS_PROXY`。也可在 `.claude/settings.json` 或系统环境变量中设置。

---

## 常见问题

| 现象 | 原因 | 解决方法 |
|------|------|----------|
| Bot 无响应，日志无报错 | Bot 未收到消息，`getUpdates` 为空 | 先向 Bot 发送一条消息再查 `getUpdates` |
| `telegram push: unexpected status 400` | `chat_id` 格式错误 | 确认 `chat_id` 为纯数字字符串，群组需带负号 |
| `telegram push: unexpected status 401` | Token 无效 | 重新从 BotFather 复制 Token |
| `telegram push: context deadline exceeded` | 网络不通 | 配置代理或检查防火墙 |
| 收到通知但 Email 未收到 | 两者独立，互不影响 | 单独检查 Email 配置 |
