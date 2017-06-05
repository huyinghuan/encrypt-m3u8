package main

import (
	"encry/encrypt"
	"encry/utils"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"strings"

	"net/url"

	"github.com/grafov/m3u8"
)

func server() {
	key := "0123456789123456"
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	http.HandleFunc("/key", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "0123456789123456")
	})
	http.HandleFunc("/encrypt", func(w http.ResponseWriter, req *http.Request) {
		values := req.URL.Query()
		urls, ok := values["url"]
		if !ok {
			w.WriteHeader(406)
			return
		}
		log.Println(urls[0])
		stream, err := utils.GetStream(urls[0])
		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, "获取数据流失败")
			return
		}
		encryptStream, err := encrypt.CBCEncryptStream(stream, key, iv)
		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, "加密流失败")
			return
		}
		w.WriteHeader(200)
		head := w.Header()
		head.Set("content-type", "application/x-mpegURL")
		w.Write(encryptStream)

	})
	http.ListenAndServe(":8080", nil)
}

func main() {
	downloadM3U8("http://175.6.246.26/c1/2017/06/03_0/70376941E914FED2B04542C0C5B02EB7_20170603_1_1_1244_mp4/0F070BFEE8B85C4C6895F48F0AB8FB98.m3u8?t=593539d3&pno=1000&sign=d5f4670e152a25c7e70484012e00f852&ld=1496631585747&win=3600&srgid=26&urgid=1556&srgids=26&nid=922&payload=usertoken%3Dhit%3D1%5Eruip%3D2095616645&rdur=21600&limitrate=0&fid=70376941E914FED2B04542C0C5B02EB7&ver=0x03&uuid=aa39b90d8e7f458ea1ece776757efb4d&arange=0&yfweb=1")
	server()
}

func download(fileURL string) {
	config, _ := utils.ReadConfig(".yaml")
	dist := config.Download
	fileURLObj, _ := url.Parse(fileURL)
	urlArr := strings.Split(fileURLObj.Path, "/")
	filename := urlArr[len(urlArr)-1]
	key := "0123456789123456"
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

func downloadTSList(url string, list []string) {
	startTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(len(list))
	resourceURLChan := make(chan string, 20)
	//开启10个并行
	for i := 0; i < 10; i++ {
		go func() {
			for {
				resourceURL := <-resourceURLChan
				download(resourceURL)
				wg.Done()
			}

		}()
	}
	for _, resource := range list {
		resourceURLChan <- (url + resource)
	}
	wg.Wait()
	config, _ := utils.ReadConfig(".yaml")
	downloadFiles, _ := ioutil.ReadDir(config.Download)
	fmt.Printf("总共需要下载文件%d 个, 实际下载文件 %d 个 \n", len(list), len(downloadFiles)-1)
	fmt.Println(time.Now().Sub(startTime))
}

func downloadM3U8(m3u8URL string) {
	resp, err := http.Get(m3u8URL)
	if err != nil {
		// handle err
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("Error: response status is %d", resp.StatusCode)
	}
	p, listType, err := m3u8.DecodeFrom(resp.Body, true)
	//_, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}
	urlObj, _ := url.Parse(m3u8URL)
	urlSplit := strings.Split(urlObj.Path, "/")
	urlSplit = urlSplit[:(len(urlSplit) - 1)]
	directory := urlObj.Scheme + "://" + urlObj.Host + strings.Join(urlSplit, "/") + "/"

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

	list := strings.Split(m3u8String, "\n")
	tsList := []string{}
	for _, line := range list {
		if strings.Index(line, ".ts") != -1 {
			tsList = append(tsList, line)
		}
	}
	//----------------------- 结束

	//生成新的m3u8
	newm3u8 := []string{}
	hasAppend := false
	for _, line := range list {
		if strings.Index(line, ".ts") == -1 {
			if strings.Index(line, "#EXTINF") != -1 && !hasAppend {
				newm3u8 = append(newm3u8, "#EXT-X-KEY:METHOD=AES-128,URI=\"http://172.28.209.248:8080/key\",IV=0x0000000000000000")
				hasAppend = true
			}
			newm3u8 = append(newm3u8, line)
			continue
		}
		newm3u8 = append(newm3u8, getNewTSLink(line))
	}
	newm3u8String := strings.Join(newm3u8, "\n")

	ioutil.WriteFile("/Users/hyh/Downloads/encrypt/download/all.m3u8", []byte(newm3u8String), 0644)
	downloadTSList(directory, tsList)
}

func getNewTSLink(line string) string {
	urlObj, _ := url.Parse(line)
	path := urlObj.Path
	arr := strings.Split(path, "/")
	filename := arr[len(arr)-1]
	return "http://172.28.209.248:14422/" + filename
}
