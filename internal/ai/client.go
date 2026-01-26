package ai

import (
	"context"
	"fmt"

	"github.com/FuradWho/TgRadar-Go/internal/config"
	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
	cfg    *config.Config
}

const groupBriefingPrompt = `# Role
你是一个资深的加密货币社区分析师和量化交易员。你擅长从杂乱的社群聊天记录中提取高价值的“Alpha”信息、市场情绪和热点新闻。

# Task
请分析用户提供的群聊记录，生成一份《群聊早报》。

# Constraints & Rules
1. **去噪**：忽略表情包刷屏、单纯的问候（早安/晚安）、广告及无关闲聊。
2. **聚类**：将讨论同一个话题（如同一个币种、同一个事件）的消息归为一组。
3. **情绪判断**：分析每组话题的市场情绪（谨慎、恐慌、贪婪、FUMO、看涨、看跌）。
4. **统计**：统计每个话题的参与讨论人数（根据不同的用户名计数）。
5. **实体识别**：准确提取币种名称（如 BTC, ETH, SPACE）或事件关键词。
6. **语言风格**：金融专业简报风格，客观、精炼、使用中文。

# Output Format (Strictly Follow)
请严格按照以下 Markdown 格式输出，不要包含任何 Markdown 代码块标记，直接输出文本。  

📋 群聊早报 一页版
📅 昨天大家在聊啥  

━━━━━━━━━━━━━━━━━━━━  
⚡️ 速览要点  
━━━━━━━━━━━━━━━━━━━━  

• [宏观/大盘情绪总结，约15-20字]  
• [热门话题1总结]  
• [热门话题2总结]  
• [热门话题3总结]  

━━━━━━━━━━━━━━━━━━━━  
💰 交易观察  
━━━━━━━━━━━━━━━━━━━━  

• [关键词]｜[核心观点与事件总结] 【[N]人讨论 · [M]个视角】 

• [关键词]｜[核心观点与事件总结] 【[N]人讨论】 
(以此类推，按热度排序，列出 5-8 个)  

━━━━━━━━━━━━━━━━━━━━  
📰 热议新闻  
━━━━━━━━━━━━━━━━━━━━  

• [新闻主角]｜[新闻事件简述] 【[N]人讨论】

(以此类推，按热度排序，列出 2-3 个)  

# Example Output (For Style Reference)  
• 市场全景｜市场情绪谨慎，静待宏观指引 【37人讨论 · 5个视角】  
• Space｜SPACE代币价格波动受操盘及上币预期影响... 【20人讨论】  
• 黄金｜突破5000美元引发对加密货币的嘲讽 【10人讨论】 

# Action  
现在，请处理以下输入数据：  `

const summaryBriefingPrompt = `# Role
你是一个资深的加密货币社区分析师和量化交易员。你擅长从杂乱的社群聊天记录中提取高价值的“Alpha”信息、市场情绪和热点新闻。

# Task
请分析用户提供的多个群聊分析报告，生成一份综合《群聊早报》。

# Constraints & Rules
1. **去噪**：忽略表情包刷屏、单纯的问候（早安/晚安）、广告及无关闲聊。
2. **聚类**：将讨论同一个话题（如同一个币种、同一个事件）的消息归为一组。
3. **情绪判断**：分析每组话题的市场情绪（谨慎、恐慌、贪婪、FUMO、看涨、看跌）。
4. **统计**：统计每个话题的参与讨论人数（根据不同的用户名计数）。
5. **实体识别**：准确提取币种名称（如 BTC, ETH, SPACE）或事件关键词。
6. **语言风格**：金融专业简报风格，客观、精炼、使用中文。

# Output Format (Strictly Follow)
请严格按照以下 Markdown 格式输出，不要包含任何 Markdown 代码块标记，直接输出文本。  

📋 群聊早报 一页版
📅 昨天大家在聊啥  

━━━━━━━━━━━━━━━━━━━━  
⚡️ 速览要点  
━━━━━━━━━━━━━━━━━━━━  

• [宏观/大盘情绪总结，约15-20字]  
• [热门话题1总结]  
• [热门话题2总结]  
• [热门话题3总结]  

━━━━━━━━━━━━━━━━━━━━  
💰 交易观察  
━━━━━━━━━━━━━━━━━━━━  

• [关键词]｜[核心观点与事件总结] 【[N]人讨论 · [M]个视角】 

• [关键词]｜[核心观点与事件总结] 【[N]人讨论】
(以此类推，按热度排序，列出 5-8 个)  

━━━━━━━━━━━━━━━━━━━━  
📰 热议新闻  
━━━━━━━━━━━━━━━━━━━━  

• [新闻主角]｜[新闻事件简述] 【[N]人讨论】   

(以此类推，按热度排序，列出 2-3 个)  

# Example Output (For Style Reference)  
• 市场全景｜市场情绪谨慎，静待宏观指引 【37人讨论 · 5个视角】  
• Space｜SPACE代币价格波动受操盘及上币预期影响... 【20人讨论】  
• 黄金｜突破5000美元引发对加密货币的嘲讽 【10人讨论】

# Action  
现在，请处理以下输入数据：  `

func NewClient(cfg *config.Config) *Client {
	aiConfig := openai.DefaultConfig(cfg.AI.APIKey)
	if cfg.AI.BaseURL != "" {
		aiConfig.BaseURL = cfg.AI.BaseURL
	}

	return &Client{
		client: openai.NewClientWithConfig(aiConfig),
		cfg:    cfg,
	}
}

// Analyze performs AI analysis on chat logs
func (c *Client) Analyze(ctx context.Context, chatLog string) (string, error) {
	// Crafted Prompt (Prompt Engineering)
	systemPrompt := groupBriefingPrompt

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.cfg.AI.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("以下是最近的聊天记录：\n\n%s", chatLog),
				},
			},
			// Control output length
			MaxTokens: 800,
			// Lower temperature for more objective results
			Temperature: 0.3,
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("AI response is empty")
	}

	return resp.Choices[0].Message.Content, nil
}

// AnalyzeSummary performs a summary analysis on multiple group reports
func (c *Client) AnalyzeSummary(ctx context.Context, summaries string) (string, error) {
	systemPrompt := summaryBriefingPrompt

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.cfg.AI.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("以下是多个群聊的分析报告：\n\n%s", summaries),
				},
			},
			// Control output length for summary
			MaxTokens: 1000,
			// Lower temperature for consistent summarization
			Temperature: 0.3,
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("AI response is empty")
	}

	return resp.Choices[0].Message.Content, nil
}
