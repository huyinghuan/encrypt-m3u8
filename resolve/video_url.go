package resolve

import (
	"encoding/json"
	"encry/utils"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"time"

	"encry/encrypt"

	"github.com/grafov/m3u8"
)

type Stream struct {
	Def   string
	Name  string
	Scale string
	Vip   string
	Url   string
}

type ResponseData struct {
	Stream       []Stream
	StreamDomain []string `json:"stream_domain"`
}

type ResponseBody struct {
	Code int
	Data ResponseData
	Msg  string
}

//GetCDNURL  根据videoid， 终端类型， 分辨率 获取cdn url列表
func GetCDNURL(vedioID string, terminalType string, resolution string) (string, error) {
	config := utils.ReadConfig()
	uri := fmt.Sprintf(config.Cdn, vedioID, terminalType)
	resp, err := http.Get(uri)
	if err != nil {
		// handle err
		log.Fatalln(err)
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("Error: response status is %d", resp.StatusCode)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	result := new(ResponseBody)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	streamDomain := result.Data.StreamDomain[0]
	streamList := result.Data.Stream
	for _, stream := range streamList {
		if stream.Def == resolution || stream.Name == resolution {
			return streamDomain + stream.Url, nil
		}
	}
	return "", nil
}

type M3U8Response struct {
	Ver        string
	Isothercdn string
	Info       string
	Status     string
	Loc        string
	T          string
	Idc        string
}

func GetM3U8URL(cdn string) (string, error) {
	resp, err := http.Get(cdn)
	if err != nil {
		// handle err
		log.Fatalln(err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("Error: response status is %d", resp.StatusCode)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	responseBody := new(M3U8Response)
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", err
	}
	return responseBody.Info, nil
}

//GetM3U8OriginSourceURL 获取 m3u8 原始地址
func GetM3U8OriginSourceURL(vedioID string, terminalType string, resolution string) (string, error) {
	var cdnURL string
	var m3u8URL string
	var err error
	cdnURL, err = GetCDNURL(vedioID, terminalType, resolution)
	if err != nil {
		return "", err
	}
	m3u8URL, err = GetM3U8URL(cdnURL)
	if err != nil {
		return "", err
	}
	return m3u8URL, nil
}

//GetM3U8OriginSource 获取m3u8原始内容
func GetM3U8OriginSource(m3u8URL string) (string, error) {
	resp, err := http.Get(m3u8URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error: response status is %d", resp.StatusCode)
	}
	p, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
		return "", err
	}
	//----------------------- 开始
	var m3u8String string
	switch listType {
	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)
		m3u8String = mediapl.String()
	case m3u8.MASTER:
		masterpl := p.(*m3u8.MasterPlaylist)
		m3u8String = masterpl.String()
	}
	return m3u8String, nil
}

var downloadTSListChan = make(chan string, 1000)

//PrepareDownloadM3U8TSList 准备下载m3u8里面的ts文件
func PrepareDownloadM3U8TSList(m3u8URL string, m3u8Content string) {
	list := strings.Split(m3u8Content, "\n")
	for _, line := range list {
		if strings.Index(line, ".ts") != -1 {
			downloadTSListChan <- m3u8URL + "/" + line
		}
	}
}

//下载文件
func DownloadTsFileByURL(url string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("下载 %s 失败!\n", url)
	}
	defer response.Body.Close()
	filename := utils.GetURLFilename(url)
	config := utils.ReadConfig()
	//open a file for writing
	file, err := os.Create(fmt.Sprintf(config.Origints, filename))
	if err != nil {
		fmt.Printf("创建下载文件失败%s", url)
	}
	defer file.Close()
	// Use io.Copy to just dump the response body to the file. This supports huge files
	body, _ := ioutil.ReadAll(response.Body)
	if err := ioutil.WriteFile(file.Name(), body, 0644); err != nil {
		fmt.Printf("写入文件 %s 失败!\n", url)
	}
}

func StartDownloadTSService() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				beDownloadURL := <-downloadTSListChan
				DownloadTsFileByURL(beDownloadURL)
			}
		}()
	}
}

func SaveOriginM3U8File(GetM3U8OriginSource string, filename string) error {
	config := utils.ReadConfig()
	distFilePath := fmt.Sprintf(config.Originm3u8, filename)
	file, err := os.Create(distFilePath)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(file.Name(), []byte(GetM3U8OriginSource), 0644); err != nil {
		fmt.Printf("写入文件 %s 失败!\n", distFilePath)
		return err
	}
	return nil
}

func EncryptM3U8(originSource string) (string, error) {
	config := utils.ReadConfig()
	list := strings.Split(originSource, "\n")
	newm3u8 := []string{}
	hasAppend := false
	//生成随机key用来加密流
	key := utils.RandString(16)
	encryptKey, err := encrypt.CFBEncryptString([]byte(config.Querykey), key)
	fmt.Printf("生成加密秘钥:%s\n", key)
	if err != nil {
		return "", err
	}
	for _, line := range list {
		if strings.Index(line, ".ts") != -1 {
			filename := strings.Split(line, "?")[0]
			//key,filename,time
			query, err := encrypt.CFBEncryptString([]byte(config.Querykey), fmt.Sprintf("%s,%s,%v", key, filename, time.Now().Unix()))
			if err != nil {
				return "", err
			}
			newm3u8 = append(newm3u8, fmt.Sprintf(config.Encrypttsurl, query))
		} else {
			if strings.Index(line, "#EXTINF") != -1 && !hasAppend {
				extXKey := fmt.Sprintf("#EXT-X-KEY:METHOD=AES-128,URI=\"%s\",IV=0x0000000000000000", fmt.Sprintf(config.Keyurl, encryptKey))
				newm3u8 = append(newm3u8, extXKey)
				hasAppend = true
			}
			newm3u8 = append(newm3u8, line)
		}
	}
	return strings.Join(newm3u8, "\n"), nil
}
