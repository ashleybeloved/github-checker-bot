package main

import (
	"context"
	"fmt"
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

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatal(err)
	}

	updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

	for update := range updates {
		if update.InlineQuery != nil {
			var result *telego.InlineQueryResultArticle
			iq := update.InlineQuery
			query := strings.TrimSpace(iq.Query)

			if query == "" {
				continue
			}

			if strings.Contains(query, "/") {
				result, err = Repository(ghClient, ctx, query)
				if err != nil {
					continue
				}
			} else {
				result, err = Profile(ghClient, ctx, query)
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
					"❌ Пользователь не найден",
					tu.TextMessage(
						"Пользователь с таким именем не существует на GitHub.",
					),
				)
				bot.AnswerInlineQuery(ctx, &telego.AnswerInlineQueryParams{
					InlineQueryID: iq.ID,
					Results:       []telego.InlineQueryResult{notFoundResult},
					CacheTime:     300,
				})
				continue
			}
		}
	}
}

func Repository(ghClient *github.Client, ctx context.Context, query string) (*telego.InlineQueryResultArticle, error) {
	pathToRepo := strings.Split(query, "/")
	repo, _, err := ghClient.Repositories.Get(ctx, pathToRepo[0], pathToRepo[1])
	if err != nil {
		log.Printf("Ошибка GitHub API для %q: %v", query, err)
		return nil, err
	}

	cardStatsURL := fmt.Sprintf("https://pixel-profile.vercel.app/api/github-stats?username=%v", query)

	profileInfo := fmt.Sprintf(
		"%v", repo.GetFullName(),
	)

	result := tu.ResultArticle(
		"repo_stats_"+query,
		"Статистика для репозитория: "+query,
		tu.TextMessage(
			profileInfo,
		).WithLinkPreviewOptions(&telego.LinkPreviewOptions{URL: cardStatsURL}).WithParseMode(telego.ModeMarkdown),
	)

	return result, nil
}

func Profile(ghClient *github.Client, ctx context.Context, query string) (*telego.InlineQueryResultArticle, error) {
	user, _, err := ghClient.Users.Get(ctx, query)
	if err != nil {
		log.Printf("Ошибка GitHub API для %q: %v", query, err)
		return nil, err
	}

	cardStatsURL := fmt.Sprintf("https://pixel-profile.vercel.app/api/github-stats?username=%v", query)

	profileInfo := fmt.Sprintf(
		"%v", user.GetName(),
	)

	result := tu.ResultArticle(
		"stats_"+query,
		"Статистика для пользователя: "+query,
		tu.TextMessage(
			profileInfo,
		).WithLinkPreviewOptions(&telego.LinkPreviewOptions{URL: cardStatsURL}).WithParseMode(telego.ModeMarkdown),
	)

	return result, nil
}
