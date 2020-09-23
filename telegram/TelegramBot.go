package telegram

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type CustomWriter struct {
	message []byte
}

func NewWrite() *CustomWriter {
	return &CustomWriter{}
}

func (w *CustomWriter) Write(p []byte) (int, error) {
	SendMessage(p)
	return len(p), nil
}

var (
	bot, _ = tgbotapi.NewBotAPI("1201939942:AAEOvtLTfqX6a1emGNCxSc8k0eupGjE67j8")
)

func SendMessage(p []byte){
	msg := tgbotapi.NewMessage(-1001170497596, string(p))
	bot.Send(msg)
}