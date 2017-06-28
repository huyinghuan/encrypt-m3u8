package main

import (
	"encry/resolve"
	"encry/service"
	"encry/utils"
	"log"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

func init() bool {
	if err := utils.InitConfig(); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func server() {
	if !init() {
		return
	}

	app := iris.New()

	app.Get("/key", func(ctx context.Context) {
		ctx.Writef(key)
	})
	app.Get("/m3u8", func(ctx context.Context) {
		videoID := ctx.URLParam("videoid")
		terminalType := ctx.URLParam("terminal")
		resolution := ctx.URLParam("resolution")

		if videoID == "" || terminalType == "" || resolution == "" {
			ctx.SetStatusCode(403)
			return
		}

		//判断该m3u8是否已存在
		if service.IsExistM3U8(terminalType, videoID, resolution) {
			//生成新的m3u8返回
			return
		}
		//获取原始m3u8地址
		originSourceURL, err := resolve.GetM3U8OriginSourceURL(videoID, terminalType, resolution)
		if err != nil {
			ctx.SetStatusCode(500)
			log.Fatal(err)
			return
		}
		originSource, err := resolve.GetM3U8OriginSource(originSourceURL)
		if err != nil {
			ctx.SetStatusCode(500)
			log.Fatal(err)
			return
		}

		resolve.DownloadM3U8TSList(originSource)

		// ctx.Redirect(originSourceURL, 302)
		// return
	})

	app.Listen(":8080")
}

func main() {
	server()
}
