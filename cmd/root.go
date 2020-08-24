package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command {
	Use: "cw",
	Short: "A tool to monitor what files are used in a container",
	Long: "A tool to monitor what files are used in a container",
}

func init() {

}