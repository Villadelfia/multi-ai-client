package multi_ai_client

import (
	"bytes"
	"errors"
	"net/http"
)

// ModelDefinition is a struct representing a model definition.
// It contains a friendly name for the model and all the settings needed
// to interact with the model.
type ModelDefinition struct {
	Name          string
	APISettings   APISettings
	ModelSettings ModelSettings
}

// NewModelDefinition creates a new ModelDefinition with the given name, API type, API key, and model name.
// By default, the API endpoint is left as the default value and only the model name is set in the settings.
func NewModelDefinition(name string, apiType APIType, apiKey string, modelName string) ModelDefinition {
	return ModelDefinition{
		Name: name,
		APISettings: APISettings{
			APIKey:      apiKey,
			APIEndpoint: "",
			APIType:     apiType,
		},
		ModelSettings: NewModelSettings(apiType, modelName),
	}
}

func (m *ModelDefinition) CreateRequest(chat Chat) (*http.Request, error) {
	url := ""
	if m.APISettings.APIEndpoint != "" {
		url = m.APISettings.APIEndpoint
	} else {
		switch m.APISettings.APIType {
		case OpenAI:
			url = "https://api.openai.com/v1/chat/completions"
		case Mistral:
			url = "https://api.mistral.ai/v1/chat/completions"
		case Anthropic:
			url = "https://api.anthropic.com/v1/messages"
		default:
			return nil, errors.New("invalid API type")
		}
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(m.ModelSettings.MakeBody(chat)))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	if m.APISettings.APIType == Anthropic {
		request.Header.Set("anthropic-version", "2023-06-01")
	}
	if m.APISettings.APIKey != "" {
		if m.APISettings.APIType == Anthropic {
			request.Header.Set("x-api-key", m.APISettings.APIKey)
		} else {
			request.Header.Set("Authorization", "Bearer "+m.APISettings.APIKey)
		}
	}
	return request, nil
}
