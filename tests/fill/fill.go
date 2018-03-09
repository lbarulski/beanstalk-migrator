package main

import (
	"github.com/kr/beanstalk"
	"beanstalk-migrator/tests/values"
	"fmt"
	"time"
)

func main()  {
	c, err := beanstalk.Dial("tcp", "127.0.0.1:11301")
	if err != nil {
		panic(err)
	}

	tube := beanstalk.Tube{Conn: c, Name: "test"}
	tubeSet := beanstalk.TubeSet{Conn:c, Name: map[string]bool{ tube.Name: true }}

	for _, element := range values.Values {
		id, err := tube.Put(element.Body, element.Pri, element.Delay, element.Ttr)
		if (err != nil) {
			panic(err)
		}

		if (element.IsToBurry) {

			reservedId, _, _ := tubeSet.Reserve(time.Duration(60 * time.Minute))

			if (id != reservedId) {
				panic(fmt.Sprintf("Ids are not equal: %d|%d", id, reservedId))
			}

			err := c.Bury(id, element.Pri)
			if (err != nil) {
				jobStats, _ := c.StatsJob(id)
				fmt.Println(jobStats)

				panic(err)
			}
		}
	}
}