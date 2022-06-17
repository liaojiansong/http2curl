/*
Copyright © 2022 liaojiansong

*/
package cmd

import (
	"github.com/spf13/cobra"
	"http2curl/impl"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web service",
	Long:  `Open the web service, and you can easily convert http message.`,
	Run: func(cmd *cobra.Command, args []string) {
		impl.Init(cmd).RunWebServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntP("port", "p", 22330, "listen port")
	serveCmd.Flags().String("log-level", "info", "set log level (debug|level|warn|error)")
	serveCmd.Flags().String("log-path", "/usr/log/http2curl", "Set log save path")
}
