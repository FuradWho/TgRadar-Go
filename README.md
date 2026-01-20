# TgRadar-Go

[English](#english) | [中文](#chinese)

---

<a name="english"></a>
## 🇬🇧 English

### Introduction
**TgRadar-Go** is a real-time Telegram group sentiment analysis and summary tool powered by AI (OpenAI/DeepSeek). It monitors specified Telegram groups, captures chat messages, and periodically generates concise briefings, helping users quickly grasp community hotspots, sentiment trends, and high-value information.

### Features
- **Real-time Monitoring**: Connects to Telegram using MTProto to capture group messages in real-time.
- **AI Analysis**: Integrates LLM (OpenAI/DeepSeek) to automatically summarize chat content.
- **Sentiment Analysis**: Identifies community sentiment (Positive/Panic/Wait-and-see, etc.).
- **Key Info Extraction**: Automatically extracts project names, URLs, error messages, etc.
- **Proxy Support**: Supports SOCKS5 proxy for use in restricted network environments.
- **Modular Design**: Clean architecture with low coupling, easy to extend and maintain.

### Configuration (`config.yml`)

Create a `config.yml` file in the root directory:

```yaml
telegram:
  app_id: 12345678             # Your Telegram App ID
  app_hash: "your_app_hash"    # Your Telegram App Hash
  session_file: "session.json" # Session storage file path
  phone: "+1234567890"         # Your phone number
  password: "your_2fa_password"# 2FA password (if enabled)
  proxy: "127.0.0.1:10808"     # SOCKS5 proxy address (optional, leave empty if not needed)

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

### License
This project is licensed under the [MIT License](LICENSE).

---

<a name="chinese"></a>
## 🇨🇳 中文

### 项目简介
**TgRadar-Go** 是一个基于 AI (OpenAI/DeepSeek) 的 Telegram 群组实时舆情分析与摘要工具。它能实时监控指定的 Telegram 群组，抓取聊天记录，并定期生成简报，帮助用户快速了解社群热点、情绪倾向和高价值信息。

### 基础功能
- **实时监控**：基于 MTProto 协议连接 Telegram，实时捕获群组消息。
- **AI 智能分析**：集成 LLM (OpenAI/DeepSeek) 自动总结聊天内容。
- **情绪判断**：识别社群整体情绪（积极/恐慌/观望等）。
- **关键信息提取**：自动提取项目名、网址、报错信息等高价值内容。
- **代理支持**：内置 SOCKS5 代理支持，适应国内网络环境。
- **模块化设计**：代码结构清晰，低耦合，易于扩展和维护。

### 配置文件 (`config.yml`)

在项目根目录下创建 `config.yml` 文件：

```yaml
telegram:
  app_id: 12345678             # 你的 Telegram App ID
  app_hash: "your_app_hash"    # 你的 Telegram App Hash
  session_file: "session.json" # 会话保存文件路径
  phone: "+1234567890"         # 你的手机号
  password: "your_2fa_password"# 两步验证密码 (如果开启)
  proxy: "127.0.0.1:10808"     # SOCKS5 代理地址 (可选，不需要则留空)

monitor:
  window_seconds: 60           # 分析周期（秒）
  debug: true                  # 是否开启调试日志

ai:
  api_key: "sk-xxxxxx"         # AI API Key
  base_url: "https://api.deepseek.com" # API Base URL (OpenAI留空，DeepSeek等需填写)
  model: "deepseek-chat"       # 模型名称 (如 gpt-4o, deepseek-chat)
  language: "zh"               # 输出语言 (预留字段)
```

### 使用方法

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

### 开源协议
本项目采用 [MIT License](LICENSE) 开源协议。
