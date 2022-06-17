/*
Copyright © 2022 liaojiansong

*/
package cmd

import (
	"github.com/spf13/cobra"
	"http2curl/impl"
)

var parserCmd = &cobra.Command{
	Use:   "cli",
	Short: "Convert the http message to the curl command line",
	Long: `This tool can read the http message in the local file and convert it to the curl command line. 
	Note that the source file can only hold one http message, and multiple messages will be treated as one.`,
	Run: func(cmd *cobra.Command, args []string) {
		impl.Cli(cmd)
	},
}

func init() {
	rootCmd.AddCommand(parserCmd)
	parserCmd.Flags().StringP("file", "f", "", "file path where http messages are stored")
	err := parserCmd.MarkFlagRequired("file")
	if err != nil {
		panic(err)
	}
}
