package ocrUtil

import (
	"testing"

	"github.com/astaxie/beego"
	"github.com/gyf841010/pz-infra-new/logging"
	//. "github.com/smartystreets/goconvey/convey"
)

func TestOCR(t *testing.T) {
	// 测试需要填入以下信息: redis链接url 秘钥,百度ai相关的3个配置 及图像路径
	logging.InitLogger("test")
	beego.AppConfig.Set("redisUrl", "")
	beego.AppConfig.Set("redisPass", "")

	beego.AppConfig.Set(BAIDU_AI_APP_ID, "")
	beego.AppConfig.Set(BAIDU_AI_API_KEY, "")
	beego.AppConfig.Set(BAIDU_AI_SECRET_KEY, "")

	// client := NewDefaultOcrClient()
	// Convey("身份证", t, func() {
	// 	img, err := client.GetImgFromFile("./id_card.jpg")
	// 	So(err, ShouldBeNil)
	// 	result, err := client.RecognizeIDCard(img)
	// 	So(err, ShouldBeNil)
	// 	fmt.Println(result)
	// })
}
