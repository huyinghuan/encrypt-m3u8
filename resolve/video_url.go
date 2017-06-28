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

//DownloadM3U8TSList 下载m3u8里面的ts文件
func DownloadM3U8TSList(m3u8Content string) {
	list := strings.Split(m3u8Content, "\n")
	tsList := []string{}
	for _, line := range list {
		if strings.Index(line, ".ts") != -1 {
			tsList = append(tsList, line)
		}
	}
	log.Println(tsList)
}
