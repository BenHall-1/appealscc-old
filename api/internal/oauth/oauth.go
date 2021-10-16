package oauth

import (
	"os"

	discord "github.com/ravener/discord-oauth2"

	"golang.org/x/oauth2"
)

func DiscordOAuth() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  os.Getenv("DISCORD_REDIRECT_URL"),
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		Scopes:       []string{discord.ScopeEmail},
		Endpoint:     discord.Endpoint,
	}
}
