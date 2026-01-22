package main

import (
	"context"
	"dada/inline"
	"fmt"
	"html"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v60/github"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	ctx := context.Background()
	godotenv.Load(".env")

	botToken := os.Getenv("TELEGRAM_TOKEN")
	githubToken := os.Getenv("GITHUB_TOKEN")

	ghClient := github.NewClient(nil).WithAuthToken(githubToken)

	bot, err := telego.NewBot(botToken)
	if err != nil {
		log.Fatal(err)
	}

	updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

	for update := range updates {
		if update.Message != nil && update.Message.Text != "" {
			msg := tu.Message(
				update.Message.Chat.ChatID(),
				fmt.Sprintf("<b>Hello, %s!</b>\n\nI only work in inline mode, to call inline write:\n<code>@%s &lt;username&gt; or &lt;username/repo&gt;</code>",
					html.EscapeString(update.Message.From.FirstName),
					bot.Username(),
				),
			).WithParseMode(telego.ModeHTML)

			bot.SendMessage(ctx, msg)
		}

		if update.InlineQuery != nil {
			var result *telego.InlineQueryResultArticle
			iq := update.InlineQuery
			query := strings.TrimSpace(iq.Query)

			if query == "" {
				continue
			}

			if strings.Contains(query, "/") {
				result, err = inline.Repository(ghClient, ctx, query)
				if err != nil {
					continue
				}
			} else {
				result, err = inline.Profile(ghClient, ctx, query)
				if err != nil {
					continue
				}
			}

			err = bot.AnswerInlineQuery(ctx, &telego.AnswerInlineQueryParams{
				InlineQueryID: iq.ID,
				Results:       []telego.InlineQueryResult{result},
				CacheTime:     300,
			})

			if err != nil {
				notFoundResult := tu.ResultArticle(
					"notfound_"+query,
					"❌ User not found",
					tu.TextMessage(
						"User or repository not found",
					).WithParseMode(telego.ModeHTML),
				)
				bot.AnswerInlineQuery(ctx, &telego.AnswerInlineQueryParams{
					InlineQueryID: iq.ID,
					Results:       []telego.InlineQueryResult{notFoundResult},
					CacheTime:     1,
				})
				continue
			}
		}
	}
}
