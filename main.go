package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/caarlos0/env/v8"
	"github.com/valyala/fasthttp"
)

type Config struct {
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	Port int    `env:"PORT" envDefault:"8000"`
}

var client = &http.Client{
	Transport: &http.Transport{
		MaxConnsPerHost: 0,
	},
	Timeout: time.Minute,
}

func main() {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		log.Fatalln(err)
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	log.Printf("listening on %s", addr)
	if err := fasthttp.ListenAndServe(addr, handler); err != nil {
		panic(err)
	}
}

func handler(ctx *fasthttp.RequestCtx) {
	url := string(ctx.URI().PathOriginal())[1:] + "?" + string(ctx.URI().QueryString())
	if url == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprint(ctx, "invalid url")
		return
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "uncaught error: %v", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "uncaught error: %v", err)
		return
	}
	defer resp.Body.Close()

	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "uncaught error: %v", err)
		return
	}

	ctx.SetContentType(resp.Header.Get("Content-Type"))
	ctx.Response.Header.SetContentLength(contentLength)

	buff := make([]byte, 1)
	for {
		_, err := resp.Body.Read(buff)
		if err == io.EOF {
			break
		}

		ctx.Write(buff)
	}
}
