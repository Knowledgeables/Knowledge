package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Sender interface {
	Send(to, subject, html string) error
}

type ResendSender struct {
	apiKey   string
	fromAddr string
	client   *http.Client
}

func NewResendSender(apiKey, fromAddr string) *ResendSender {
	return &ResendSender{
		apiKey:   apiKey,
		fromAddr: fromAddr,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

type sendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

func (r *ResendSender) Send(to, subject, html string) error {
	if r.apiKey == "" {
		slog.Info("email (dev mode — no RESEND_API_KEY set)",
			"to", to,
			"subject", subject,
		)
		return nil
	}

	body, err := json.Marshal(sendEmailRequest{
		From:    r.fromAddr,
		To:      []string{to},
		Subject: subject,
		HTML:    html,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend: status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
