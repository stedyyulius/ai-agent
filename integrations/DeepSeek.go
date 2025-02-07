package integrations

import (
	"ai-agent/knowledge"
	"context"
	"log"
	"regexp"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

var LLM *ollama.LLM

func ConnectDeepSeek() error {
	llm, err := ollama.New(ollama.WithModel("deepseek-r1"))
	if err != nil {
		log.Fatal(err)
	}

	LLM = llm

	return nil
}

func ProcessChat(prompt string, history []llms.MessageContent) (string, error) {

	// intention := CheckUserIntention(prompt)
	intention := "asking question"

	customKnowledge, err := EnrichedKnowledge(prompt, intention)
	if err != nil {
		history = append(history, llms.TextParts(llms.ChatMessageTypeSystem, err.Error()))
	}

	history = append(history, customKnowledge...)

	history = append(history, llms.TextParts(llms.ChatMessageTypeHuman, prompt))

	maxTokensOption := llms.WithMaxTokens(5000)
	temperatureOption := llms.WithTemperature(1.5)
	response, err := LLM.GenerateContent(context.Background(), history, maxTokensOption, temperatureOption)
	if err != nil {
		log.Printf("Failed to get completion from LangChain: %v", err)
		return "", nil
	}

	responseText := response.Choices[0].Content

	return responseText, nil
}

func ConvertMarkdownToHTML(text string) string {

	codeBlockRegex := regexp.MustCompile("```(.*?)```")

	text = codeBlockRegex.ReplaceAllString(text, "<pre><code>$1<code></pre>")

	boldRegex := regexp.MustCompile(`\*\*(.*?)\*\*`)
	text = boldRegex.ReplaceAllString(text, `<b>$1</b>`)

	return text
}

// func CheckUserIntention(prompt string) string {

// 	checker := []llms.MessageContent{
// 		llms.TextParts(llms.ChatMessageTypeHuman, fmt.Sprintf("What is user intention based on the prompt? "+
// 			"Answer only with one of the following options "+
// 			"and if none of the options fit user intention just answer 'asking question'."+
// 			"The options are: "+
// 			"buying ticket"+
// 			"user prompt: %s", prompt)),
// 	}

// 	maxTokensOption := llms.WithMaxTokens(50)
// 	temperatureOption := llms.WithTemperature(0)

// 	response, err := LLM.GenerateContent(context.Background(), checker, maxTokensOption, temperatureOption)
// 	if err != nil {
// 		log.Printf("Failed to get completion from intention checker: %v", err)
// 	}

// 	if len(response.Choices) == 0 {
// 		log.Println("Intention checker returned no response")
// 		return "asking question"
// 	}

// 	return response.Choices[0].Content

// }

func EnrichedKnowledge(prompt string, intention string) ([]llms.MessageContent, error) {
	var enrichedKnowledge []llms.MessageContent

	enrichedKnowledge = append(enrichedKnowledge, llms.TextParts(llms.ChatMessageTypeSystem, knowledge.Identity()))
	enrichedKnowledge = append(enrichedKnowledge, llms.TextParts(llms.ChatMessageTypeSystem, knowledge.FlatEarth()))

	return enrichedKnowledge, nil
}
