package multi_ai_client

// Message is a struct representing a message in a chat.
// It has a type and a text.
type Message struct {
	Type MessageType
	Text string
}

// NewMessage Creates a new Message with the given type and text.
func NewMessage(t MessageType, text string) *Message {
	return &Message{
		Type: t,
		Text: text,
	}
}

// NewSystemMessage Creates a new Message with the given text of type
// SystemMessage.
func NewSystemMessage(text string) *Message {
	return NewMessage(SystemMessage, text)
}

// NewUserMessage Creates a new Message with the given text of type
// UserMessage.
func NewUserMessage(text string) *Message {
	return NewMessage(UserMessage, text)
}

// NewAssistantMessage Creates a new Message with the given text of type
// AssistantMessage.
func NewAssistantMessage(text string) *Message {
	return NewMessage(AssistantMessage, text)
}
