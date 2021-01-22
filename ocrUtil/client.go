package ocrUtil

import (
	"github.com/astaxie/beego"
	gosdk "github.com/chenqinghe/baidu-ai-go-sdk"
	"github.com/chenqinghe/baidu-ai-go-sdk/vision/ocr"
)

type MyOCRClient struct {
	clt       *ocr.OCRClient
	AppID     string
	ApiKey    string
	SecretKey string
}

func NewMyOCRClient(appID, apiKey, secretKey string) *MyOCRClient {
	clt := &gosdk.Client{
		ClientID:     apiKey,
		ClientSecret: secretKey,
		Authorizer: OcrAuthorizer{
			AppID: appID,
		},
	}
	return &MyOCRClient{
		clt: &ocr.OCRClient{
			Client: clt,
		},
		AppID:     appID,
		ApiKey:    apiKey,
		SecretKey: secretKey,
	}
}

func NewDefaultOcrClient() *MyOCRClient {
	return NewMyOCRClient(
		beego.AppConfig.String(BAIDU_AI_APP_ID),
		beego.AppConfig.String(BAIDU_AI_API_KEY),
		beego.AppConfig.String(BAIDU_AI_SECRET_KEY))
}

func (c *MyOCRClient) Token() (string, error) {
	return GetAccessToken(c.AppID, c.ApiKey, c.SecretKey)
}
