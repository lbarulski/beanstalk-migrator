package migrator

import (
	"github.com/kr/beanstalk"
	"github.com/olekukonko/tablewriter"
	"os"
	"fmt"
)

func showTableWithTubesStats(c *beanstalk.Conn, caption string) {
	tubesList, err := c.ListTubes()
	if err != nil {
		panic(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetCaption(true, "⇧ "+caption)
	table.SetHeader([]string{"Tube", "Delayed", "Buried", "Ready"})

	for _, tubeName := range tubesList {
		tube := beanstalk.Tube{Conn:c, Name:tubeName}
		tubeStats, _ := tube.Stats()
		delayedCount, _ := tubeStats["current-jobs-delayed"]
		readyCount, _ := tubeStats["current-jobs-ready"]
		buriedCount, _ := tubeStats["current-jobs-buried"]

		table.Append([]string{tubeName, delayedCount, buriedCount, readyCount})
	}
	table.Render()
}

func ShowTubesStats(sourceConnection, destinationConnection *beanstalk.Conn) {
	showTableWithTubesStats(sourceConnection, "Source")

	fmt.Print("\n")

	showTableWithTubesStats(destinationConnection, "Destination")
}