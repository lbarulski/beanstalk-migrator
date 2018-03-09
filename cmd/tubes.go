package cmd

import (
	"github.com/spf13/cobra"
	"beanstalk-migrator/migrator"
	"github.com/kr/beanstalk"
)

var tubesCmd = &cobra.Command{
	Use:   "tubes",
	Short: "List all tubes and show stats",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		sourceConnection, _ := beanstalk.Dial("tcp", srcAddr)
		destinationConnection, _ := beanstalk.Dial("tcp", dstAddr)

		migrator.ShowTubesStats(sourceConnection, destinationConnection)
	},
}

func init() {
	rootCmd.AddCommand(tubesCmd)
}
