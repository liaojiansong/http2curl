package impl

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"http2curl/pkg/log"
	"io/ioutil"
	"os"
)

func Cli(cmd *cobra.Command) {
	log.Init()
	file, err := cmd.Flags().GetString("file")
	if err != nil {
		log.Fatal("Flags file failed", zap.Error(err))
	}
	if file == "" {
		log.Fatal("file path is empty")
	}
	RunCli(file)
	os.Exit(0)
}
func readMsgFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func RunCli(filePath string) {
	msg, err := readMsgFile(filePath)
	if err != nil {
		log.Fatal("read http msg from file failed", zap.String("err", err.Error()))
	}
	if len(msg) == 0 {
		log.Fatal("http msg is empty,exit!")
	}
	converter := NewConverter(msg)
	cmds, err := converter.toCommands()
	if err != nil {
		log.Error("conv failed", zap.Error(err))
		os.Exit(-1)
	}
	s := cmds.String()
	println(s)
}
