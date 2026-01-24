package inline

import (
	"context"
	"fmt"
	"html"
	"log"
	"math/rand"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/google/go-github/v60/github"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func Repository(ghClient *github.Client, ctx context.Context, query string) (*telego.InlineQueryResultArticle, error) {
	pathToRepo := strings.Split(query, "/")
	repo, _, err := ghClient.Repositories.Get(ctx, pathToRepo[0], pathToRepo[1])
	if err != nil {
		log.Printf("Ошибка GitHub API для %q: %v", query, err)
		return nil, err
	}

	safeQuery := html.EscapeString(query)
	safeDesc := html.EscapeString(repo.GetDescription())
	safeLang := html.EscapeString(repo.GetLanguage())

	repoInfo := fmt.Sprintf("<b>%s</b>\n\n", safeQuery)

	if safeDesc != "" {
		repoInfo += fmt.Sprintf("📝 <b>Description:</b> %s\n\n", safeDesc)
	}

	repoInfo += fmt.Sprintf(
		"⭐ <b>Stars:</b> %d\n"+
			"🍴 <b>Forks:</b> %d\n"+
			"👁️ <b>Subscribers:</b> %d\n\n",
		repo.GetStargazersCount(),
		repo.GetForksCount(),
		repo.GetSubscribersCount(),
	)

	if safeLang != "" {
		repoInfo += fmt.Sprintf("💻 <b>Top Language:</b> %s\n", safeLang)
	}

	repoInfo += fmt.Sprintf(
		"📅 <b>Created at:</b> %s\n\n"+
			`<a href="%s">Repository</a> | <a href="%s">Releases</a>`,
		humanize.Time(repo.GetCreatedAt().Time),
		html.EscapeString("https://github.com/"+query),
		html.EscapeString("https://github.com/"+query+"/releases"),
	)

	result := tu.ResultArticle(
		"repo_stats_"+query,
		"Repository: "+query,
		tu.TextMessage(
			repoInfo,
		).WithParseMode(telego.ModeHTML),
	)

	return result, nil
}

func Profile(ghClient *github.Client, ctx context.Context, query string) (*telego.InlineQueryResultArticle, error) {
	user, _, err := ghClient.Users.Get(ctx, query)
	if err != nil {
		log.Printf("Ошибка GitHub API для %q: %v", query, err)
		return nil, err
	}

	cardStatsURL := fmt.Sprintf("https://github-streak-generator-moranr123-production.up.railway.app/api/streak/card/%v?theme=58a6ff&fontSize=large&cardWidth=500&_t=%v", query, rand.Intn(100000))
	safeName := html.EscapeString(user.GetName())
	safeBio := html.EscapeString(strings.TrimSpace(user.GetBio()))
	safeLocation := html.EscapeString(user.GetLocation())

	profileInfo := fmt.Sprintf("<b>%s</b>\n\n", html.EscapeString(query))

	if safeName != "" {
		profileInfo += fmt.Sprintf("<b>👤 Name:</b> %s\n", safeName)
	}

	if safeBio != "" {
		profileInfo += fmt.Sprintf("<b>📝 Bio:</b> %s\n", safeBio)
	}

	if safeLocation != "" {
		profileInfo += fmt.Sprintf("<b>📍 Location:</b> %s\n\n", safeLocation)
	}

	profileInfo += fmt.Sprintf(
		"<b>👥 Followers:</b> %d\n<b>📦 Repositories:</b> %d\n<b>📅 Created at:</b> %v\n\n",
		user.GetFollowers(),
		user.GetPublicRepos(),
		humanize.Time(user.GetCreatedAt().Time),
	)

	githubProfileURL := "https://github.com/" + query
	githubReposURL := githubProfileURL + "?tab=repositories"

	profileInfo += fmt.Sprintf(
		`<a href="%s">Profile</a> | <a href="%s">Repositories</a>`,
		html.EscapeString(githubProfileURL),
		html.EscapeString(githubReposURL),
	)

	result := tu.ResultArticle(
		"stats_"+query,
		"Profile: "+query,
		tu.TextMessage(profileInfo).
			WithLinkPreviewOptions(&telego.LinkPreviewOptions{
				URL:           cardStatsURL,
				ShowAboveText: false,
			}).
			WithParseMode(telego.ModeHTML),
	)

	return result, nil
}
