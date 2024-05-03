package multi_ai_client

// MessageType is an enum representing the type of a message in a chat.
// It can be a SystemMessage, UserMessage, or AssistantMessage.
// Chats usually start with a SystemMessage, followed by an alternating
// sequence of UserMessage and AssistantMessage.
type MessageType int

const (
	SystemMessage MessageType = iota
	UserMessage
	AssistantMessage
)
