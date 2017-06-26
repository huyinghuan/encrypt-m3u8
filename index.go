package main

import (
	"encry/resolve"

	"log"

	iris "gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
)

func server() {
	key := "0123456789123456"
	app := iris.New()

	app.Adapt(httprouter.New())

	app.Get("/key", func(ctx *iris.Context) {
		ctx.Writef(key)
	})
	app.Get("/m3u8", func(ctx *iris.Context) {
		videoID := ctx.URLParam("videoid")
		terminalType := ctx.URLParam("terminal")
		resolution := ctx.URLParam("resolution")
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
	server()
}
