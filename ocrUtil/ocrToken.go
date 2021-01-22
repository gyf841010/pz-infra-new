package ocrUtil

import (
	"bytes"
	"encoding/json"
	"fmt"

	"errors"

	gosdk "github.com/chenqinghe/baidu-ai-go-sdk"
	"github.com/go-resty/resty/v2"
	. "github.com/gyf841010/pz-infra-new/logging"
	"github.com/gyf841010/pz-infra-new/redisUtil"
)

const (
	BAIDU_AI_ACCESS_TOKEN_API            = "https://aip.baidubce.com/oauth/2.0/token"
	BAIDU_AI_ACCESS_TOKEN_KEY            = "INFRA:baidu_ai_access_token:%s"
	BAIDU_AI_ACCESS_TOKEN_KEY_VAILD_TIME = "baidu_ai_access_token_vaild_time"
	// token 过期时间
	BAIDU_AI_TOKEN_EXPIRE = 3600 * 24 * 3

	//百度图像识别-字符识别
	BAIDU_AI_APP_ID     = "BaiduAI.AppID"
	BAIDU_AI_API_KEY    = "BaiduAI.ApiKey"
	BAIDU_AI_SECRET_KEY = "BaiduAI.SecretKey"
)

// 自定义认证方法:增加通过redis获取access_token的方法
type OcrAuthorizer struct {
	AppID string
}

var _ gosdk.Authorizer = OcrAuthorizer{}

func (ocr OcrAuthorizer) Authorize(client *gosdk.Client) error {
	accessToken, err := GetAccessToken(ocr.AppID, client.ClientID, client.ClientSecret)
	client.AccessToken = accessToken
	return err
}

func GetAccessToken(appID, apiKey, secretKey string) (accessToken string, err error) {
	accessToken, err = GetAccessTokenFromRedis(appID)
	if err != nil || len(accessToken) <= 0 {
		Log.ErrorWithStack("从redis获取百度ai token失败", err)
	} else {
		return accessToken, nil
	}

	Log.Info("从Redis获取token失败,尝试通过http获取")
	tokenResp, err := HttpAccessToken(apiKey, secretKey)
	if err != nil {
		return "", err
	}
	accessToken = tokenResp.AccessToken
	err = AddAccessTokenToRedis(appID, tokenResp)
	if err != nil {
		Log.Info("fail to AddAccessTokenToRedis" + fmt.Sprintf("%+v", err))
	}
	return tokenResp.AccessToken, nil
}

func GetBaiduAiTokenKey(appID string) string {
	return fmt.Sprintf(BAIDU_AI_ACCESS_TOKEN_KEY, appID)
}

func GetAccessTokenFromRedis(appID string) (string, error) {
	key := GetBaiduAiTokenKey(appID)
	accessToken, err := redisUtil.GetString(key)
	return accessToken, err
}

func HttpAccessToken(apiKey, secretKey string) (accessToken *AuthResponse, err error) {
	url := BAIDU_AI_ACCESS_TOKEN_API
	accessTokenResponse := &AuthResponse{}
	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     apiKey,
			"client_secret": secretKey,
		}).
		Post(url)
	err = json.NewDecoder(bytes.NewReader(resp.Body())).Decode(accessTokenResponse)

	if accessTokenResponse.ERROR != "" || accessTokenResponse.AccessToken == "" {
		return nil, errors.New("通过http获取百度ai access_token失败")
	}
	return accessTokenResponse, err
}

// 设置token到redis
func AddAccessTokenToRedis(appID string, accessToken *AuthResponse) error {
	key := GetBaiduAiTokenKey(appID)
	err := redisUtil.SetStringWithExpire(key, accessToken.AccessToken, BAIDU_AI_TOKEN_EXPIRE)
	return err
}
