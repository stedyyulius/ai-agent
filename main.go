package main

import (
	"ai-agent/integrations"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = integrations.ConnectDeepSeek()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("LIVE"))
	})

	http.HandleFunc("/webhook", integrations.ListenToWhatsapp)

	port := "8080"
	fmt.Printf("Listening for WhatsApp messages on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}