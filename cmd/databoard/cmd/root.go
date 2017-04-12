package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "databoard",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
}

func exit(code int, v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(code)
}

func exitf(code int, format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
	if !strings.HasSuffix(format, "\n") {
		os.Stderr.Write([]byte{'\n'})
	}
	os.Exit(code)
}
