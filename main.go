package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	botToken := os.Getenv("TELEGRAM_TOKEN")
	_ = os.Getenv("GITHUB_TOKEN")

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

	for update := range updates {
		if update.InlineQuery != nil {
			iq := update.InlineQuery
			name := iq.From.FirstName

			bot.AnswerInlineQuery(ctx, tu.InlineQuery(
				iq.ID,

				tu.ResultArticle(
					"GitHub Profile Stats",
					"Hello",
					tu.TextMessage(
						fmt.Sprintf("Hello %s\n\nYour query:\n```%#+v```", name, iq),
					).WithParseMode(telego.ModeMarkdownV2),
				).WithDescription(fmt.Sprintf("Query: %q", iq.Query))))
		}
	}
}
