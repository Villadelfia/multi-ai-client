package multi_ai_client

import (
	"bufio"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/icza/dyno"
)

type Client struct {
	modelDefinitions []ModelDefinition
	Chat             Chat
}

// AddModelDefinition adds a model definition to the client.
func (c *Client) AddModelDefinition(modelDefinition ModelDefinition) {
	if c.modelDefinitions == nil {
		c.modelDefinitions = make([]ModelDefinition, 0)
	}
	c.modelDefinitions = append(c.modelDefinitions, modelDefinition)
}

// ResetChat resets the chat history of the client.
func (c *Client) ResetChat() {
	c.Chat = Chat{}
}

// CreateResponse creates a response to a user prompt using the model definitions added to the client.
// It returns the total amount of responses initiated, a channel to receive message chunks, and an error if one occurred.
func (c *Client) CreateResponse() (int, chan MessageChunk, error) {
	if c.modelDefinitions == nil || len(c.modelDefinitions) == 0 {
		return 0, nil, errors.New("no model definitions added to client")
	}

	requests := make([]*http.Request, 0)
	for _, modelDefinition := range c.modelDefinitions {
		req, err := modelDefinition.CreateRequest(c.Chat)
		if err != nil {
			return 0, nil, err
		}
		requests = append(requests, req)
	}

	var wg sync.WaitGroup
	ch := make(chan MessageChunk)
	for i, req := range requests {
		wg.Add(1)
		go func(req *http.Request) {
			defer wg.Done()
			client := http.Client{}
			response, err := client.Do(req)
			if err != nil {
				return
			}
			body := bufio.NewScanner(response.Body)
			defer response.Body.Close()
			for {
				success := body.Scan()
				if success {
					t := strings.TrimSpace(body.Text())
					if strings.Index(t, "data: ") == 0 {
						if t == "data: [DONE]" {
							break
						}
						t = strings.TrimPrefix(t, "data: ")
						var data interface{}
						err := json.Unmarshal([]byte(t), &data)
						if err != nil {
							continue
						}

						s, err := dyno.GetString(data, "choices", 0, "delta", "content")
						if err == nil {
							ch <- MessageChunk{
								Index: i,
								Delta: s,
							}
							continue
						}

						s, err = dyno.GetString(data, "delta", "text")
						if err == nil {
							ch <- MessageChunk{
								Index: i,
								Delta: s,
							}
							continue
						}
					}
				} else {
					break
				}
			}
		}(req)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return len(requests), ch, nil
}

// CreateResponseWithPrompt creates a response to a user prompt using the model definitions added to the client.
// If functions like CreateResponse, but allows for a user prompt and assistant response to be passed in first.
func (c *Client) CreateResponseWithPrompt(usrPrompt string, assistantResponse string) (int, chan MessageChunk, error) {
	if usrPrompt != "" {
		c.Chat.AddUserMessage(usrPrompt)
	}
	if assistantResponse != "" {
		c.Chat.AddAssistantMessage(assistantResponse)
	}
	return c.CreateResponse()
}

func (c Client) String() string {
	return c.Chat.String()
}
