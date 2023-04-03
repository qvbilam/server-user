package business

import (
	"encoding/json"
	"user/enum"
	"user/utils"
)

type OAuthGitHubBusiness struct {
	AppId      string `json:"app_id"`
	AppSecrete string `json:"app_secrete"`
	Uri        string `json:"uri"`
}

type GithubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Bearer      string `json:"bearer"`
}

type GithubUserResponse struct {
	Name      string `json:"name"`       // 昵称
	Login     string `json:"login"`      // 账号
	AvatarUrl string `json:"avatar_url"` // 头像
}

func (b *OAuthGitHubBusiness) getToken(code string) (string, error) {
	requestUrl := "https://github.com/login/oauth/access_token"
	params := map[string]interface{}{}
	params["client_id"] = b.AppId
	params["client_secret"] = b.AppSecrete
	params["code"] = code

	headers := map[string]interface{}{}
	headers["Accept"] = "application/json"

	res, err := utils.Post(requestUrl, params, headers)
	if err != nil {
		return "", err
	}

	response := GithubTokenResponse{}

	if err := json.Unmarshal(res, &response); err != nil {
		return "", err
	}
	return response.AccessToken, nil
}

func (b *OAuthGitHubBusiness) getUser(token string) (*GithubUserResponse, error) {
	requestUrl := "https://api.github.com/user"
	headers := map[string]interface{}{}
	headers["Authorization"] = "Bearer " + token
	headers["Accept"] = "application/json"
	res, err := utils.Get(requestUrl, nil, headers)
	if err != nil {
		return nil, err
	}

	response := GithubUserResponse{}
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (b *OAuthGitHubBusiness) token(code string) (string, error) {
	return b.getToken(code)
}

func (b *OAuthGitHubBusiness) User(code string) *OAuthUserResponse {
	var user *GithubUserResponse
	var token string
	var err error
	if token, err = b.token(code); err != nil {
		return nil
	}
	if user, err = b.getUser(token); err != nil {
		return nil
	}

	return &OAuthUserResponse{
		Type:       enum.LoginMethodPlatformWGitHub,
		PlatformId: user.Login,
		Nickname:   user.Name,
		Gender:     enum.GenderTypeMale,
		Avatar:     user.AvatarUrl,
	}
}
