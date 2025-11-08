package grok

const (
	systemRole = "system"
	userRole   = "user"
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
	return append([]GrokMessage(nil), messages...)
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
