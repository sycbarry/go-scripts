package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)


var listCmd = &cobra.Command{
    Use: "list",
    Short: "list all of the tasks.",
    Long: "long" ,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("list called")
    },
}



func init() {
    RootCmd.AddCommand(listCmd)
}
