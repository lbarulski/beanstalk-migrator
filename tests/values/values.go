package values

import "time"

type JobArgs struct {
	Body      []byte
	Pri       uint32
	Delay     time.Duration
	Ttr       time.Duration
	IsToBurry bool
}

// jobs order is important
// first we are looking for quques in order: buried, delayed, ready
// inside specific queue (delayed, buried, ready) all jobs needs to be putted in order with which they will be peeked from beanstalkd
var Values = [...]JobArgs{
	//BURIED
	JobArgs{Body: []byte("test buried 1"), Pri: 1, Delay: 0, Ttr: time.Duration(60*time.Second), IsToBurry: true},
	// DELAYED
	JobArgs{Body: []byte("test delayed 1"), Pri: 128, Delay: time.Duration(60*time.Minute), Ttr: time.Duration(60*time.Second), IsToBurry: false},
	// READY
	JobArgs{Body: []byte("test ready 1"), Pri: 128, Delay: 0, Ttr: time.Duration(1*time.Second), IsToBurry: false},
	JobArgs{Body: []byte("test ready 2"), Pri: 1024, Delay: 0, Ttr: time.Duration(60*time.Second), IsToBurry: false},
}
