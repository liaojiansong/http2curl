/*
Copyright © 2022 liaojiansong

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "http2curl",
	Short: "A toolkit that can convert http messages into curl commands",
	Long: `You can specify the '-serve' option to start a web service and use it in the browser, which is more recommended. 
Or specify the '-cli' and '--file' options to read the http message in the local file and output the conversion result to the console`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
