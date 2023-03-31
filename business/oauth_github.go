package business

import (
	"encoding/json"
	"fmt"
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
		return "", nil
	}

	response := GithubTokenResponse{}

	// {"access_token":"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","token_type":"bearer","scope":""}
	//fmt.Println(string(res))

	if err := json.Unmarshal(res, &response); err != nil {
		return "", err
	}
	return response.AccessToken, nil
}

func (b *OAuthGitHubBusiness) getUser(token string) {
	requestUrl := "https://api.github.com/user"
	headers := map[string]interface{}{}
	headers["Authorization"] = "token " + token
	headers["Accept"] = "application/json"
	res, err := utils.Post(requestUrl, nil, headers)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(res))
}

func (b *OAuthGitHubBusiness) token(code string) (string, error) {
	return b.getToken(code)
}

func (b *OAuthGitHubBusiness) User(code string) *OAuthUserResponse {
	token, _ := b.token(code)
	b.getUser(token)
	return nil
}
