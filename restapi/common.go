package restapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/valyala/fasthttp"
)

func ping(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("pong")
}

func writeError(ctx *fasthttp.RequestCtx, statusCode int, message string) {
	res := map[string]interface{}{
		"message": message,
	}

	ctx.SetStatusCode(statusCode)
	writeJSON(ctx, res)
}

func writeOK(ctx *fasthttp.RequestCtx, data interface{}) {
	ctx.SetStatusCode(http.StatusOK)
	writeJSON(ctx, data)
}

func writeJSON(ctx *fasthttp.RequestCtx, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Println("error : ", err)
		ctx.WriteString(err.Error())
		return
	}

	_, err = ctx.Write(b)
	if err != nil {
		log.Println("error : ", err)
	}
}
