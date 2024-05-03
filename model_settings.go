package multi_ai_client

import (
	"encoding/json"
	"errors"
	"github.com/guregu/null/v5"
	"github.com/icza/dyno"
	"slices"
)

// ModelSettings is an interface representing the settings needed to interact
// with a model.
type ModelSettings interface {
	// MakeBody creates the body of the request to the model API.
	MakeBody(chat Chat) []byte

	// Set sets a value in the model settings.
	// The available keys are specific to the model settings implementation.
	// If the key is not valid, an error is returned.
	// If the value may not be nilled and a nil value is passed, an error is returned.
	//
	// For OpenAI, the valid keys are:
	//  - model                 (required, any valid model name as string)
	//  - frequency_penalty     (float64, [-2.0; 2.0])
	//  - logit_bias            (map[string]int, {"token_id": bias [-100; 100]})
	//  - logprobs				(bool)
	//  - top_logprobs			(int, [0; 20])
	//  - max_tokens			(int, [1; +inf])
	//  - presence_penalty		(float64, [-2.0; 2.0])
	//  - response_format		(string, "json" or "plain_text")
	//  - seed					(int)
	//  - stop					([]string)
	//  - temperature			(float64, [0.0; 2.0])
	//  - top_p					(float64, [0.0; 1.0])
	//  - user					(string)
	//
	// For Mistral, the valid keys are:
	//  - model                 (required, any valid model name as string)
	//  - response_format		(string, "json" or "plain_text")
	//  - temperature			(float64, [0.0; 1.0])
	//  - top_p					(float64, [0.0; 1.0])
	//  - max_tokens			(int, [1; +inf])
	//  - safe_prompt			(bool)
	//  - random_seed			(int)
	//
	// For Anthropic, the valid keys are:
	//  - model                 (required, any valid model name as string)
	//  - max_tokens			(required, int, [1; +inf])
	//  - metadata, user_id		(string)
	//  - stop_sequences		([]string)
	//  - temperature			(float64, [0.0; 1.0])
	//  - top_k					(int)
	//  - top_p					(float64, [0.0; 1.0])
	//

	Set(key string, value interface{}) error
}

type ModelSettingsOpenAI struct {
	ModelSettings    `json:"-"`
	Model            string                `json:"model"`
	FrequencyPenalty null.Float            `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int        `json:"logit_bias,omitempty"`
	Logprobs         null.Bool             `json:"logprobs,omitempty"`
	TopLogprobs      null.Int              `json:"top_logprobs,omitempty"`
	MaxTokens        null.Int              `json:"max_tokens,omitempty"`
	PresencePenalty  null.Float            `json:"presence_penalty,omitempty"`
	ResponseFormat   *OpenAIResponseFormat `json:"response_format,omitempty"`
	Seed             null.Int              `json:"seed,omitempty"`
	Stop             []string              `json:"stop,omitempty"`
	Stream           null.Bool             `json:"stream,omitempty"`
	Temperature      null.Float            `json:"temperature,omitempty"`
	TopP             null.Float            `json:"top_p,omitempty"`
	User             null.String           `json:"user,omitempty"`
	Messages         []JsonMessage         `json:"messages"`
	// Tools            []OpenAITool          `json:"tools,omitempty"`
	// ToolChoice       *OpenAIToolChoice     `json:"tool_choice,omitempty"`
}

func (m ModelSettingsOpenAI) MakeBody(chat Chat) []byte {
	messages := chat.GetMessages()
	if len(messages) > 0 {
		m.Messages = make([]JsonMessage, 0, len(messages))
		for _, message := range messages {
			m.Messages = append(m.Messages, NewJsonMessageFromMessage(message))
		}
	} else {
		m.Messages = nil
	}
	m.Stream = null.BoolFrom(true)
	body, _ := json.Marshal(m)
	return body
}

func (m ModelSettingsOpenAI) Set(key string, value interface{}) error {
	keys := []string{"model", "frequency_penalty", "logit_bias", "logprobs", "top_logprobs", "max_tokens", "presence_penalty", "response_format", "seed", "stop", "temperature", "top_p", "user"}
	if !slices.Contains(keys, key) {
		return errors.New("invalid key")
	}
	switch key {
	case "model":
		if value == nil {
			return errors.New("value may not be nil")
		}
		m.Model = value.(string)
	case "frequency_penalty":
		if value == nil {
			m.FrequencyPenalty = null.FloatFromPtr(nil)
		} else {
			m.FrequencyPenalty = null.FloatFrom(value.(float64))
		}
	case "logit_bias":
		if value == nil {
			m.LogitBias = nil
		} else {
			m.LogitBias = value.(map[string]int)
		}
	case "logprobs":
		if value == nil {
			m.Logprobs = null.BoolFromPtr(nil)
		} else {
			m.Logprobs = null.BoolFrom(value.(bool))
		}
	case "top_logprobs":
		if value == nil {
			m.TopLogprobs = null.IntFromPtr(nil)
		} else {
			m.TopLogprobs = null.IntFrom(int64(value.(int)))
		}
	case "max_tokens":
		if value == nil {
			m.MaxTokens = null.IntFromPtr(nil)
		} else {
			m.MaxTokens = null.IntFrom(int64(value.(int)))
		}
	case "presence_penalty":
		if value == nil {
			m.PresencePenalty = null.FloatFromPtr(nil)
		} else {
			m.PresencePenalty = null.FloatFrom(value.(float64))
		}
	case "seed":
		if value == nil {
			m.Seed = null.IntFromPtr(nil)
		} else {
			m.Seed = null.IntFrom(int64(value.(int)))
		}
	case "stop":
		if value == nil {
			m.Stop = nil
		} else {
			m.Stop = value.([]string)
		}
	case "temperature":
		if value == nil {
			m.Temperature = null.FloatFromPtr(nil)
		} else {
			m.Temperature = null.FloatFrom(value.(float64))
		}
	case "top_p":
		if value == nil {
			m.TopP = null.FloatFromPtr(nil)
		} else {
			m.TopP = null.FloatFrom(value.(float64))
		}
	case "user":
		if value == nil {
			m.User = null.StringFromPtr(nil)
		} else {
			m.User = null.StringFrom(value.(string))
		}
	case "response_format":
		if value == nil {
			m.ResponseFormat = nil
		} else {
			m.ResponseFormat = &OpenAIResponseFormat{Type: value.(string)}
		}
	}
	return nil
}

type ModelSettingsMistral struct {
	ModelSettings  `json:"-"`
	Model          string                 `json:"model"`
	ResponseFormat *MistralResponseFormat `json:"response_format,omitempty"`
	Temperature    null.Float             `json:"temperature,omitempty"`
	TopP           null.Float             `json:"top_p,omitempty"`
	MaxTokens      null.Int               `json:"max_tokens,omitempty"`
	Stream         null.Bool              `json:"stream,omitempty"`
	SafePrompt     null.Bool              `json:"safe_prompt,omitempty"`
	RandomSeed     null.Int               `json:"random_seed,omitempty"`
	Messages       []JsonMessage          `json:"messages"`
	// Tools            []MistralTool     `json:"tools,omitempty"`
	// ToolChoice       null.String       `json:"tool_choice,omitempty"`
}

func (m ModelSettingsMistral) MakeBody(chat Chat) []byte {
	messages := chat.GetMessages()
	if len(messages) > 0 {
		m.Messages = make([]JsonMessage, 0, len(messages))
		for _, message := range messages {
			m.Messages = append(m.Messages, NewJsonMessageFromMessage(message))
		}
	} else {
		m.Messages = nil
	}
	m.Stream = null.BoolFrom(true)

	// We need to alter the body to remove the null values
	body, _ := json.Marshal(m)
	var altered interface{}
	_ = json.Unmarshal(body, &altered)
	_, err := dyno.GetFloating(altered, "temperature")
	if err != nil {
		_ = dyno.Delete(altered, "temperature")
	}
	_, err = dyno.GetFloating(altered, "top_p")
	if err != nil {
		_ = dyno.Delete(altered, "top_p")
	}
	_, err = dyno.GetInt(altered, "max_tokens")
	if err != nil {
		_ = dyno.Delete(altered, "max_tokens")
	}
	_, err = dyno.GetBoolean(altered, "safe_prompt")
	if err != nil {
		_ = dyno.Delete(altered, "safe_prompt")
	}
	_, err = dyno.GetInt(altered, "random_seed")
	if err != nil {
		_ = dyno.Delete(altered, "random_seed")
	}
	body, _ = json.Marshal(altered)
	return body
}

func (m ModelSettingsMistral) Set(key string, value interface{}) error {
	keys := []string{"model", "response_format", "temperature", "top_p", "max_tokens", "safe_prompt", "random_seed"}
	if !slices.Contains(keys, key) {
		return errors.New("invalid key")
	}
	switch key {
	case "model":
		if value == nil {
			return errors.New("value may not be nil")
		}
		m.Model = value.(string)
	case "response_format":
		if value == nil {
			m.ResponseFormat = nil
		} else {
			m.ResponseFormat = &MistralResponseFormat{Type: value.(string)}
		}
	case "temperature":
		if value == nil {
			m.Temperature = null.FloatFromPtr(nil)
		} else {
			m.Temperature = null.FloatFrom(value.(float64))
		}
	case "top_p":
		if value == nil {
			m.TopP = null.FloatFromPtr(nil)
		} else {
			m.TopP = null.FloatFrom(value.(float64))
		}
	case "max_tokens":
		if value == nil {
			m.MaxTokens = null.IntFromPtr(nil)
		} else {
			m.MaxTokens = null.IntFrom(int64(value.(int)))
		}
	case "safe_prompt":
		if value == nil {
			m.SafePrompt = null.BoolFromPtr(nil)
		} else {
			m.SafePrompt = null.BoolFrom(value.(bool))
		}
	case "random_seed":
		if value == nil {
			m.RandomSeed = null.IntFromPtr(nil)
		} else {
			m.RandomSeed = null.IntFrom(int64(value.(int)))
		}
	}
	return nil
}

type ModelSettingsAnthropic struct {
	ModelSettings `json:"-"`
	Model         string             `json:"model"`
	MaxTokens     int                `json:"max_tokens"`
	Metadata      *AnthropicMetadata `json:"metadata,omitempty"`
	StopSequences []string           `json:"stop_sequences,omitempty"`
	Stream        null.Bool          `json:"stream,omitempty"`
	Temperature   null.Float         `json:"temperature,omitempty"`
	TopK          null.Int           `json:"top_k,omitempty"`
	TopP          null.Float         `json:"top_p,omitempty"`
	Messages      []JsonMessage      `json:"messages"`
	System        null.String        `json:"system,omitempty"`
	// Tools         []AnthropicTool    `json:"tools,omitempty"`
}

func (m ModelSettingsAnthropic) MakeBody(chat Chat) []byte {
	if chat.GetSystemMessage() != "" {
		m.System = null.StringFrom(chat.GetSystemMessage())
	} else {
		m.System = null.StringFromPtr(nil)
	}

	messages := chat.GetMessagesWithoutSystemMessage()
	if len(messages) > 0 {
		m.Messages = make([]JsonMessage, 0, len(messages))
		for _, message := range messages {
			m.Messages = append(m.Messages, NewJsonMessageFromMessage(message))
		}
	} else {
		m.Messages = nil
	}

	m.Stream = null.BoolFrom(true)
	if m.MaxTokens == 0 {
		m.MaxTokens = 4096
	}

	// We need to alter the body to remove the null values
	body, _ := json.Marshal(m)
	var altered interface{}
	_ = json.Unmarshal(body, &altered)
	_, err := dyno.GetFloating(altered, "temperature")
	if err != nil {
		_ = dyno.Delete(altered, "temperature")
	}
	_, err = dyno.GetFloating(altered, "top_p")
	if err != nil {
		_ = dyno.Delete(altered, "top_p")
	}
	_, err = dyno.GetInt(altered, "top_k")
	if err != nil {
		_ = dyno.Delete(altered, "top_k")
	}
	_, err = dyno.GetString(altered, "system")
	if err != nil {
		_ = dyno.Delete(altered, "system")
	}
	body, _ = json.Marshal(altered)
	return body
}

func (m ModelSettingsAnthropic) Set(key string, value interface{}) error {
	keys := []string{"model", "max_tokens", "metadata", "stop_sequences", "temperature", "top_k", "top_p"}
	if !slices.Contains(keys, key) {
		return errors.New("invalid key")
	}
	switch key {
	case "model":

		if value == nil {
			return errors.New("value may not be nil")
		}
		m.Model = value.(string)
	case "max_tokens":
		if value == nil {
			return errors.New("value may not be nil")
		}
		m.MaxTokens = value.(int)
	case "metadata", "user_id":
		if value == nil {
			m.Metadata = nil
		} else {
			m.Metadata = &AnthropicMetadata{UserID: value.(string)}
		}
	case "stop_sequences":
		if value == nil {
			m.StopSequences = nil
		} else {
			m.StopSequences = value.([]string)
		}
	case "temperature":
		if value == nil {
			m.Temperature = null.FloatFromPtr(nil)
		} else {
			m.Temperature = null.FloatFrom(value.(float64))
		}
	case "top_k":
		if value == nil {
			m.TopK = null.IntFromPtr(nil)
		} else {
			m.TopK = null.IntFrom(int64(value.(int)))
		}
	case "top_p":
		if value == nil {
			m.TopP = null.FloatFromPtr(nil)
		} else {
			m.TopP = null.FloatFrom(value.(float64))
		}
	}
	return nil
}

type OpenAIResponseFormat struct {
	Type string `json:"type"`
}

type MistralResponseFormat struct {
	Type string `json:"type"`
}

type AnthropicMetadata struct {
	UserID string `json:"user_id"`
}

type JsonMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewJsonMessageFromMessage(message Message) JsonMessage {
	role := ""
	switch message.Type {
	case UserMessage:
		role = "user"
	case AssistantMessage:
		role = "assistant"
	case SystemMessage:
		role = "system"
	}
	return JsonMessage{
		Role:    role,
		Content: message.Text,
	}
}

func NewModelSettings(apiType APIType, modelName string) ModelSettings {
	switch apiType {
	case OpenAI:
		return ModelSettingsOpenAI{
			Model: modelName,
		}
	case Mistral:
		return ModelSettingsMistral{
			Model: modelName,
		}
	case Anthropic:
		return ModelSettingsAnthropic{
			Model: modelName,
		}
	default:
		return nil
	}
}
