# TgRadar-Go

[English](README.md) | [中文](README.zh-CN.md)

---

## 项目简介
**TgRadar-Go** 是一个基于 AI (OpenAI/DeepSeek) 的 Telegram 群组监控与简报工具。它抓取指定群聊消息，按时间窗口生成一份汇总简报，聚焦交易情绪、热点项目与关键事件。

## 基础功能
- **多群监控**：可配置多个群组，或监控所有群。
- **周期汇总**：每个窗口输出一份汇总简报。
- **交易视角**：突出情绪、热点项目与关键事件。
- **Bot 推送**：通过 Telegram Bot 自动发送汇总。
- **代理支持**：内置 SOCKS5 代理。
- **模块化设计**：结构清晰，易扩展。

## 配置文件 (`config.yml`)

在项目根目录下创建 `config.yml` 文件：

```yaml
telegram:
  app_id: 12345678             # 你的 Telegram App ID
  app_hash: "your_app_hash"    # 你的 Telegram App Hash
  session_file: "session.json" # 会话保存文件路径
  phone: "+1234567890"         # 你的手机号
  password: "your_2fa_password"# 两步验证密码 (如果开启)
  proxy: "127.0.0.1:10808"     # SOCKS5 代理地址 (可选)
  target_groups: [1234567890]  # 目标群组ID (留空则监控所有)
  bot_token: "123456:ABCDEF"   # Bot token (可选)
  bot_chat_id: -1001234567890  # Bot 接收 chat_id (可选)

monitor:
  window_seconds: 60           # 分析周期（秒）
  debug: true                  # 是否开启调试日志

ai:
  api_key: "sk-xxxxxx"         # AI API Key
  base_url: "https://api.deepseek.com" # API Base URL (OpenAI留空，DeepSeek等需填写)
  model: "deepseek-chat"       # 模型名称 (如 gpt-4o, deepseek-chat)
  language: "zh"               # 输出语言 (预留字段)
```

## 使用方法

1.  **准备工作**：
    *   安装 Go 1.21+ 环境。
    *   获取 Telegram API ID & Hash (前往 [my.telegram.org](https://my.telegram.org) 申请)。
    *   获取 AI API Key (OpenAI 或 DeepSeek)。

2.  **运行程序**：
    ```bash
    # 克隆项目
    git clone https://github.com/FuradWho/TgRadar-Go.git
    cd TgRadar-Go

    # 安装依赖
    go mod tidy

    # 启动程序
    go run .
    ```

3.  **首次登录**：
    *   程序首次运行会提示输入 Telegram 验证码（发送到你的 TG 客户端）。
4.  **Bot 推送（可选）**：
    *   配置 `bot_token` 与 `bot_chat_id`，即可在 Telegram 中接收汇总。

## 开源协议
本项目采用 [MIT License](LICENSE) 开源协议。
