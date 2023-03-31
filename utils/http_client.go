package utils

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

func params(d map[string]interface{}) *url.Values {
	p := url.Values{}
	for k, v := range d {
		p.Set(k, v.(string))
		p.Set(k, v.(string))
	}
	return &p
}

//func Get(u string, d map[string]interface{}) ([]byte, error) {
//	p := params(d)
//	urlWithParams := u + "?" + p.Encode()
//	// 创建一个GET请求
//	resp, err := http.Get(urlWithParams)
//	defer resp.Body.Close()
//	if err != nil {
//		return nil, err
//	}
//
//	// 读取响应内容
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	return body, nil
//}

func Get(url string, data map[string]interface{}, headers map[string]interface{}) ([]byte, error) {
	p := params(data)
	urlWithParams := url + "?" + p.Encode()
	return sendGet(urlWithParams, nil, headers)
}

func Post(url string, data map[string]interface{}, headers map[string]interface{}) ([]byte, error) {
	if headers == nil {
		headers = make(map[string]interface{})
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return sendPost(url, data, headers)
}

func PostJson(url string, data map[string]interface{}, headers map[string]interface{}) ([]byte, error) {
	if headers == nil {
		headers = make(map[string]interface{})
	}
	headers["Content-Type"] = "application/json"
	return sendPost(url, data, headers)
}

func sendPost(url string, data map[string]interface{}, headers map[string]interface{}) ([]byte, error) {
	return send("POST", url, data, headers)
}

func sendGet(url string, data map[string]interface{}, headers map[string]interface{}) ([]byte, error) {
	return send("GET", url, data, headers)
}

func send(method, url string, data map[string]interface{}, headers map[string]interface{}) ([]byte, error) {
	p := params(data)

	client := &http.Client{}
	request, err := http.NewRequest(method, url, strings.NewReader(p.Encode()))
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v.(string))
		}
	}
	resp, err := client.Do(request)
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
