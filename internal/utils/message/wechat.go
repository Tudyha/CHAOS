package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	CORPID     = "wwfa6272be4a939e8b"
	CORPSECRET = "A0ijphI-2l-XNmuhZoQfrT0laS9XtQe28ACFcxXUhow"
	AGENTID    = "1000004"
	TOUSER     = "@all"
	TOKEN_URL  = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	MSG_URL    = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
	TOKEN_TTL  = 7200 // access_token 的有效期，单位秒
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type MessageResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

var (
	accessTokenMu sync.Mutex
	accessToken   string
	tokenExpiry   time.Time
)

func getAccessToken() (string, error) {
	accessTokenMu.Lock()
	defer accessTokenMu.Unlock()

	if time.Now().Before(tokenExpiry) {
		return accessToken, nil
	}

	url := fmt.Sprintf(TOKEN_URL, CORPID, CORPSECRET)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response AccessTokenResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	accessToken = response.AccessToken
	tokenExpiry = time.Now().Add(time.Duration(response.ExpiresIn-100) * time.Second) // 提前100秒刷新

	return accessToken, nil
}

func SendTextMessage(message string) error {
	accessToken, err := getAccessToken()
	if err != nil {
		return err
	}
	url := fmt.Sprintf(MSG_URL, accessToken)
	data := map[string]interface{}{
		"touser":  TOUSER,
		"msgtype": "text",
		"text":    map[string]string{"content": message},
		"agentid": AGENTID,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response MessageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	if response.ErrCode != 0 {
		return fmt.Errorf("error sending message: %s", response.ErrMsg)
	}

	return nil
}
