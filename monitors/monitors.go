package monitors

import (
	"time"

	"github.com/andmar/fraudion/config"
	"github.com/andmar/fraudion/softswitches"
)

const (
	// RunModeNormal ...
	RunModeNormal = iota
	// RunModeInWarning ...
	RunModeInWarning
	// RunModeInAlarm ...
	RunModeInAlarm
)

// Monitor ...
type Monitor interface {
	Run()
}

// monitorBase ...
type monitorBase struct {
	Softswitch softswitches.Softswitch
}

// DangerousDestinations ...
type DangerousDestinations struct {
	monitorBase
	Config *config.MonitorDangerousDestinations
	State  StateDangerousDestinations
}

// SimultaneousCalls ...
type SimultaneousCalls struct {
	monitorBase
	Config *config.MonitorSimultaneousCalls
	State  StateSimultaneousCalls
}

type stateBase struct {
	LastActionChainRunTime time.Time
	ActionChainRunCount    uint32
	RunMode                int
}

// StateDangerousDestinations ...
type StateDangerousDestinations struct {
	stateBase
}

// StateSimultaneousCalls ...
type StateSimultaneousCalls struct {
	stateBase
}
