package ocrUtil

type AuthResponse struct {
	AccessToken      string `json:"access_token"`  //要获取的Access Token
	ExpireIn         int64  `json:"expires_in"`    //Access Token的有效期(秒为单位，一般为1个月)；
	RefreshToken     string `json:"refresh_token"` //以下参数忽略，暂时不用
	Scope            string `json:"scope"`
	SessionKey       string `json:"session_key"`
	SessionSecret    string `json:"session_secret"`
	ERROR            string `json:"error"`             //错误码；关于错误码的详细信息请参考鉴权认证错误码(http://ai.baidu.com/docs#/Auth/top)
	ErrorDescription string `json:"error_description"` //错误描述信息，帮助理解和解决发生的错误。
}

type IDCardRecognitionRequest struct {
	IdCardSide int `json:"id_card_side" description:"图像是正面还是反面，1：正面（含照片的一面） 2：反面（带国徽的一面）"`
}

type IdCardRes struct {
	Direction        int         `json:"direction" description:"图像方向，当图像旋转时，返回该参数。- 1:未定义，- 0:正向，- 1: 逆时针90度，- 2:逆时针180度，- 3:逆时针270度"`
	ImageStatus      string      `json:"image_status" description:" normal-识别正常 reversed_side-身份证正反面颠倒 non_idcard-上传的图片中不包含身份证 blurred-身份证模糊 other_type_card-其他类型证照 over_exposure-身份证关键字段反光或过曝 over_dark-身份证欠曝（亮度过低）unknown-未知状态"`
	RiskType         string      `json:"risk_type" description:"输入参数 detect_risk = true 时，则返回该字段识别身份证类型: normal-正常身份证；copy-复印件；temporary-临时身份证；screen-翻拍；unknown-其他未知情况"`
	EditTool         string      `json:"edit_tool" description:"如果参数 detect_risk = true 时，则返回此字段。如果检测身份证被编辑过，该字段指定编辑软件名称，如:Adobe Photoshop CC 2014 (Macintosh),如果没有被编辑过则返回值无此参数"`
	LogID            int64       `json:"log_id" description:"唯一的log id，用于问题定位"`
	Photo            string      `json:"photo" description:"当请求参数 detect_photo = true时返回，头像切图的 base64 编码（无编码头，需自行处理）"`
	PhotoLocation    LocationInt `json:"photo_location" description:"当请求参数 detect_photo = true时返回，头像的位置信息（坐标0点为左上角）"`
	IdCardNumberType int         `json:"idcard_number_type" description:"用于校验身份证号码、性别、出生是否一致，输出结果及其对应关系如下 -1: 身份证正面所有字段全为空 0: 身份证证号识别错误 1: 身份证证号和性别、出生信息一致 2: 身份证证号和性别、出生信息都不一致 3: 身份证证号和出生信息不一致 4: 身份证证号和性别信息不一致"`
	WordResults      struct {
		Name         WordResult `json:"姓名"`
		Nation       WordResult `json:"民族"`
		Address      WordResult `json:"住址"`
		Sex          WordResult `json:"性别"`
		BirthDay     WordResult `json:"出生"`
		IDCradNumber WordResult `json:"公民身份号码"`
		TimeOutData  WordResult `json:"失效日期"`
		Organization WordResult `json:"发证机关"`
		IssueDate    WordResult `json:"签发日期"`
	} `json:"words_result" description:"定位和识别结果数组"`
	WordsResultNum uint32 `json:"words_result_num" description:"识别结果数，表示words_result的元素个数"`
}

type IDCardResponse struct {
}

type LocationInt struct {
	Top    int `json:"top" description:"位置,顶部"`
	Left   int `json:"left" description:"位置,左"`
	Width  int `json:"width" description:"位置,宽度"`
	Height int `json:"height" description:"位置,高度"`
}

type Locationfloat struct {
	Top    float64 `json:"top" description:"位置,顶部"`
	Left   float64 `json:"left" description:"位置,左"`
	Width  float64 `json:"width" description:"位置,宽度"`
	Height float64 `json:"height" description:"位置,高度"`
}

type WordResult struct {
	PosInImg LocationInt `json:"location" description:"在图像中的位置"`
	Words    string      `json:"words" description:"字段内容"`
}

type IDCardRecognitionResponse struct {
	Name         string `json:"姓名"`
	Nation       string `json:"民族"`
	Address      string `json:"住址"`
	Sex          string `json:"性别"`
	BirthDay     int64  `json:"出生"`
	IDCradNumber string `json:"公民身份证号码"`
	TimeOutDate  int64  `json:"失效日期"`
	Organization string `json:"签发机关"`
	IssueDate    int64  `json:"签发日期"`
}

type GetAccessTokenResponse struct {
	Token string `json:"access_token" description:"token"`
}

type BankCardRes struct {
	Direction int   `json:"direction" description:"图像方向，当图像旋转时，返回该参数。- 1:未定义，- 0:正向，- 1: 逆时针90度，- 2:逆时针180度，- 3:逆时针270度"`
	LogID     int64 `json:"log_id" description:"请求标识码，随机数，唯一的log id，用于问题定位"`
	Result    struct {
		Number    string `json:"bank_card_number" description:"银行卡卡号,样例:3568 8900 8000 0005"`
		VaildDate string `json:"valid_date" description:"有效期,样例:07/21"`
		CardType  int    `json:"bank_card_type" description:"银行卡类型，0:不能识别; 1: 借记卡; 2: 贷记卡（原信用卡大部分为贷记卡）; 3: 准贷记卡; 4: 预付费卡;"`
		Bank      string `json:"bank_name" description:"银行名,招商银行"`
	} `json:"result" description:"返回结果"`
}

type DrivingLicenseRes struct {
	LogID          int64 `json:"log_id" description:"请求标识码，随机数，唯一的log id，用于问题定位"`
	WordsResultNum int   `json:"words_result_num" description:"识别结果数，表示words_result的元素个数"`
	Result         struct {
		Name             Words `json:"姓名" description:"王桃桃"`
		Expire           Words `json:"至" description:"有效期结束时间,20210518"`
		Birth            Words `json:"出生日期" description:"19880929"`
		Number           Words `json:"证号" description:"210282198809294228"`
		Address          Words `json:"住址" description:"辽宁省大连市甘井子区"`
		FirstLicenseTime Words `json:"初次领证日期" description:"20150518"`
		Nation           Words `json:"国籍" description:"中国"`
		Level            Words `json:"准驾车型" description:"C1"`
		Sex              Words `json:"性别" description:"女"`
		VaildTimeStart   Words `json:"有效期限" description:"20150518"`
	} `json:"words_result" description:"识别结果"`
}

type Words struct {
	Word string `json:"words" description:""`
}

type VehicleLicenseRes struct {
	LogID          int64 `json:"log_id" description:"请求标识码，随机数，唯一的log id，用于问题定位"`
	WordsResultNum int   `json:"words_result_num" description:"识别结果数，表示words_result的元素个数"`
	Result         struct {
		CarCode         Words `json:"车辆识别代号" description:"SSVUDDTT2J2022558"`
		Address         Words `json:"住址" description:"中牟县三刘寨村"`
		LicenseTime     Words `json:"发证日期" description:"20180313"`
		CarModel        Words `json:"品牌型号" description:"大众汽车牌SVW6474DFD"`
		CarType         Words `json:"车辆类型" description:"小型普通客车"`
		Owner           Words `json:"所有人" description:"郑昆"`
		Usage           Words `json:"使用性质" description:"非营运"`
		EngineNumber    Words `json:"发动机号码" description:"111533"`
		PlateNumber     Words `json:"号牌号码" description:"豫A99RR9"`
		RegisterTime    Words `json:"注册日期" description:"20180312"`
		RecordID        Words `json:"档案编号" description:"320601972272"`
		Load            Words `json:"核定载人数" description:"2人"`
		Weight          Words `json:"总质量" description:"242kg'"`
		EquipmentWeight Words `json:"装备质量" description:"91kg"`
		LoadWeight      Words `json:"核定载质量" description:"--"`
		Remark          Words `json:"备注" description:"强制报废期止2030-05-15"`
		CheckRecord     Words `json:"检验记录" description:"检验有效期至2021年05月苏F(00)"`
		Size            Words `json:"尺寸" description:"1770X735X1060mm"`
	} `json:"words_result" description:"识别结果"`
}

type PlateNumberRes struct {
	ErrorNum int    `json:"errno" description:"0"`
	Msg      string `json:"msg" description:"success"`
	Data     struct {
		LogID  int64 `json:"log_id" description:"请求标识码，随机数，唯一的log id，用于问题定位"`
		Result struct {
			Color       string `json:"color" description:"车牌颜色：支持blue、green、yellow"`
			Number      string `json:"number" description:"车牌号,苏AD12267"`
			Probability []int  `json:"probability" description:"车牌中每个字符的置信度，区间为0-1"`
			Location    []struct {
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"vertexes_location" description:"返回文字外接多边形顶点位置"`
		} `json:"words_result" description:"结果"`
	} `json:"data" description:"车牌数据"`
}

type VinRes struct {
	LogID          int64  `json:"log_id" description:"请求标识码，随机数，唯一的log id，用于问题定位"`
	WordsResultNum int    `json:"words_result_num" description:"识别结果数，表示words_result的元素个数"`
	ErrorCode      int    `json:"error_code" description:"错误编码"`
	ErrorMsg       string `json:"error_msg" description:"错误编码"`
	Result         []struct {
		Location LocationInt `json:"location" description:"在图像中的位置"`
		Vin      string      `json:"words" description:"vin码识别结果"`
	} `json:"words_result" description:"结果"`
}

type CarTypeRes struct {
	LogID       int64         `json:"log_id" description:"请求标识码，随机数，唯一的log id，用于问题定位"`
	Location    Locationfloat `json:"location_result" description:"在图像中的位置"`
	ColorResult string        `json:"color_result" description:"车身颜色"`
	Result      []struct {
		BaikeInfo struct {
			BaikeURL    string `json:"baike_url" description:"对应车型识别结果百度百科页面链接"`
			Description string `json:"description" description:"对应车型识别结果百科内容描述"`
			ImgUrl      string `json:"img_url" description:"对应车型识别结果百科图片链接"`
		} `json:"baike_info,omitempty" description:"对应车型识别结果的百科词条名称"`
		Score float64 `json:"score" description:"置信度，取值0-1，示例：0.5321"`
		Name  string  `json:"name" description:"车型名称，示例：宝马x6"`
		Year  string  `json:"year" description:"年份"`
	} `json:"result"`
}

type BusinessLicenseRes struct {
	LogID          int64 `json:"log_id"`
	WordsResultNum int   `json:"words_result_num"`
	Direction      int   `json:"direction" description:"方向,-1:未定义 0:正向， 1: 逆时针90度,2:逆时针180度， 3:逆时针270度"`
	WordsResult    struct {
		Code struct {
			Words    string      `json:"words" description:"10440119MA06M8503"`
			Location LocationInt `json:"location"`
		} `json:"社会信用代码"`
		MakeUp struct {
			Words    string      `json:"words" description:"无"`
			Location LocationInt `json:"location"`
		} `json:"组成形式"`
		Range struct {
			Words    string      `json:"words" description:"商务服务业"`
			Location LocationInt `json:"location"`
		} `json:"经营范围"`
		SetUpTime struct {
			Words    string      `json:"words" description:"2019年01月01日"`
			Location LocationInt `json:"location"`
		} `json:"成立日期"`
		Corporate struct {
			Words    string      `json:"words" description:"方平"`
			Location LocationInt `json:"location"`
		} `json:"法人"`
		RegisteredCapital struct {
			Words    string      `json:"words" description:"200万元"`
			Location LocationInt `json:"location"`
		} `json:"注册资本"`
		ID struct {
			Words    string      `json:"words" description:"921MA190538210301"`
			Location LocationInt `json:"location"`
		} `json:"证件编号"`
		Address struct {
			Words    string      `json:"words" description:"广州市"`
			Location LocationInt `json:"location"`
		} `json:"地址"`
		Name struct {
			Words    string      `json:"words" description:"有限公司"`
			Location LocationInt `json:"location"`
		} `json:"单位名称"`
		Validity struct {
			Words    string      `json:"words" description:"长期"`
			Location LocationInt `json:"location"`
		} `json:"有效期"`
		Type struct {
			Words    string      `json:"words" description:"有限责任公司(自然人投资或控股)"`
			Location LocationInt `json:"location"`
		} `json:"类型"`
	} `json:"words_result"`
}

type NumberRes struct {
	LogID          int `json:"log_id"`
	WordsResultNum int `json:"words_result_num"`
	WordsResult    []struct {
		Location LocationInt `json:"location"`
		Words    string      `json:"words"`
	} `json:"words_result"`
}
