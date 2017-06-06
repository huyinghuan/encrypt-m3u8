package main

import (
	"encry/resolve"

	"log"

	"fmt"

	iris "gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
)

func server() {
	key := "0123456789123456"
	// //iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	// http.HandleFunc("/key", func(w http.ResponseWriter, req *http.Request) {
	// 	io.WriteString(w, key)
	// })

	// http.ListenAndServe(":8080", nil)
	app := iris.New()

	app.Adapt(httprouter.New())

	app.Get("/key", func(ctx *iris.Context) {
		ctx.Writef(key)
	})
	app.Get("/m3u8", func(ctx *iris.Context) {
		videoID := ctx.URLParam("videoid")
		terminalType := ctx.URLParam("terminal")
		resolution := ctx.URLParam("resolution")
		fmt.Println(videoID, terminalType, resolution)
		if videoID == "" || terminalType == "" || resolution == "" {
			ctx.SetStatusCode(403)
			return
		}
		encryptURL, err := resolve.GetEncryptURL(videoID, terminalType, resolution)
		if err != nil {
			ctx.SetStatusCode(500)
			log.Fatal(err)
			return
		}
		ctx.Redirect(encryptURL, 302)
		return
	})

	app.Listen(":8080")
}

func main() {
	//encrypttype.DownloadM3U8("http://175.6.246.26/c1/2017/06/03_0/70376941E914FED2B04542C0C5B02EB7_20170603_1_1_1244_mp4/0F070BFEE8B85C4C6895F48F0AB8FB98.m3u8?t=593539d3&pno=1000&sign=d5f4670e152a25c7e70484012e00f852&ld=1496631585747&win=3600&srgid=26&urgid=1556&srgids=26&nid=922&payload=usertoken%3Dhit%3D1%5Eruip%3D2095616645&rdur=21600&limitrate=0&fid=70376941E914FED2B04542C0C5B02EB7&ver=0x03&uuid=aa39b90d8e7f458ea1ece776757efb4d&arange=0&yfweb=1")
	server()
}
