package controller

import (
	"encry/encrypt"
	"encry/utils"
	"log"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

func DecryptKey(ctx context.Context) {
	key := ctx.URLParam("key")
	config := utils.ReadConfig()
	if querystring, err := encrypt.CFBDecryptString([]byte(config.Querykey), key); err != nil {
		log.Println(err)
		ctx.StatusCode(iris.StatusNotAcceptable)
		return
	} else {
		ctx.WriteString(querystring)
	}
}
