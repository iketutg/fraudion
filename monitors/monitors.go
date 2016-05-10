package monitors

// TODO: This will hold the Monitors interface type, if needed!

// Monitor ...
type Monitor interface {
	Run()
}

// DangerousDestinations ...
type DangerousDestinations struct{}

// Run ...
func (monitor *DangerousDestinations) Run() {}
