package cmd

import "github.com/spf13/cobra"

var watchUsage = `Usage: cw watch IMAGEDIGEST

Options: 
  -h, --help  help for watch
`

var WatchCmd = &cobra.Command {
	Use: "watch",
	Short: "Watch what files are used in an image",
	Long: "Watch what files are used in an image",
	Args: ,
	Run: func(cmd *cobra.Command, args []string){

	},
}

func init() {
	RootCmd.AddCommand(WatchCmd)

}