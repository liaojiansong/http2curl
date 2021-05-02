package impl

import (
	_ "embed"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	HINT = "something was wrong,please try again!"
	PORT = ":4887"
)

//go:embed textarea.html
var indexHtml string

func Web() {
	http.HandleFunc("/", index)
	http.HandleFunc("/conv", warp(conv))
	log.Printf("server was start,listen port%s", PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatalf("start server failed;\n%s", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(indexHtml))
	if err != nil {
		badRsp(w, HINT)
		return
	}
	return
}

func conv(w http.ResponseWriter, r *http.Request) {
	msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		badRsp(w, err.Error())
		return
	}
	if len(msg) == 0 {
		badRsp(w, "msg is empty")
		return
	}
	bMsg := string(msg)
	log.Printf("reveiced msg:\n%s", bMsg)
	converter, err := NewConverter(bMsg)
	if err != nil {
		badRsp(w, err.Error())
		return
	}
	command, err := converter.do()
	if err != nil {
		badRsp(w, err.Error())
		return
	}
	log.Printf("return command:%s\n", command)
	w.Write([]byte(command))
	return
}

type CurlRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func badRsp(writer http.ResponseWriter, mgs string) {
	writer.Write([]byte(mgs))
	return
}

func warp(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		//w.Header().Set("content-type", "application/json;charset=UTF-8")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		f(w, r)
	}
}
