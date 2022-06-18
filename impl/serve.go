package impl

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"http2curl/pkg/log"
	syslog "log"
	"net/http"
)

type serve struct {
	port     int
	logLevel string
	logPath  string
}

func parseFlags(cmd *cobra.Command) *serve {
	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		syslog.Fatalf("port invliad err:%v", err)
	}
	logLevel, err := cmd.Flags().GetString("log-level")
	if err != nil {
		syslog.Fatalf("loglevel invliad err:%v", err)
	}
	logPath, err := cmd.Flags().GetString("log-path")
	if err != nil {
		syslog.Fatalf("log-path invliad err:%v", err)
	}
	a := &serve{
		port:     port,
		logLevel: logLevel,
		logPath:  logPath,
	}
	return a
}

func Init(cmd *cobra.Command) *serve {
	a := parseFlags(cmd)
	logOpts := []log.Options{
		log.SetLevel(a.logLevel),
		log.AddPath(a.logPath),
	}
	log.Init(logOpts...)
	return a
}

func (a *serve) RunWebServer() {
	Web(a.port)
}

func Web(port int) {
	http.HandleFunc("/", index)
	http.HandleFunc("/conv", warp(conv))
	log.Info("starting web serve", zap.Int("listen port", port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("start web serve fail", zap.String("err", err.Error()))
	}
}
