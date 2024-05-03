package multi_ai_client

import (
	"errors"
	"strconv"
	"strings"
)

// Chat is a struct representing a chat between a user and an assistant.
type Chat struct {
	systemMessage *Message
	messages      []Message
}

// SetSystemMessage Adds a system message to the chat.
func (c *Chat) SetSystemMessage(s string) {
	if s == "" {
		c.systemMessage = nil
		return
	}

	if c.systemMessage == nil {
		c.systemMessage = NewSystemMessage(s)
		return
	} else {
		c.systemMessage.Text = s
	}
}

// ClearSystemMessage Removes the system message from the chat.
func (c *Chat) ClearSystemMessage() {
	c.systemMessage = nil
}

// GetSystemMessage returns the system message in the chat.
func (c *Chat) GetSystemMessage() string {
	if c.systemMessage == nil {
		return ""
	}
	return c.systemMessage.Text
}

// AddUserMessage Adds a user message to the chat.
func (c *Chat) AddUserMessage(s string) {
	if c.messages == nil {
		c.messages = make([]Message, 0)
	}
	c.messages = append(c.messages, *NewUserMessage(s))
}

// AddAssistantMessage Adds an assistant message to the chat.
func (c *Chat) AddAssistantMessage(s string) {
	if c.messages == nil {
		c.messages = make([]Message, 0)
	}
	c.messages = append(c.messages, *NewAssistantMessage(s))
}

// ReplaceLastAssistantMessage Replaces the last assistant message in the chat with a new message.
// If the last message is not an assistant message, this function does nothing.
func (c *Chat) ReplaceLastAssistantMessage(s string) {
	if len(c.messages) == 0 {
		return
	}

	if c.messages[len(c.messages)-1].Type == AssistantMessage {
		c.messages[len(c.messages)-1].Text = s
	}
}

// ClearMessages Removes all non-system messages from the chat.
func (c *Chat) ClearMessages() {
	c.messages = nil
}

// GetMessages returns all messages in the chat.
func (c *Chat) GetMessages() []Message {
	messages := make([]Message, 0)
	if c.systemMessage != nil {
		messages = append(messages, *c.systemMessage)
	}
	if c.messages != nil {
		messages = append(messages, c.messages...)
	}
	return messages
}

// GetMessagesWithoutSystemMessage returns all messages in the chat except the system message.
func (c *Chat) GetMessagesWithoutSystemMessage() []Message {
	return c.messages
}

// NewChatFromMessages creates a new Chat from a list of messages.
// The messages are added to the chat in the order they are provided.
// There may only be one system message in the list of messages, and it must be the first message.
func NewChatFromMessages(messages []Message) (*Chat, error) {
	chat := &Chat{}
	if len(messages) == 0 {
		return chat, nil
	}

	idx := 0
	if messages[0].Type == SystemMessage {
		idx = 1
		chat.SetSystemMessage(messages[0].Text)
	}

	for i := idx; i < len(messages); i++ {
		if messages[i].Type == UserMessage {
			chat.AddUserMessage(messages[i].Text)
		} else if messages[i].Type == AssistantMessage {
			chat.AddAssistantMessage(messages[i].Text)
		} else {
			return nil, errors.New("system messages must be the first message in the list of messages")
		}
	}

	return chat, nil
}

func (c *Chat) String() string {
	i := 1
	str := ""
	if c.systemMessage != nil {
		str += "# Message: " + strconv.Itoa(i) + "\n# Type: System\n"
		str += c.systemMessage.Text + "\n\n"
		i += 1
	}
	if c.messages == nil {
		return strings.TrimSpace(str)
	}
	for _, m := range c.messages {
		str += "# Message: " + strconv.Itoa(i) + "\n"
		if m.Type == UserMessage {
			str += "# Type: User\n"
		} else {
			str += "# Type: Assistant\n"
		}
		str += m.Text + "\n\n"
		i += 1
	}
	return strings.TrimSpace(str)
}
