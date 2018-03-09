package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var srcAddr string
var dstAddr string

var rootCmd = &cobra.Command{
	Use:   "beanstalk-migrator",
	Short: "",
	Long: `Allows to migrate jobs between Beanstalkd instances.
Migrates jobs from all tubes present on source instance to corresponding tubes on destination instance.
Keeps body, pri, delay and ttr of the job.

!!!WARNING!!! Migration of buried jobs is not possible strait into buried queue on destination instance. All buried jobs will be migrated with state "ready" instead of "buried".`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&srcAddr, "source-addr", "", "Source beanstalkd address")
	rootCmd.PersistentFlags().StringVar(&dstAddr, "destination-addr", "", "Destination beanstalkd address")
	rootCmd.MarkPersistentFlagRequired("source-addr")
	rootCmd.MarkPersistentFlagRequired("destination-addr")

}
