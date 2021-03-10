package ocrUtil

import (
	"io"
	"net/http"
	"time"

	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"errors"

	"github.com/chenqinghe/baidu-ai-go-sdk/vision"
	"github.com/chenqinghe/baidu-ai-go-sdk/vision/ocr"
	"github.com/go-resty/resty/v2"
	. "github.com/gyf841010/pz-infra-new/logging"
	timeutil "github.com/gyf841010/pz-infra-new/timeUtil"
)

const (
	BAIDU_AI_TIME_FORMAT_OCR = "2006/01/02"
	OCR_VIN_URL              = "https://aip.baidubce.com/rest/2.0/ocr/v1/vin_code?access_token=%s"
	OCR_CAR_TYPE_URL         = "https://aip.baidubce.com/rest/2.0/image-classify/v1/car?access_token=%s"
	OCR_NUMBER_URL           = "https://aip.baidubce.com/rest/2.0/ocr/v1/numbers?access_token=%s"
	OCR_BUSINESS_LICENSE_URL = "https://aip.baidubce.com/rest/2.0/ocr/v1/business_license?access_token=%s"
)

const (
	//ocr 错误
	ERROR_CODE_OCR_HTTP_TOKEN = "通过http请求百度ai token失败"
	ERROR_CODE_OCR_IMG_EMPTY  = "OCR,图像为空,无法识别"
	ERROR_CODE_DOWNLOAD_IMG   = "通过url下载图像数据失败,无法识别"
)

func (c *MyOCRClient) GetImgFromFile(file string) (*vision.Image, error) {
	img, err := vision.FromFile(file)
	return img, err
}

func (c *MyOCRClient) GetImgFromUrl(url string) (*vision.Image, error) {
	img, err := vision.FromUrl(url)
	return img, err
}

func (c *MyOCRClient) GetImgFromBytes(bts []byte) (*vision.Image, error) {
	img, err := vision.FromBytes(bts)
	return img, err
}

// 身份证识别(正/反面),接口可自动识别图像是正面还是反面,
func (c *MyOCRClient) RecognizeIDCard(img *vision.Image) (*IDCardRecognitionResponse, error) {
	resp, err := c.clt.IdCardRecognize(
		img,
		ocr.DetectDirection(),
	)
	if err != nil {
		return nil, err
	}
	resIdCard := &IdCardRes{}
	err = resp.ToJSON(&resIdCard)

	Log.Info(fmt.Sprintln(resIdCard))

	return TransHttpResToIDCardRes(resIdCard), err
}

// 银行卡识别
func (c *MyOCRClient) RecognizeBankCard(img *vision.Image) (*BankCardRes, error) {
	resp, err := c.clt.BankcardRecognize(
		img,
		ocr.DetectDirection())
	if err != nil {
		return nil, err
	}
	bankCard := &BankCardRes{}
	err = resp.ToJSON(bankCard)
	return bankCard, err
}

// 驾驶证
func (c *MyOCRClient) RecognizeDrivingLicense(img *vision.Image) (*DrivingLicenseRes, error) {
	resp, err := c.clt.DriverLicenseRecognize(
		img,
		ocr.DetectDirection(),
	)
	if err != nil {
		return nil, err
	}
	drivingLicense := &DrivingLicenseRes{}
	err = resp.ToJSON(drivingLicense)
	return drivingLicense, err
}

// 行驶证 side: front/back - front：默认值，识别行驶证主页 - back：识别行驶证副页
func (c *MyOCRClient) RecognizeVehicleLicense(img *vision.Image, side string) (*VehicleLicenseRes, error) {
	resp, err := c.clt.VehicleLicenseRecognize(
		img,
		ocr.DetectDirection(),
		vehicleParam(side),
	)
	if err != nil {
		return nil, err
	}
	vehicleLicense := &VehicleLicenseRes{}
	err = resp.ToJSON(vehicleLicense)
	return vehicleLicense, err
}

//行驶证 side: front/back - front：默认值，识别行驶证主页 - back：识别行驶证副页
func vehicleParam(side string) ocr.RequestParam {
	if side != "front" && side != "back" {
		side = "front"
	}
	return func(m map[string]interface{}) {
		m["vehicle_license_side"] = side
	}
}

// 车牌识别
func (c *MyOCRClient) RecognizePlateNumber(img *vision.Image) (*PlateNumberRes, error) {
	resp, err := c.clt.LicensePlateRecognize(
		img,
		ocr.DetectDirection(),
	)
	if err != nil {
		return nil, err
	}
	plateNum := &PlateNumberRes{}
	err = resp.ToJSON(plateNum)
	return plateNum, err
}

// Vin码识别
func (c *MyOCRClient) RecognizeVin(img *vision.Image) (*VinRes, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(OCR_VIN_URL, token)

	body := make(map[string]string, 0)
	if img.Reader == nil {
		if img.Url == "" {
			return nil, errors.New(ERROR_CODE_OCR_IMG_EMPTY + "Vin码识别")
		} else {
			body["url"] = img.Url
		}
	} else {
		base64Str, err := img.Base64Encode()
		if err != nil {
			return nil, err
		}
		body["image"] = base64Str
	}

	// application/x-www-form-urlencoded格式
	// 请求数据通过SetFormData 而不是 SetBody()
	// https://github.com/go-resty/resty/blob/b90c855d687abc2b9efef0e936ea5462836954a1/client.go
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(body).
		Post(url)
	if err != nil {
		return nil, err
	}
	vin := &VinRes{}
	err = json.NewDecoder(bytes.NewReader(resp.Body())).Decode(vin)
	return vin, err
}

func (c *MyOCRClient) RecognizeVinV2(img *vision.Image) (*VinRes, error) {
	resp, err := c.clt.VinRecognize(
		img,
	)
	if err != nil {
		return nil, err
	}
	vin := &VinRes{}
	err = resp.ToJSON(vin)
	return vin, err
}

// 验证码,数字识别 对图片中的数字进行提取和识别，自动过滤非数字内容，仅返回数字内容及其位置信息
func (c *MyOCRClient) RecognizeNumber(img *vision.Image) (*NumberRes, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(OCR_NUMBER_URL, token)

	body := make(map[string]string, 0)
	if img.Reader == nil {
		if img.Url == "" {
			return nil, errors.New(ERROR_CODE_OCR_IMG_EMPTY + "数字/验证码识别")
		} else {
			body["url"] = img.Url
		}
	} else {
		base64Str, err := img.Base64Encode()
		if err != nil {
			return nil, err
		}
		body["image"] = base64Str
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(body).
		Post(url)
	if err != nil {
		return nil, err
	}
	number := &NumberRes{}
	err = json.NewDecoder(bytes.NewReader(resp.Body())).Decode(number)
	return number, err
}

func (c *MyOCRClient) RecognizeNumberV2(img *vision.Image) (*NumberRes, error) {
	resp, err := c.clt.NumberRecognize(
		img,
	)
	if err != nil {
		return nil, err
	}
	number := &NumberRes{}
	err = resp.ToJSON(number)
	return number, err
}

// 车型识别
// 百度车型识别接口不支持通过Url识别,需要将图像下载到本地转换为图像
// 一张图像接口推测得到的结果有多种,需要返回多少种
// topNum,返回车型数量,不设置默认为1,传入设置使用第一个参数作为车型数量
func (c *MyOCRClient) RecognizeCarType(img *vision.Image, topNum ...int) (*CarTypeRes, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(OCR_CAR_TYPE_URL, token)

	body := make(map[string]string, 0)
	if img.Reader == nil {
		if img.Url == "" {
			return nil, errors.New(ERROR_CODE_OCR_IMG_EMPTY + "车型识别")
		} else {
			file := "./cat_type.png"
			err := downloadFile(img.Url, file)
			if err != nil {
				return nil, errors.New(ERROR_CODE_DOWNLOAD_IMG + "车型识别")
			}
			defer func() {
				os.Remove(file)
			}()
			img, err = vision.FromFile(file)
			if err != nil {
				return nil, err
			}
		}
	}

	base64Str, err := img.Base64Encode()
	if err != nil {
		return nil, err
	}
	body["image"] = base64Str
	if len(topNum) <= 0 {
		body["top_num"] = fmt.Sprintln(1)
	} else {
		body["top_num"] = fmt.Sprintln(topNum[0])
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(body).
		Post(url)
	if err != nil {
		return nil, err
	}
	carType := &CarTypeRes{}
	err = json.NewDecoder(bytes.NewReader(resp.Body())).Decode(carType)
	fmt.Println("resp", string(resp.Body()))
	return carType, err
}

func (c *MyOCRClient) RecognizeCarTypeV2(img *vision.Image, topNum ...int) (*CarTypeRes, error) {
	if img.Reader == nil && img.Url != "" {
		file := "./cat_type.png"
		err := downloadFile(img.Url, file)
		if err != nil {
			return nil, errors.New(ERROR_CODE_DOWNLOAD_IMG + "车型识别")
		}
		defer func() {
			os.Remove(file)
		}()
		img, err = vision.FromFile(file)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.clt.CarTypeRecognize(
		img,
		ocr.CarTypeTopNum(1),
		ocr.CarTypeBaikeNum(0),
	)
	if err != nil {
		return nil, err
	}
	carType := &CarTypeRes{}
	err = resp.ToJSON(carType)
	return carType, err
}

// 营业执照
func (c *MyOCRClient) RecognizeBusinessLicense(img *vision.Image) (*BusinessLicenseRes, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(OCR_BUSINESS_LICENSE_URL, token)

	body := make(map[string]string, 0)
	if img.Reader == nil {
		if img.Url == "" {
			return nil, errors.New(ERROR_CODE_OCR_IMG_EMPTY + "营业执照识别")
		} else {
			body["url"] = img.Url
		}
	} else {
		base64Str, err := img.Base64Encode()
		if err != nil {
			return nil, err
		}
		body["image"] = base64Str
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(body).
		Post(url)
	if err != nil {
		return nil, err
	}
	businessLicense := &BusinessLicenseRes{}
	err = json.NewDecoder(bytes.NewReader(resp.Body())).Decode(businessLicense)
	return businessLicense, err
}

func (c *MyOCRClient) RecognizeBusinessLicenseV2(img *vision.Image) (*BusinessLicenseRes, error) {
	resp, err := c.clt.BusinessLicenseRecognize(
		img,
	)
	if err != nil {
		return nil, err
	}
	business := &BusinessLicenseRes{}
	err = resp.ToJSON(business)
	return business, err
}

// 从身份证号获取生日
func GetBirthdayFromIDCard(idCardNum string) (int64, error) {
	if len(idCardNum) != 18 {
		return 0, errors.New("wrong id_card value")
	}
	nYear, err := strconv.Atoi(string(idCardNum[6:10]))
	if err != nil {
		return 0, err
	}
	nMonth, err := strconv.Atoi(string(idCardNum[10:12]))
	if err != nil {
		return 0, err
	}
	nDay, err := strconv.Atoi(string(idCardNum[12:14]))
	if err != nil {
		return 0, err
	}

	birthday := time.Date(nYear, time.Month(nMonth), nDay, 0, 0, 0, 0, time.Local)
	return timeutil.ConvertTimeToUnixBigInt(birthday), nil
}

const (
	//性别 0:未知 1：男 2：女
	SEX_UNKNOWN = iota
	SEX_MAN
	SEX_WOMEN
)

// 由身份证获取性别
func GetGenderFromIDCard(idCardNum string) int {
	if len(idCardNum) != 18 {
		return SEX_UNKNOWN
	}
	val, err := strconv.Atoi(string(idCardNum[16]))
	if err != nil {
		return SEX_UNKNOWN
	}
	gender := val % 2
	//偶数为女 奇数为男
	if gender == 0 {
		return SEX_WOMEN
	} else if gender == 1 {
		return SEX_MAN
	}
	return SEX_UNKNOWN
}

func GetGenderFromString(gender string) int {
	if gender == "男" {
		return SEX_MAN
	} else if gender == "女" {
		return SEX_WOMEN
	}
	return SEX_UNKNOWN
}

// 百度AI识别出来的时间字符串转时间戳
func TransStrToTime(strTime string) int64 {
	if len(strTime) < 8 {
		return 0
	}
	formatTime := strTime[0:4] + "/" + strTime[4:6] + "/" + strTime[6:]
	res, err := time.Parse(BAIDU_AI_TIME_FORMAT_OCR, formatTime)
	if err != nil {
		return 0
	}
	return timeutil.ConvertTimeToUnixBigInt(res)
}

func TransHttpResToIDCardRes(in *IdCardRes) *IDCardRecognitionResponse {
	resp := &IDCardRecognitionResponse{}
	resp.Name = in.WordResults.Name.Words
	resp.Nation = in.WordResults.Nation.Words
	resp.Address = in.WordResults.Address.Words
	resp.Sex = in.WordResults.Sex.Words
	resp.BirthDay = TransStrToTime(in.WordResults.BirthDay.Words)
	resp.IDCradNumber = in.WordResults.IDCradNumber.Words
	resp.TimeOutDate = TransStrToTime(in.WordResults.TimeOutData.Words)
	resp.Organization = in.WordResults.Organization.Words
	resp.IssueDate = TransStrToTime(in.WordResults.IssueDate.Words)

	return resp
}

func downloadFile(url, localPath string) error {
	out, err := os.Create(localPath)
	if err != nil {
		Log.Error("Failed to Create File", With("localPath", localPath), WithError(err))
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		Log.Error("Failed to Download From", With("url", url), WithError(err))
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		Log.Error("Failed to Copy Body For File", WithError(err))
		return err
	}
	return nil
}
