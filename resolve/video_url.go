package resolve

import (
	"encoding/json"
	"encry/encrypt"
	"encry/utils"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

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
	config, _ := utils.ReadResourceConfig()
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

func download(fileURL string, config *utils.Config) {
	dist := config.Download
	filename := utils.GetURLFilename(fileURL)
	key := config.Key
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if filename == "" {
		fmt.Printf("文件名为空%s", fileURL)
		return
	}
	response, e := http.Get(fileURL)
	if e != nil {
		fmt.Printf("下载 %s 失败!\n", fileURL)
	}

	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(dist + filename)
	if err != nil {
		fmt.Printf("创建下载文件失败%s", filename)
	}
	defer file.Close()
	// Use io.Copy to just dump the response body to the file. This supports huge files
	body, _ := ioutil.ReadAll(response.Body)
	//加密
	content, _ := encrypt.CBCEncryptStream(body, key, iv)
	err = ioutil.WriteFile(file.Name(), content, 0644)
	if err != nil {
		fmt.Printf("写入文件 %s 失败!\n", fileURL)
	}
}

func downloadTSList(url string, list []string, config *utils.Config) {
	startTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(len(list))
	resourceURLChan := make(chan string, 20)
	//开启10个并行
	for i := 0; i < 10; i++ {
		go func() {
			for {
				resourceURL := <-resourceURLChan
				download(resourceURL, config)
				wg.Done()
			}

		}()
	}
	for _, resource := range list {
		resourceURLChan <- (url + resource)
	}
	wg.Wait()
	fmt.Printf("总共下载ts文件 %d， 耗时:", len(list))
	fmt.Println(time.Now().Sub(startTime))
}

func DownloadM3U8(m3u8URL string) string {
	config, _ := utils.ReadConfig()
	resp, err := http.Get(m3u8URL)
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("Error: response status is %d", resp.StatusCode)
		return ""
	}
	p, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	directory := utils.GetDirname(m3u8URL) + "/"
	//获取m3u8 里面的ts列表
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

	//生成新的m3u8
	newm3u8 := []string{}
	hasAppend := false

	list := strings.Split(m3u8String, "\n")
	//待下载列表
	tsList := []string{}
	for _, line := range list {
		if strings.Index(line, ".ts") != -1 {
			tsList = append(tsList, line)
			newm3u8 = append(newm3u8, config.Tsurl+utils.GetURLFilename(line))
		} else {
			if strings.Index(line, "#EXTINF") != -1 && !hasAppend {
				extXKey := fmt.Sprintf("#EXT-X-KEY:METHOD=AES-128,URI=\"%s\",IV=0x0000000000000000", config.Keyurl)
				newm3u8 = append(newm3u8, extXKey)
				hasAppend = true
			}
			newm3u8 = append(newm3u8, line)
		}

	}
	//----------------------- 结束

	//生成新的m3u8
	newm3u8String := strings.Join(newm3u8, "\n")
	randomFileName := utils.RandStringRunes(8)
	go func() {
		ioutil.WriteFile(fmt.Sprintf(config.M3u8, randomFileName), []byte(newm3u8String), 0644)
	}()
	go func() { downloadTSList(directory, tsList, config) }()
	return fmt.Sprintf(config.Finalurl, randomFileName)
}

func GetEncryptURL(vedioID string, terminalType string, resolution string) (string, error) {
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
	fmt.Println(m3u8URL)
	return DownloadM3U8(m3u8URL), nil
}
