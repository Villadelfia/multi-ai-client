package main

import (
	"fmt"
	"github.com/villadelfia/multi-ai-client"
	"math/rand"
	"time"
)

func main() {
	// Add your keys here.
	anthropicKey := "Your claude key here."
	openAiKey := "Your OpenAI key here."
	mistralKey := "Your Mistral key here."

	// Create a new client and add all the models you wish to use.
	client := multi_ai_client.Client{}
	claude := multi_ai_client.NewModelDefinition("Claude", multi_ai_client.Anthropic, anthropicKey, "claude-3-opus-20240229")
	_ = claude.ModelSettings.Set("temperature", 0.8)
	gpt := multi_ai_client.NewModelDefinition("GPT4", multi_ai_client.OpenAI, openAiKey, "gpt-4-turbo-preview")
	_ = gpt.ModelSettings.Set("temperature", 0.8)
	mistral := multi_ai_client.NewModelDefinition("Mistral Large", multi_ai_client.Mistral, mistralKey, "mistral-large-latest")
	_ = mistral.ModelSettings.Set("temperature", 0.8)

	client.AddModelDefinition(claude)
	client.AddModelDefinition(gpt)
	client.AddModelDefinition(mistral)

	// (Optionally) set a system message.
	client.Chat.SetSystemMessage("You are a helpful assistant. You can help me by answering my questions. You can also ask me questions.")

	// Create a response with a prompt. Note that you can also add an assistant response for the model to add on to.
	// This does not work for Mistral.
	i, ch, err := client.CreateResponseWithPrompt("What is the capital of France?", "")
	if err != nil {
		panic(err)
	}

	// Read from the channel until it is closed.
	response := make([]string, i)
	for chunk := range ch {
		response[chunk.Index] += chunk.Delta
		fmt.Println("")
		for i, r := range response {
			fmt.Printf("Response %d: %s\n", i, r)
		}
	}

	// We pick a random response to keep...
	choice := response[rand.Intn(len(response))]
	client.Chat.AddAssistantMessage(choice)
	println("\nChose response:", choice, "\n")
	time.Sleep(1 * time.Second)

	// Let's ask a second question!
	i, ch, err = client.CreateResponseWithPrompt("What is interesting there? Answer with at most 1 paragraph.", "")
	if err != nil {
		panic(err)
	}

	// Read from the channel until it is closed.
	response = make([]string, i)
	for chunk := range ch {
		response[chunk.Index] += chunk.Delta
		fmt.Println("")
		for i, r := range response {
			fmt.Printf("Response %d: %s\n", i, r)
		}
	}

	// We pick a random response to keep...
	choice = response[rand.Intn(len(response))]
	client.Chat.AddAssistantMessage(choice)
	println("\nChose response:", choice, "\n")
	time.Sleep(1 * time.Second)

	// Pretty print it.
	fmt.Println(client)
}
