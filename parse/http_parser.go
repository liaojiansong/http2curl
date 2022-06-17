package parse

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
)

func contentLengthPos(msg []byte) (start, end int) {
	start = bytes.Index(msg, []byte("Content-Length"))
	if start == -1 {
		return -1, -1
	}
	for i := start; i < len(msg); i++ {
		if msg[i] == '\n' {
			return start, i
		}
	}
	return -1, -1
}

func realLength(msg []byte) int {
	bl := len(msg)
	for i := 0; i < bl; i++ {
		if i+1 < bl && msg[i] == '\n' && msg[i+1] == '\n' {
			return bl - i - 2
		}
	}
	return -1
}

func fixMsg(msg []byte) []byte {
	length := realLength(msg)
	if length == -1 {
		return msg
	}
	start, end := contentLengthPos(msg)
	if start*end < 0 {
		return msg
	}
	// todo copy
	s := fmt.Sprintf("Content-Length: %d", length)
	newMsg := make([]byte, 0, len(msg)+len(s))
	newMsg = append(newMsg, msg[:start]...)
	newMsg = append(newMsg, s...)
	newMsg = append(newMsg, msg[end:]...)
	return newMsg
}

func EchoCurlCommand(msg string) {
	fix := fixMsg([]byte(msg))
	newReader := bytes.NewReader(fix)
	reader := bufio.NewReader(newReader)
	request, err := http.ReadRequest(reader)
	if err != nil {
		fmt.Printf("解析错误:%s", err.Error())
		return
	}
	command, err := GetCurlCommand(request)
	if err != nil {
		fmt.Printf("%#v\n", err)
		return
	}
	println(command.String())
}
