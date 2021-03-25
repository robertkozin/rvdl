package server

import (
	"log"
	"net/http"
)


func ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.URL)

	switch req.URL.Path {
	case "/":
		home(res, req)
	case "/favicon.ico", "/apple-touch-icon.png":
		favicon(res, req)
	case "/robots.txt":
		robots(res, req)
	default:
		handleRvdl(res, req)
	}
}

func home(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Cache-Control", "public, max-age=86400")
	//if req.URL.Host != "redditvideodownload.com" {
	//	redirect(res, req, "https://redditvideodownload.com/", http.StatusFound)
	//	return
	//}
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.ServeFile(res, req, "web/static/index.txt")
}

func favicon(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Cache-Control", "public, max-age=604800")
	switch req.URL.Path {
	case "/favicon.ico":
		res.Header().Set("Content-Type", "image/x-icon")
		http.ServeFile(res, req, "./web/static/favicon.ico")
	case "/apple-touch-icon.png":
		res.Header().Set("Content-Type", "image/png")
		http.ServeFile(res, req, "./web/static/apple-touch-icon.png")
	}
}

func robots(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Cache-Control", "public, max-age=604800")
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.ServeFile(res, req, "web/static/robots.txt")
}
