package impl

import (
	_ "embed"
	"go.uber.org/zap"
	"http2curl/pkg/log"
	"io/ioutil"
	"net/http"
)

//go:embed static/textarea.html
var indexHtml string

func index(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(indexHtml))
	if err != nil {
		log.Error("send index failed", zap.Error(err))
		badRsp(w, "Send index.html failed,please try again")
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
		badRsp(w, "http msg is empty")
		return
	}
	command, err := NewConverter(msg).toCommands()
	if err != nil {
		badRsp(w, err.Error())
		return
	}
	_, err = w.Write([]byte(command.String()))
	if err != nil {
		badRsp(w, err.Error())
		return
	}
	return
}

func badRsp(writer http.ResponseWriter, msg string) {
	writer.Write([]byte(msg))
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
