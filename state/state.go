package state

import (
	"time"
)

// System ...
var System *system

// Triggers ...
var Triggers *triggers

func init() {
	System = new(system)
	Triggers = new(triggers)
}

type system struct {
	StartUpTime time.Time
}

// StateTriggers ...
type triggers struct {
	StateSimultaneousCalls     simultaneousCalls
	StateDangerousDestinations dangerousDestinations
	StateExpectedDestinations  expectedDestinations
	StateSmallDurationCalls    smallCallDurations
}

type simultaneousCalls struct {
	LastActionChainRunTime time.Time
	ActionChainRunCount    uint32
}

type dangerousDestinations struct {
	LastActionChainRunTime time.Time
	ActionChainRunCount    uint32
}

type expectedDestinations struct {
	LastActionChainRunTime time.Time
	ActionChainRunCount    uint32
}

type smallCallDurations struct {
	LastActionChainRunTime time.Time
	ActionChainRunCount    uint32
}
