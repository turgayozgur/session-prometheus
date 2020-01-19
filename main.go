package main

import (
	"github.com/valyala/fasthttp"
	"log"
)

func main(){
	m := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/pay":
			pay(ctx)
		}
	}
	log.Print("Listening on port 8080...")
	fasthttp.ListenAndServe(":8080", m)
}

func pay(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("OK")
}