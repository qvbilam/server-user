package business

type OAuthUserResponse struct {
	Type       string `json:"type"`
	PlatformId string `json:"platform_id"`
	Nickname   string `json:"nickname"`
	Gender     string `json:"gender"`
	Avatar     string `json:"avatar"`
}
