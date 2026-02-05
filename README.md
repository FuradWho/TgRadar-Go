# TgRadar-Go

[English](README.md) | [中文](README.zh-CN.md)

### Introduction
**TgRadar-Go** is an AI-powered Telegram group monitoring and briefing tool (OpenAI/DeepSeek). It captures messages from selected groups and periodically generates a consolidated market brief with trading sentiment, hot topics, and key signals.

### Features
- **Multi-group monitoring**: Track multiple groups or all groups.
- **Periodic AI briefing**: Generates a single consolidated summary per window.
- **Trading-focused insights**: Highlights sentiment, hot projects, and key events.
- **Telegram Bot delivery**: Pushes summaries to a Bot chat.
- **Proxy support**: SOCKS5 proxy for restricted networks.
- **Clean architecture**: Modular design, easy to extend.

### Configuration (`config.yml`)

Create a `config.yml` file in the root directory:

```yaml
telegram:
  app_id: 12345678             # Your Telegram App ID
  app_hash: "your_app_hash"    # Your Telegram App Hash
  session_file: "session.json" # Session storage file path
  phone: "+1234567890"         # Your phone number
  password: "your_2fa_password"# 2FA password (if enabled)
  proxy: "127.0.0.1:10808"     # SOCKS5 proxy address (optional)
  target_groups: [1234567890]  # Target group IDs (empty = all)
  bot_token: "123456:ABCDEF"   # Bot token (optional)
  bot_chat_id: -1001234567890  # Bot target chat_id (optional)

monitor:
  window_seconds: 60           # Analysis interval (seconds)
  debug: true                  # Enable debug logs

ai:
  api_key: "sk-xxxxxx"         # Your AI API Key
  base_url: "https://api.deepseek.com" # API Base URL (optional, e.g., for DeepSeek)
  model: "deepseek-chat"       # Model name (e.g., gpt-4o, deepseek-chat)
  language: "en"               # Output language (reserved for future use)
```

### Usage

1.  **Prerequisites**:
    *   Go 1.21+ installed.
    *   Telegram API ID & Hash (Get from [my.telegram.org](https://my.telegram.org)).
    *   AI API Key (OpenAI or DeepSeek).

2.  **Run**:
    ```bash
    # Clone the repository
    git clone https://github.com/FuradWho/TgRadar-Go.git
    cd TgRadar-Go

    # Install dependencies
    go mod tidy

    # Run the application
    go run .
    ```

3.  **Login**:
    *   On the first run, the terminal will prompt you to enter the Telegram verification code sent to your app.
4.  **Bot delivery (optional)**:
    *   Set `bot_token` and `bot_chat_id` to receive summaries in Telegram.

### License
This project is licensed under the [MIT License](LICENSE).

