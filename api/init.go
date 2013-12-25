package api

var Ops *Operations

// init initialize operations by making a Redis connection
// and connecting to any other outside resources so we are
// aware of failure early on.
func init() {
	Ops = NewOperations()
}
