package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "log"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Welcome to goShort URL shortner written in Golang!\n")
}

func Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := r.FormValue("url")
	fmt.Fprint(w, "url : %s", url)
}

func Redirect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    code := ps.ByName("code")
	url := "https://github.com/julienschmidt/httprouter?code=" + code
	http.Redirect(w, r, url, http.StatusFound)
}

func main() {
    router := httprouter.New()
    router.GET("/", Index)
    router.POST("/create/", Create)
	router.GET("/r/:code", Redirect)
    log.Fatal(http.ListenAndServe(":8080", router))
}

