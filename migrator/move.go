package migrator

import (
	"fmt"
	"github.com/kr/beanstalk"
	"strconv"
	"time"
	"reflect"
	"gopkg.in/cheggaaa/pb.v1"
	"strings"
)

func getJob(conn *beanstalk.Conn, tubeName string) (uint64, []byte, error) {
	tube := beanstalk.Tube{conn, tubeName}
	tubeStats, _ := tube.Stats()

	switch true {
	case tubeStats["current-jobs-delayed"] != "0":
		return tube.PeekDelayed()
	case tubeStats["current-jobs-ready"] != "0":
		return tube.PeekReady()
	case tubeStats["current-jobs-buried"] != "0":
		return tube.PeekBuried()
	}

	return 0, nil, JobsNotFoundError{fmt.Sprintf("No jobs found in all queues of tube '%s'", tubeName)}
}

/**
 * wantedTubeNamePart == "" means that you want to move all tubes
 */
func MoveJobs(sourceConnection, destinationConnection *beanstalk.Conn, wantedTubeNamePart string, useProgressBar bool) {
	tubes, err := sourceConnection.ListTubes()
	if err != nil {
		panic(err)
	}

	var tubesBar *pb.ProgressBar
	var jobsBar *pb.ProgressBar
	var barPool *pb.Pool
	if (useProgressBar) {
		if "" != wantedTubeNamePart {
			var wantedTubes []string;
			for _, tubeName := range tubes {
				if strings.Contains(tubeName, wantedTubeNamePart) {
					wantedTubes = append(wantedTubes, tubeName);
				}
			}

			tubes = wantedTubes
		}

		tubesCount := len(tubes)
		tubesBar = pb.New(tubesCount).Prefix("Tubes")
		jobsBar = pb.New(0).Prefix("Jobs")

		barPool, _ = pb.StartPool(tubesBar, jobsBar)
	}

	for _, tubeName := range tubes {
		if (useProgressBar) {
			tube := beanstalk.Tube{sourceConnection, tubeName}
			tubeStats, _ := tube.Stats()

			delayedCount, _ := strconv.Atoi(tubeStats["current-jobs-delayed"])
			readyCount, _ := strconv.Atoi(tubeStats["current-jobs-ready"])
			buriedCount, _ := strconv.Atoi(tubeStats["current-jobs-buried"])

			jobsBar.Total = int64(delayedCount + readyCount + buriedCount)
			tubesBar.Prefix(fmt.Sprintf("Tubes [%s]", tubeName))
		}

		delayedMovedCounter := int64(0)
		buriedMovedCounter := int64(0)
		readyMovedCounter := int64(0)

		for {
			id, body, err := getJob(sourceConnection, tubeName)
			if err != nil {
				if (reflect.TypeOf(err).String() == reflect.TypeOf(JobsNotFoundError{}).String()) {
					break;
				}

				panic(err)
			}

			jobStats, err := sourceConnection.StatsJob(id)
			if err != nil {
				panic(err)
			}

			pri, _ := strconv.ParseUint(jobStats["pri"], 10, 32)
			delay, _ := strconv.ParseInt(jobStats["time-left"], 10, 64)
			ttr, _ := strconv.ParseInt(jobStats["ttr"], 10, 32)

			tube := beanstalk.Tube{Conn: destinationConnection, Name: tubeName}
			_, err = tube.Put(body, uint32(pri), time.Duration(delay)*time.Second, time.Duration(ttr)*time.Second)
			if err != nil {
				panic(err)
			}

			sourceConnection.Delete(id)

			if (useProgressBar) {
				switch jobStats["state"] {
				case "delayed":
					delayedMovedCounter++
					break
				case "buried":
					buriedMovedCounter++
					break
				case "ready":
					readyMovedCounter++
					break
				}

				jobsBar.Increment()
				jobsBar.Prefix(fmt.Sprintf("Jobs [d: %d, b: %d, r: %d]", delayedMovedCounter, buriedMovedCounter, readyMovedCounter))
			}
		}

		if (useProgressBar) {
			jobsBar.Finish()
			tubesBar.Increment()
		}
	}
	if (useProgressBar) {
		tubesBar.Finish()
		barPool.Stop()
	}
}
