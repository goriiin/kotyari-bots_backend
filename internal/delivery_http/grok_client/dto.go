package grok_client

const (
	systemRole   = "system"
	userRole     = "user"
	defaultModel = "grok-3-mini"
)

type GrokRequest struct {
	Model    string        `json:"model"`
	Messages []GrokMessage `json:"messages"`
}

type GrokMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func messageFromPrompt(role, prompt string) GrokMessage {
	return GrokMessage{
		Role:    role,
		Content: prompt,
	}
}

func messagesFromGrokMessage(messages ...GrokMessage) []GrokMessage {
	messagesSlice := make([]GrokMessage, 0, len(messages))
	for _, message := range messages {
		messagesSlice = append(messagesSlice, message)
	}

	return messagesSlice
}

type GrokResponse struct {
	Choices []GrokChoice `json:"choices"`
	// Можно будет добавить поле токен, если нужно будет добавить rate-limiter в будущем
}

type GrokChoice struct {
	Index   int                 `json:"index"`
	Message GrokResponseMessage `json:"message"`
}

type GrokResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
