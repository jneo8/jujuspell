package data

type JobStatus int

const (
	Pending JobStatus = iota
	Running
	Canceling
)
