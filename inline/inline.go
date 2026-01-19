package inline

import (
	"context"
	"fmt"
	"log"
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

	repoInfo := fmt.Sprintf("*%v*\n\n", query)

	if desc := repo.GetDescription(); desc != "" {
		repoInfo = fmt.Sprintf(repoInfo+"📝 *Description:* %s\n\n", desc)
	}

	repoInfo = fmt.Sprintf(repoInfo+"⭐ *Stars:* %d\n🍴 *Forks:* %d\n👁️ *Subscribers:* %d\n\n",
		repo.GetStargazersCount(),
		repo.GetForksCount(),
		repo.GetSubscribersCount())

	if lang := repo.GetLanguage(); lang != "" {
		repoInfo = fmt.Sprintf(repoInfo+"💻 *Top Language:* %v\n", lang)
	}

	repoInfo = fmt.Sprintf(repoInfo+"📅 Created at: %s\n\n[Repository](%s) | [Releases](%s)",
		humanize.Time(repo.GetCreatedAt().Time),
		"https://github.com/"+query,
		"https://github.com/"+query+"/releases")

	result := tu.ResultArticle(
		"repo_stats_"+query,
		"Repository: "+query,
		tu.TextMessage(
			repoInfo,
		).WithParseMode(telego.ModeMarkdown),
	)

	return result, nil
}

func Profile(ghClient *github.Client, ctx context.Context, query string) (*telego.InlineQueryResultArticle, error) {
	user, _, err := ghClient.Users.Get(ctx, query)
	if err != nil {
		log.Printf("Ошибка GitHub API для %q: %v", query, err)
		return nil, err
	}

	cardStatsURL := "https://github-streak-generator-moranr123-production.up.railway.app/api/streak/card/" + query + "?theme=58a6ff&fontSize=large&cardWidth=500&_t=1768842469398"
	profileInfo := fmt.Sprintf("*%s*\n\n", query)

	if name := user.GetName(); name != "" {
		profileInfo = fmt.Sprintf(profileInfo+"*👤 Name: *%s\n", name)
	}

	if bio := user.GetBio(); bio != "" {
		profileInfo = fmt.Sprintf(profileInfo+"*📝 Bio: *%s\n", strings.TrimSpace(bio))
	}

	if location := user.GetLocation(); location != "" {
		profileInfo = fmt.Sprintf(profileInfo+"*📍 Location: *%s\n\n", location)
	}

	profileInfo = fmt.Sprintf(profileInfo+"*👥 Followers:* %d\n*📦 Repositories: *%d\n*📅 Created at: *%v\n\n",
		user.GetFollowers(),
		user.GetPublicRepos(),
		humanize.Time(user.GetCreatedAt().Time),
	)

	profileInfo = fmt.Sprintf(profileInfo+"[Profile](%v) | [Repositories](%v)", "https://github.com/"+query, "https://github.com/"+query+"?tab=repositories")

	result := tu.ResultArticle(
		"stats_"+query,
		"Profile: "+query,
		tu.TextMessage(
			profileInfo,
		).WithLinkPreviewOptions(&telego.LinkPreviewOptions{URL: cardStatsURL, ShowAboveText: false}).WithParseMode(telego.ModeMarkdown),
	)

	return result, nil
}
