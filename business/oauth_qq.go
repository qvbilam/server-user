package business

import (
	"encoding/json"
	"fmt"
	"user/cache"
	"user/enum"
	"user/utils"
)

type OAuthQQBusiness struct {
	AppId     string
	AppSecret string
	Uri       string
}

type QQTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type QQOpenIdResponses struct {
	ClientId string `json:"client_id"`
	OpenId   string `json:"openid"`
}

type QQUserResponse struct {
	Ret         int    `json:"ret"`
	Msg         string `json:"msg"`
	Nickname    string `json:"nickname"`
	GanderType  string `json:"ganderType"` // 1女 2男
	FigureurlQq string `json:"figureurl_qq"`
}

func (b *OAuthQQBusiness) getToken(code string) (*QQTokenResponse, error) {
	requestUrl := "https://graph.qq.com/oauth2.0/token"
	params := map[string]interface{}{}
	params["grant_type"] = "authorization_code"
	params["client_id"] = b.AppId
	params["client_secret"] = b.AppSecret
	params["code"] = code
	params["redirect_uri"] = b.Uri
	params["fmt"] = "json"

	res, err := utils.Get(requestUrl, params, nil)
	if err != nil {
		return nil, err
	}

	// access_token=5EEBE60F15AC2E27F58B5B5C1B1E456D&expires_in=7776000&refresh_token=D604CEEFB75C9037307F57040047E226
	fmt.Println(string(res))

	response := QQTokenResponse{}
	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (b *OAuthQQBusiness) getOpenId(token string) (*QQOpenIdResponses, error) {
	requestUrl := "https://graph.qq.com/oauth2.0/me"
	params := map[string]interface{}{
		"access_token": token,
		"fmt":          "json",
	}
	res, err := utils.Get(requestUrl, params, nil)
	if err != nil {
		return nil, err
	}
	// {"client_id":"101840991","openid":"FA800C280F8B1151C7214A7EC9ED3FE9"}
	//fmt.Println(string(res))

	response := QQOpenIdResponses{}
	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (b *OAuthQQBusiness) getUser(token, openId string) (*QQUserResponse, error) {
	requestUrl := "https://graph.qq.com/user/get_user_info"
	params := map[string]interface{}{
		"access_token":       token,
		"oauth_consumer_key": b.AppId,
		"openid":             openId,
		"fmt":                "json",
	}

	res, err := utils.Get(requestUrl, params, nil)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(res))

	response := QQUserResponse{}
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (b *OAuthQQBusiness) token(code string) (string, error) {
	server := cache.RedisServer{}
	token := server.GetOAuthQQToken(b.AppId)
	if token != "" {
		return token, nil
	}

	res, err := b.getToken(code)
	if err != nil {
		return "", err
	}
	return res.AccessToken, nil
}

func (b *OAuthQQBusiness) User(code string) *OAuthUserResponse {
	var token string
	var openId *QQOpenIdResponses
	var user *QQUserResponse
	var err error

	if token, err = b.token(code); err != nil {
		return nil
	}
	if openId, err = b.getOpenId(token); err != nil {
		return nil
	}
	if user, err = b.getUser(token, openId.OpenId); err != nil {
		return nil
	}

	gender := enum.GenderTypeMale
	if user.GanderType == "1" { // 2 = 男, 1=女
		gender = enum.GenderTypeFemale
	}

	return &OAuthUserResponse{
		Type:       enum.LoginMethodPlatformQQ,
		PlatformId: openId.OpenId,
		Nickname:   user.Nickname,
		Gender:     gender,
		Avatar:     user.FigureurlQq,
	}
}
