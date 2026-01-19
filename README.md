# GitHub Stats Inline Bot

A Telegram bot that provides GitHub profile and repository statistics via inline queries.

## Features

- **Profile Search**: Type `@bot username` (e.g., `@your_bot octocat`).
- **Repository Search**: Type `@bot user/repo` (e.g., `@your_bot google/go-github`).
- **Visuals**: Displays GitHub Streak cards for profiles and key metrics (stars, forks, etc.) for repositories.

## Tech Stack

- **Go** (Golang)
- [telego](https://github.com/mymmrac/telego) — Telegram Bot API library
- [go-github](https://github.com/google/go-github) — GitHub API client

## Setup

1. **Clone the repository**:
    ```bash
    git clone https://github.com/ashleybeloved/github-checker-bot.git
    cd github-checker-bot
    ```
2. **Configure environment variables**:
    ```bash
    TELEGRAM_TOKEN=your_telegram_bot_token
    GITHUB_TOKEN=your_github_personal_access_token
    ```
3. **Run the bot**:
    ```bash
    go run main.go
    ```
## Telegram Configuration

To enable inline mode:

1. Open @BotFather
2. Go to **Edit Bot > Inline Mode**
3. Click **Turn On**


