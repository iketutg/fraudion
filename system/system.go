package system

import (
	"time"
)

// State ...
var State *stateSystem

func init() {
	State = new(stateSystem)
}

type stateSystem struct {
	StartUpTime time.Time
}
