package multi_ai_client

// APIType is an enum representing the type of API.
// It determines how the client interacts with the API, and which settings are
// required or supported.
type APIType int

const (
	OpenAI APIType = iota
	Mistral
	Anthropic
)
