package impl

import (
	"bufio"
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"http2curl/pkg/log"
	"io"
	"net/http"
	"strings"
)

type commands []string

func InitCommand() *commands {
	return &commands{"curl"}
}

func (c *commands) add(s ...string) {
	*c = append(*c, s...)
}

func (c *commands) String() string {
	return strings.Join(*c, " ")
}

type Converter struct {
	msg       []byte
	msgReader *bufio.Reader
}

func NewConverter(msg []byte) *Converter {
	c := &Converter{
		msg:       msg,
		msgReader: bufio.NewReader(bytes.NewReader(msg)),
	}
	return c
}

func (c *Converter) toCommands() (*commands, error) {
	req, err := http.ReadRequest(c.msgReader)
	if err != nil {
		log.Error("read request failed", zap.String("err", err.Error()))
		return nil, err
	}
	command := InitCommand()

	schema := req.URL.Scheme
	url := req.URL.String()
	if schema == "" {
		schema = "http"
		if req.TLS != nil {
			schema = "https"
		}
		url = schema + "://" + req.Host + req.URL.Path
	}
	command.add("-X", req.Method)

	for key, values := range req.Header {
		h := fmt.Sprintf("%s: %s", key, strings.Join(values, " "))
		command.add("-H", c.escape(h))
	}

	if req.Body != nil {
		var buff bytes.Buffer
		n, err := buff.ReadFrom(req.Body)
		// sometime Content-Length maybe wrong,log this error and go on
		if err == io.ErrUnexpectedEOF {
			log.Info("Content-Length wrong", zap.Int64("original len", req.ContentLength), zap.Int64("actual", n))
		}

		if err != nil && err != io.ErrUnexpectedEOF {
			log.Error("read request body failed", zap.String("err", err.Error()))
			return nil, err
		}
		body := buff.String()
		if len(body) > 0 {
			command.add("-d", c.escape(string(body)))
		}
	}
	command.add(url)
	return command, nil
}

func (c *Converter) escape(s string) string {
	return `'` + strings.Replace(s, `'`, `'\''`, -1) + `'`
}
