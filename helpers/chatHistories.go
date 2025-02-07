package helpers

import (
	"log"
	"sync"
	"time"

	"github.com/tmc/langchaingo/llms"
)

var chatHistories sync.Map


type MessageHistory struct {
	Messages  []string
	Timestamp time.Time
}

func init() {
	go CleanupExpiredMessages()
}

func StoreChatHistory(key string, message string) {
	now := time.Now()

	historyInterface, found := chatHistories.Load(key)
	if found {
		if history, ok := historyInterface.(MessageHistory); ok {
			history.Messages = append(history.Messages, message)
			history.Timestamp = now
			chatHistories.Store(key, history)
			return
		}
	}

	chatHistories.Store(key, MessageHistory{
		Messages:  []string{message},
		Timestamp: now,
	})
}

func RetrieveChatHistory(to string) []llms.MessageContent {
	historyInterface, found := chatHistories.Load(to)
	if found {
		if history, ok := historyInterface.(MessageHistory); ok {
			var chatHistory []llms.MessageContent
			log.Println("Retrieved history for:", to)
			for _, msg := range history.Messages {
				log.Println("Stored message:", msg)
				chatHistory = append(chatHistory, llms.TextParts(llms.ChatMessageTypeHuman, msg))
			}
			return chatHistory
		}
	}
	log.Println("No history found for:", to)
	return nil
}

func CleanupExpiredMessages() {
	for {
		time.Sleep(1 * time.Minute)

		now := time.Now()
		chatHistories.Range(func(key, value interface{}) bool {
			if history, ok := value.(MessageHistory); ok {
				// Remove messages older than 10 minutes
				if now.Sub(history.Timestamp) > 10*time.Minute {
					chatHistories.Delete(key)
					log.Printf("Deleted chat history for %v due to expiration.", key)
				}
			}
			return true
		})
	}
}
