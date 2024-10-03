package aianswer

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"log"
	"net/http"
	"os"
)

// ChatGPT处理函数
func ChatGPTHandler(c *gin.Context) {
	var req struct {
		Prompt string `json:"prompt"`
	}
	// 绑定JSON输入
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 从环境变量中获取API密钥
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("API Key not set")
	}

	// 初始化 OpenAI 客户端
	client := openai.NewClient(apiKey)

	// 调用 OpenAI API
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo, // 或者其他模型
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: req.Prompt,
			},
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to communicate with ChatGPT"})
		return
	}

	// 返回 ChatGPT 的回复
	c.JSON(http.StatusOK, gin.H{"response": resp.Choices[0].Message.Content})
}
