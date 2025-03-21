package fileUtil

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gyf841010/pz-infra-new/commonUtil"
	"github.com/gyf841010/pz-infra-new/log"
	"github.com/h2non/filetype"

	"github.com/astaxie/beego"
)

var fileTypeMap sync.Map

func init() {
	fileTypeMap.Store("ffd8ffe000104a464946", "jpg")  //JPEG (jpg)
	fileTypeMap.Store("89504e470d0a1a0a0000", "png")  //PNG (png)
	fileTypeMap.Store("47494638396126026f01", "gif")  //GIF (gif)
	fileTypeMap.Store("49492a00227105008037", "tif")  //TIFF (tif)
	fileTypeMap.Store("424d228c010000000000", "bmp")  //16色位图(bmp)
	fileTypeMap.Store("424d8240090000000000", "bmp")  //24位位图(bmp)
	fileTypeMap.Store("424d8e1b030000000000", "bmp")  //256色位图(bmp)
	fileTypeMap.Store("41433130313500000000", "dwg")  //CAD (dwg)
	fileTypeMap.Store("3c21444f435459504520", "html") //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c68746d6c3e0", "html")        //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c21646f637479706520", "htm")  //HTM (htm)
	fileTypeMap.Store("48544d4c207b0d0a0942", "css")  //css
	fileTypeMap.Store("696b2e71623d696b2e71", "js")   //js
	fileTypeMap.Store("7b5c727466315c616e73", "rtf")  //Rich Text Format (rtf)
	fileTypeMap.Store("38425053000100000000", "psd")  //Photoshop (psd)
	fileTypeMap.Store("46726f6d3a203d3f6762", "eml")  //Email [Outlook Express 6] (eml)
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "doc")  //MS Excel 注意：word、msi 和 excel的文件头一样
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "vsd")  //Visio 绘图
	fileTypeMap.Store("5374616E64617264204A", "mdb")  //MS Access (mdb)
	fileTypeMap.Store("252150532D41646F6265", "ps")
	fileTypeMap.Store("255044462d312e350d0a", "pdf")  //Adobe Acrobat (pdf)
	fileTypeMap.Store("2e524d46000000120001", "rmvb") //rmvb/rm相同
	fileTypeMap.Store("464c5601050000000900", "flv")  //flv与f4v相同
	fileTypeMap.Store("00000020667479706d70", "mp4")
	fileTypeMap.Store("49443303000000002176", "mp3")
	fileTypeMap.Store("000001ba210001000180", "mpg") //
	fileTypeMap.Store("3026b2758e66cf11a6d9", "wmv") //wmv与asf相同
	fileTypeMap.Store("52494646e27807005741", "wav") //Wave (wav)
	fileTypeMap.Store("52494646d07d60074156", "avi")
	fileTypeMap.Store("4d546864000000060001", "mid") //MIDI (mid)
	fileTypeMap.Store("504b0304140000000800", "zip")
	fileTypeMap.Store("526172211a0700cf9073", "rar")
	fileTypeMap.Store("235468697320636f6e66", "ini")
	fileTypeMap.Store("504b03040a0000000000", "jar")
	fileTypeMap.Store("4d5a9000030000000400", "exe")        //可执行文件
	fileTypeMap.Store("3c25402070616765206c", "jsp")        //jsp文件
	fileTypeMap.Store("4d616e69666573742d56", "mf")         //MF文件
	fileTypeMap.Store("3c3f786d6c2076657273", "xml")        //xml文件
	fileTypeMap.Store("494e5345525420494e54", "sql")        //xml文件
	fileTypeMap.Store("7061636b616765207765", "java")       //java文件
	fileTypeMap.Store("406563686f206f66660d", "bat")        //bat文件
	fileTypeMap.Store("1f8b0800000000000000", "gz")         //gz文件
	fileTypeMap.Store("6c6f67346a2e726f6f74", "properties") //bat文件
	fileTypeMap.Store("cafebabe0000002e0041", "class")      //bat文件
	fileTypeMap.Store("49545346030000006000", "chm")        //bat文件
	fileTypeMap.Store("04000000010000001300", "mxp")        //bat文件
	fileTypeMap.Store("504b0304140006000800", "docx")       //docx文件
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "wps")        //WPS文字wps、表格et、演示dps都是一样的
	fileTypeMap.Store("6431303a637265617465", "torrent")
	fileTypeMap.Store("6D6F6F76", "mov")         //Quicktime (mov)
	fileTypeMap.Store("FF575043", "wpd")         //WordPerfect (wpd)
	fileTypeMap.Store("CFAD12FEC5FD746F", "dbx") //Outlook Express (dbx)
	fileTypeMap.Store("2142444E", "pst")         //Outlook (pst)
	fileTypeMap.Store("AC9EBD8F", "qdf")         //Quicken (qdf)
	fileTypeMap.Store("E3828596", "pwl")         //Windows Password (pwl)
	fileTypeMap.Store("2E7261FD", "ram")         //Real Audio (ram)
}

const chunkSize = 64000

//二进制读取比较两个文件是否相同
func CompareFile(file1, file2 string) (bool, error) {
	f1, err := os.Open(file1)
	if err != nil {
		log.Error(err)
		return false, err
	}

	f2, err := os.Open(file2)
	if err != nil {
		log.Error(err)
		return false, err
	}

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true, nil
			} else if err1 == io.EOF || err2 == io.EOF {
				return false, nil
			} else {
				log.Error(err)
				return false, err
			}
		}

		if !bytes.Equal(b1, b2) {
			return false, nil
		}
	}

	return false, nil
}

func ExistFile(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	}
	return false
}

//获取文件的md5
func FileMD5(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	var result []byte
	result = hash.Sum(result)
	return fmt.Sprintf("%x", result), nil
}

//获取文件大小,单位字节
func GetFileSize(file string) (int, error) {
	fi, err := os.Stat(file)
	if err != nil {
		log.Error(err)
		return -1, err
	}
	return int(fi.Size()), nil
}

func GetLocalTempDir() string {
	dir := "../temp/" + strings.Replace(commonUtil.UUID(), "-", "", -1)
	os.MkdirAll(dir, 0755)
	return dir
}

// 获取前面结果字节的二进制
func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}

// 用文件前面几个字节来判断
// fSrc: 文件字节流（就用前面几个字节）
func GetFileType(fSrc []byte) string {
	var fileType string
	fileCode := bytesToHexString(fSrc)

	fileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
			strings.HasPrefix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	return fileType
}

// 通过文件头判断文件类型
// 实现:github.com/h2non/filetype
// 不支持txt格式,纯文本无文件头 2021-2-23
func GetFileTypeNew(fSrc []byte) string {
	fileType, err := filetype.Match(fSrc)
	if err != nil {
		return ""
	}

	if fileType == filetype.Unknown || fileType.Extension == "" {
		return ""
	}
	return fileType.Extension
}

func GetFileKey(fSrc []byte) string {
	datePrefix, _ := beego.AppConfig.Bool("datePrefixFlag")
	fileType := GetFileType(fSrc)
	if fileType == "" {
		fileType = "png"
	}

	var fileKey string
	if datePrefix {
		fileKey = fmt.Sprintf("%d/%d/", time.Now().Year(), time.Now().Month()) + commonUtil.UUID() + "." + fileType
	} else {
		fileKey = commonUtil.UUID() + "." + fileType
	}

	return fileKey
}

/**
 * @description: 生成文件oss key
 * @param {[]byte} fSrc 文件内容
 * @param {...string} fileTypes,文件类型(filetype失败xlsx文件类型还是容易错误,修改函数,支持传入文件类型,对于知道类型的场景优化)
 * @return {*} 文件oss key
 */
func GetFileKeyNew(fSrc []byte, fileTypes ...string) string {
	datePrefix, _ := beego.AppConfig.Bool("datePrefixFlag")
	fileType := ""
	if len(fileTypes) > 0 {
		for _, item := range fileTypes {
			if item != "" {
				fileType = item
				break
			}
		}
	}
	if fileType == "" {
		fileType = GetFileTypeNew(fSrc)
		if fileType == "" {
			fileType = "png"
		}
	}

	var fileKey string
	if datePrefix {
		fileKey = fmt.Sprintf("%d/%d/", time.Now().Year(), time.Now().Month()) + commonUtil.UUID() + "." + fileType
	} else {
		fileKey = commonUtil.UUID() + "." + fileType
	}

	return fileKey
}
