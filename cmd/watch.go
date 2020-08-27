package cmd

import (
	"ContainerWatcher/pkg/file"
	"ContainerWatcher/pkg/log"
	"ContainerWatcher/pkg/watcher"
	"fmt"
	"os"
	"os/signal"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/spf13/cobra"
)

var watchUsage = `Usage: cw watch IMAGEDIGEST

Options: 
  -h, --help  help for watch
`

// WatchCmd watch command
var WatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch what files are used in an image",
	Long:  "Watch what files are used in an image",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		digest := args[0]

		w := watcher.NewWatcher(digest)

		c := w.Watch()

		var totalSize int64 = 0
		var uniqueFiles map[string]*file.File = make(map[string]*file.File)

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)

		for {
			select {
			case e := <-c:
				if _, ok := uniqueFiles[e.Name]; ok {
					uniqueFiles[e.Name].Times++
				} else {
					log.Logger.Info("File path: ", e.Name, "File size: ", e.Size)
					uniqueFiles[e.Name] = &file.File{Name: e.Name, Size: e.Size, Times: 1}
					totalSize += e.Size
				}
			case <-stop:
				log.Logger.Infof("Total size: %d MB", totalSize/1024/1024)

				func() {
					f := excelize.NewFile()
					r := 1
					for _, fi := range uniqueFiles {
						f.SetSheetRow("Sheet1", fmt.Sprint("A", r), &[]interface{}{fi.Name, fi.Size, fi.Times})
						r++
					}
					err := f.SaveAs("records.xlsx")
					if err != nil {
						log.Logger.Warn(err)
					}
				}()

				os.Exit(0)
			default:
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(WatchCmd)
	WatchCmd.SetUsageTemplate(watchUsage)

	// pullCmd.Flags().BoolVarP(&pullPublicFlag, "public", "p", false, "pull a public image")
}
