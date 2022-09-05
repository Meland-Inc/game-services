package util

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func makeHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 不检测服务器数据证书
			},
		},
	}
}

// 不带授权get请求 method= "GET" || "POST"
func HttpRequest(urlstr string, method string, reqBody string, token string) ([]byte, error) {
	request, e := http.NewRequest(method, urlstr, strings.NewReader(reqBody))
	if e != nil {
		return nil, e
	}
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Authorization", token)

	client := makeHttpClient()
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Error response.StatusCode=%v urlstr=%s", response.StatusCode, urlstr)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

// 带授权的http请求，支持传入方式 method= "GET" || "POST"
func HttpRequestAuth(url, method string, requestBody string, auth string) ([]byte, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))

	client := makeHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error response.StatusCode=%v urlstr=%s", resp.StatusCode, url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}
