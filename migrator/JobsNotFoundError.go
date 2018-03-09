package migrator

type JobsNotFoundError struct {
	msg string
}

func (err JobsNotFoundError) Error() string {
	return err.msg
}