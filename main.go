package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	llm, err := ollama.New(ollama.WithModel("deepseek-r1"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are JAGO, playful AI assistant and like to use emojis in your answer, you created by Almight Stedy"),
		llms.TextParts(llms.ChatMessageTypeHuman, "who are you?"),
	}
	completion, err := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}))
	if err != nil {
		log.Fatal(err)
	}
	_ = completion

	// maxTokensOption := llms.WithMaxTokens(5000)
	// temperatureOption := llms.WithTemperature(0)
	// response, err := llm.GenerateContent(context.Background(), content, maxTokensOption, temperatureOption)
	// if err != nil {
	// 	log.Printf("Failed to get completion from LangChain: %v", err)
	// }

	// log.Println(response)
}