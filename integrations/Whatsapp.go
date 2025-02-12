package integrations

import (
	"ai-agent/helpers"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

type WhatsAppWebhook struct {
	EventType   string `json:"event_type"`
	InstanceID  string `json:"instanceId"`
	ID          string `json:"id"`
	ReferenceID string `json:"referenceId"`
	Hash        string `json:"hash"`
	Data        struct {
		ID           string   `json:"id"`
		From         string   `json:"from"`
		To           string   `json:"to"`
		Author       string   `json:"author"`
		PushName     string   `json:"pushname"`
		Ack          string   `json:"ack"`
		Type         string   `json:"type"`
		Body         string   `json:"body"`
		Media        string   `json:"media"`
		FromMe       bool     `json:"fromMe"`
		Self         bool     `json:"self"`
		IsForwarded  bool     `json:"isForwarded"`
		IsMentioned  bool     `json:"isMentioned"`
		QuotedMsg    struct{} `json:"quotedMsg"`
		MentionedIDs []string `json:"mentionedIds"`
		Time         int64    `json:"time"`
	} `json:"data"`
}

func ListenToWhatsapp(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var webhook WhatsAppWebhook
	err := json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	data := webhook.Data
	removedTaggedMessage := strings.ReplaceAll(data.Body, os.Getenv("WHATSAPP_NUMBER"), "")
	// convertNumberToName := strings.ReplaceAll(data.Body, data.From), "")

	var additionalKnowledge []llms.MessageContent

	additionalKnowledge = helpers.RetrieveChatHistory(data.To)

	if (data.PushName == "Std") {
		additionalKnowledge = append(additionalKnowledge, llms.TextParts(llms.ChatMessageTypeSystem, "Orang yang berbicara ke kamu adalah yang mulia Stedy, bicaralah dengan sopan, dan selalu panggil dia Yang Mulia Stedy dia setiap jawabanmu"))
	}

	if strings.Contains(data.From, "g.us") && data.IsMentioned {

		responseText, err := ProcessChat(removedTaggedMessage, additionalKnowledge)

		if err != nil {
			SendWhatsappMessage(data.From, "Error: "+err.Error())
			return
		}

		err = SendWhatsappMessage(data.From, responseText)
		if err != nil {
			log.Printf("Failed to send reply: %v\n", err)
		}

		helpers.StoreChatHistory(data.To, removedTaggedMessage)
		helpers.StoreChatHistory(data.To, responseText)

	} else if strings.Contains(data.From, "c.us") {

		responseText, err := ProcessChat(removedTaggedMessage, additionalKnowledge)
		if err != nil {
			SendWhatsappMessage(data.From, "Error: "+err.Error())
			return
		}

		err = SendWhatsappMessage(data.From, responseText)
		if err != nil {
			log.Printf("Failed to send private reply: %v", err)
		}

		helpers.StoreChatHistory(data.From, removedTaggedMessage)
		helpers.StoreChatHistory(data.From, responseText)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message received"))
}

func SendWhatsappMessage(receiverID, replyMessage string) error {
	instanceID := os.Getenv("WHATSAPP_INSTANCE_ID")
	apiToken := os.Getenv("WHATSAPP_TOKEN")

	url := fmt.Sprintf("https://api.ultramsg.com/%s/messages/chat", instanceID)

	payload := fmt.Sprintf("token=%s&to=%s&body=%s", apiToken, receiverID, replyMessage)
	
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	fmt.Println("Reply sent successfully!")
	return nil
}
