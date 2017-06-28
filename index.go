package main

import (
	"encry/resolve"
	"encry/service"
	"encry/utils"
	"log"
	"os"

	"fmt"

	"encry/encrypt"

	"strings"

	"io/ioutil"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

func initConfig() bool {
	if err := utils.InitConfig(); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func server() {
	if !initConfig() {
		return
	}
	//开启下载线程
	resolve.StartDownloadTSService()

	app := iris.New()
	config := utils.ReadConfig()
	app.Get("/key", func(ctx context.Context) {
		key := ctx.URLParam("key")
		if len(key) < 16 {
			ctx.StatusCode(iris.StatusNotAcceptable)
			return
		}
		decryptKey, err := encrypt.CFBDecryptString([]byte(config.Querykey), key)
		if err != nil {
			log.Println(err)
			ctx.StatusCode(iris.StatusNotAcceptable)
			return
		}
		ctx.WriteString(decryptKey)
	})
	app.Get("/encrypt.ts", func(ctx context.Context) {
		f := ctx.URLParam("f")
		if len(f) < 16 {
			ctx.StatusCode(iris.StatusNotAcceptable)
			return
		}
		querystring, err := encrypt.CFBDecryptString([]byte(config.Querykey), f)
		if err != nil {
			log.Println(err)
			ctx.StatusCode(iris.StatusNotAcceptable)
			return
		}
		//key,time,filename
		decryptContentArr := strings.Split(querystring, ",")
		if len(decryptContentArr) != 3 {
			ctx.StatusCode(iris.StatusNotAcceptable)
			return
		}
		key := decryptContentArr[0]
		filename := decryptContentArr[1]
		distFilePath := fmt.Sprintf(config.Origints, filename)
		if _, err := os.Stat(distFilePath); err != nil {
			ctx.StatusCode(iris.StatusNotFound)
			return
		}
		iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		fileContent, err := ioutil.ReadFile(distFilePath)
		if err != nil {
			log.Println(err)
			ctx.StatusCode(iris.StatusServiceUnavailable)
			return
		}
		body, err := encrypt.CBCEncryptStream(fileContent, key, iv)
		if err != nil {
			log.Println(distFilePath, err)
			ctx.StatusCode(iris.StatusServiceUnavailable)
			return
		}
		//ctx.Header("Content-Type", "video/mp2t")
		ctx.Write(body)
	})
	app.Get("/all.m3u8", func(ctx context.Context) {
		videoID := ctx.URLParam("videoid")
		terminalType := ctx.URLParam("terminal")
		resolution := ctx.URLParam("resolution")

		if videoID == "" || terminalType == "" || resolution == "" {
			ctx.StatusCode(403)
			return
		}
		originSource := ""
		//判断该m3u8是否已存在
		if service.IsExistM3U8(terminalType, videoID, resolution) {
			body, err := service.GetExistM3U8(terminalType, videoID, resolution)
			if err != nil {
				fmt.Println(err)
				ctx.StatusCode(500)
				return
			}
			originSource = string(body)
		} else {
			//获取原始m3u8地址
			originSourceURL, err := resolve.GetM3U8OriginSourceURL(videoID, terminalType, resolution)
			if err != nil {
				ctx.StatusCode(500)
				log.Println(err)
				return
			}
			originSourceDirURL := utils.GetDirname(originSourceURL)
			originSource, err = resolve.GetM3U8OriginSource(originSourceURL)
			if err != nil {
				ctx.StatusCode(500)
				log.Println(err)
				return
			}
			//放入数据通道
			go resolve.PrepareDownloadM3U8TSList(originSourceDirURL, originSource)
			//写入文件
			resolve.SaveOriginM3U8File(originSource, fmt.Sprintf(config.M3u8rule, terminalType, videoID, resolution))
		}
		//返回编码后的m3u8
		//编码m3u8
		content, err := resolve.EncryptM3U8(originSource)
		if err != nil {
			ctx.StatusCode(500)
			log.Println(err)
			return
		}
		ctx.Header("Content-Type", "application/x-mpegURL")
		ctx.Write([]byte(content))
	})

	app.Listen(":8080")
}

func main() {
	server()
}
