package cmd

import (
	"github.com/spf13/cobra"
	"beanstalk-migrator/migrator"
	"github.com/kr/beanstalk"
)

var quietFlag bool;

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move jobs between Beanstalkd instances",
	Long: `Allows to migrate jobs between Beanstalkd instances.
Migrates jobs from all tubes present on source instance to corresponding tubes on destination instance.
Keeps body, pri, delay and ttr of the job.

!!!WARNING!!! Migration of buried jobs is not possible strait into buried queue on destination instance. All buried jobs will be migrated with state "ready" instead of "buried".`,
	Run: func(cmd *cobra.Command, args []string) {
		sourceConnection, _ := beanstalk.Dial("tcp", srcAddr)
		destinationConnection, _ := beanstalk.Dial("tcp", dstAddr)

		migrator.MoveJobs(sourceConnection, destinationConnection,!quietFlag)
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)
	moveCmd.Flags().BoolVarP(&quietFlag, "quiet", "q", false, "Don't show progress bar")
}
