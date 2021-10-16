package discordmodel

type DiscordUser struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	Public_flags  int    `json:"public_flags"`
	Flags         int    `json:"flags"`
	Banner        string `json:"banner"`
	Banner_color  string `json:"banner_color"`
	Accent_color  string `json:"accent_color"`
	Locale        string `json:"locale"`
	Mfa_enabled   bool   `json:"mfa_enabled"`
	Premium_type  int    `json:"premium_type"`
	Email         string `json:"email"`
	Verified      bool   `json:"verified"`
}
