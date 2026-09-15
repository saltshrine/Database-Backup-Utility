package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SlackNotifier struct {
	WebhookURL string
}

type Payload struct {
	Text string `json:"text"`
}

func (s *SlackNotifier) Send(message string) error {
	if s.WebhookURL == "" {
		return nil
	}

	data, _ := json.Marshal(Payload{Text: message})
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(s.WebhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("slack webhook error status: %d", resp.StatusCode)
	}

	return nil
}
