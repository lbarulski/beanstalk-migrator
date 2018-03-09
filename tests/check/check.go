package main

import (
	"github.com/kr/beanstalk"
	"beanstalk-migrator/tests/values"
	"fmt"
	"strconv"
	"time"
)

func main() {
	c, err := beanstalk.Dial("tcp", "127.0.0.1:11302")
	if err != nil {
		panic(err)
	}

	tube := beanstalk.Tube{Conn: c, Name: "test"}

	for _, e := range values.Values {
		id, body, err := getJob(tube, e)
		if (err != nil) {
			panic(err)
		}

		jobStats, _ := c.StatsJob(id)

		panicIfFalse(string(e.Body) == string(body), "Body not equal", jobStats, body, e)
		panicIfFalse(jobStats["pri"] == strconv.FormatUint(uint64(e.Pri), 10), "Priority is not equal", jobStats, body, e)
		panicIfFalse(jobStats["ttr"] == strconv.Itoa(int(e.Ttr/time.Second)), "Ttr is not equal", jobStats, body, e)

		if (e.Delay > 0) {
			panicIfFalse(jobStats["state"] == "delayed", "Job is delayed", jobStats, body, e)

			jobStatsDelay, _ := strconv.Atoi(jobStats["delay"])
			panicIfFalse(jobStatsDelay < int(e.Delay/time.Second), "Delay is smaller than in base job", jobStats, body, e)
		}

		if (e.IsToBurry) {
			panicIfFalse(jobStats["state"] == "ready", "Job is ready", jobStats, body, e)
		}

		fmt.Println(fmt.Sprintf("[OK] %s", string(e.Body)))
		c.Delete(id)
	}
}

func panicIfFalse(isOk bool, test string, jobStats map[string]string, body []byte, e values.JobArgs) {
	if (!isOk) {
		fmt.Println(jobStats, string(body), string(e.Body), e.Pri, e.Delay, e.Ttr, e.IsToBurry)
		panic(test);
	}
}
func getJob(tube beanstalk.Tube, e values.JobArgs) (id uint64, body []byte, err error) {
	switch true {
	case e.IsToBurry:
		return tube.PeekReady()
	case e.Delay > 0:
		return tube.PeekDelayed()
	default:
		return tube.PeekReady()
	}

	panic(fmt.Sprintf("Job for '%s' not found", e.Body))
}
