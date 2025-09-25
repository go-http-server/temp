package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Bot interface {
	SendMessage(message string) error
}

type BotTelegram struct {
	BotToken string
	ChatID   string
}

type PayloadBotTelegram struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func NewBotTelegramService(token, chatID string) Bot {
	return &BotTelegram{
		BotToken: token,
		ChatID:   chatID,
	}
}

func (bot BotTelegram) SendMessage(message string) error {
	payload := PayloadBotTelegram{
		ChatID: bot.ChatID,
		Text:   message,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("cannot marshal payload: %s", err)
	}

	reader := bytes.NewReader(data)
	urlSendNotify := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", bot.BotToken)

	_, err = http.Post(urlSendNotify, "application/json", reader)
	if err != nil {
		return fmt.Errorf("cannot request to bot tele: %s", err)
	}

	return nil
}
