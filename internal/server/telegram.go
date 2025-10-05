package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

func (s *Server) sendTelegramMessage(text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.cfg.TelegramBotToken)
	body := fmt.Sprintf(`{"chat_id": %d, "text": "%s"}`, s.cfg.TelegramChatID, text)

	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API error: %s", string(respBody))
	}
	return nil
}
