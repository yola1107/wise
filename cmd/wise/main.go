package wise

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "wise",
	Short:   "wise: An elegant toolkit for Go microservices.",
	Long:    `wise: An elegant toolkit for Go microservices.`,
	Version: release,
}

func init() {
	//rootCmd.AddCommand(project.CmdNew)
	//rootCmd.AddCommand(proto.CmdProto)
	//rootCmd.AddCommand(upgrade.CmdUpgrade)
	//rootCmd.AddCommand(change.CmdChange)
	//rootCmd.AddCommand(run.CmdRun)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
