package system

import (
	"time"
)

const (
	DEBUG   = true
	VERSION = "v0.0.0-alpha.1"
)

// State ...
var State *stateSystem

func init() {
	State = new(stateSystem)
}

type stateSystem struct {
	StartUpTime time.Time
}
