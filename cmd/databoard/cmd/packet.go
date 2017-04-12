package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// packetCmd represents the packet command
var packetCmd = &cobra.Command{
	Use: "packet",
}

func init() {
	RootCmd.AddCommand(packetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// packetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	packetCmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
	}
}
