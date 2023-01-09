package main

import (
	"context"
	"github.com/joho/godotenv"
	gogpt "github.com/sashabaranov/go-gpt3"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButton("2"),
		tgbotapi.NewKeyboardButton("3"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("4"),
		tgbotapi.NewKeyboardButton("5"),
		tgbotapi.NewKeyboardButton("6"),
	),
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	//CREATE TG BOT CONNECTION
	log.Printf("Start to create connection to Telegram Bot")
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_KEY"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	//START CONSUME TG MESSAGES

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Text {
		case "open":
			msg.ReplyMarkup = numericKeyboard
		case "close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}

		//CREATE CHAT GPT CONNECTION
		log.Printf("Start to create connection to Chat GPT")

		gptClient := gogpt.NewClient(os.Getenv("OPENAI_API_KEY"))
		ctx := context.Background()

		req := gogpt.CompletionRequest{
			Model:     "text-davinci-003",
			MaxTokens: 1000,
			Prompt:    msg.Text,
		}
		resp, err := gptClient.CreateCompletion(ctx, req)
		if err != nil {
			return
		}
		answer := tgbotapi.NewMessage(update.Message.Chat.ID, resp.Choices[0].Text)

		if _, err := bot.Send(answer); err != nil {
			log.Panic(err)
		}
	}
}
