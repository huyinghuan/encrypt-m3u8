package main

import (
	"encry/controller"
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
		decryptContent, err := encrypt.CFBDecryptString([]byte(config.Querykey), key)
		if err != nil {
			log.Println(err)
			ctx.StatusCode(iris.StatusNotAcceptable)
			return
		}
		decryptContentArr := []string{}
		if decryptContentArr = strings.Split(decryptContent, ";"); len(decryptContentArr) != 2 {
			ctx.StatusCode(iris.StatusNotAcceptable)
			return
		}
		//检测特征码
		if !service.MatchSignature(ctx.Request(), decryptContentArr[1]) {
			ctx.StatusCode(iris.StatusNotAcceptable)
			return
		}
		ctx.WriteString(decryptContentArr[0])
	})
	app.Get("/encrypt.ts", func(ctx context.Context) {
		f := ctx.URLParam("f")
		if len(f) < 16 {
			ctx.StatusCode(iris.StatusNotAcceptable)
			return
		}
		querystring := ""
		decryptContentArr := []string{}

		var err error
		if querystring, err = encrypt.CFBDecryptString([]byte(config.Querykey), f); err != nil {
			log.Println(err)
			ctx.StatusCode(iris.StatusNotAcceptable)
			return
		}
		decryptContentArr = strings.Split(querystring, ",")

		//key,time,filename
		if len(decryptContentArr) != 3 {
			ctx.StatusCode(iris.StatusNotAcceptable)
			return
		}
		//解密密钥
		key := ""
		keyArray := []string{}
		if keyArray = strings.Split(decryptContentArr[0], ";"); len(keyArray) != 2 {
			ctx.StatusCode(iris.StatusNotAcceptable)
			return
		}
		//检测特征码
		if !service.MatchSignature(ctx.Request(), keyArray[1]) {
			ctx.StatusCode(iris.StatusNotAcceptable)
			return
		}
		key = keyArray[0]

		filename := decryptContentArr[1]
		distFilePath := fmt.Sprintf(config.Origints, filename)
		if _, err := os.Stat(distFilePath); err != nil {
			ctx.StatusCode(iris.StatusNotFound)
			return
		}
		iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		//读取目标ts文件，并加密返回
		fileContent, err := ioutil.ReadFile(distFilePath)
		if err != nil {
			log.Println(err)
			ctx.StatusCode(iris.StatusServiceUnavailable)
			return
		}
		if key == config.Donotencrypt {
			ctx.Write(fileContent)
			return
		}
		body, err := encrypt.CBCEncryptStream(fileContent, key, iv)
		//fmt.Println("加密流密码：", key)
		if err != nil {
			log.Println(distFilePath, err)
			ctx.StatusCode(iris.StatusServiceUnavailable)
			return
		}
		ctx.Write(body)
	})
	app.Get("/decrypt", controller.Decrypt)
	app.Get("/decrypt-key", controller.DecryptKey)
	app.Get("/all.m3u8", func(ctx context.Context) {
		videoID := ctx.URLParam("videoid")
		terminalType := ctx.URLParam("terminal")
		resolution := ctx.URLParam("resolution")

		if videoID == "" || terminalType == "" || resolution == "" {
			ctx.StatusCode(403)
			return
		}
		originSource := ""
		//判断该m3u8是否已存在, 存在则读取源内容
		if service.IsExistM3U8(terminalType, videoID, resolution) {
			body, err := service.GetExistM3U8(terminalType, videoID, resolution)
			if err != nil {
				fmt.Println(err)
				ctx.StatusCode(500)
				return
			}
			originSource = string(body)
		} else {
			//不存在
			//获取原始m3u8地址
			originSourceURL, err := resolve.GetM3U8OriginSourceURL(videoID, terminalType, resolution)
			if err != nil {
				ctx.StatusCode(500)
				log.Println(err)
				return
			}
			originSourceDirURL := utils.GetDirname(originSourceURL)
			//获取原始m3u8的内容
			originSource, err = resolve.GetM3U8OriginSource(originSourceURL)
			if err != nil {
				ctx.StatusCode(500)
				log.Println(err)
				return
			}
			//放入数据通道，下载ts片段
			go resolve.PrepareDownloadM3U8TSList(originSourceDirURL, originSource)
			//将原始m3u8的内容写入文件，缓存以便二次使用
			resolve.SaveOriginM3U8File(originSource, fmt.Sprintf(config.M3u8rule, terminalType, videoID, resolution))
		}
		//返回编码后的m3u8
		//编码并加密m3u8
		content, err := resolve.EncryptM3U8(originSource, service.GetSignature(ctx.Request()))
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
