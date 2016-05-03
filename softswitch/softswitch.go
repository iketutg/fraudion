package softswitch

import (
	"database/sql"
)

// Monitored ...
var Monitored Softswitch

// Softswitch ...
type Softswitch interface {
	GetCDRsSource() CDRsSource
}

// Asterisk ...
type Asterisk struct {
	CDRsSource CDRsSource
}

// CDRsSource ...
type CDRsSource interface{}

// CDRsSourceDatabase ...
type CDRsSourceDatabase struct {
	Connection *sql.DB
	TableName  string
}

// Init ...
func Init() error {

	// TODO: This will require config to be parsed, validated and loaded

	// TODO: Check what Softswitch type it is

	// TODO: Check what CDRs Source it is

	dbConnections, err := sql.Open("mysql", "root:@tcp(localhost:3306)/asteriskcdrdb?allowOldPasswords=1")
	if err != nil {

	}

	sourceDatabase := new(CDRsSourceDatabase)
	sourceDatabase.Connection = dbConnections
	sourceDatabase.TableName = "cdr"

	newMonitored := new(Asterisk)
	newMonitored.CDRsSource = sourceDatabase

	Monitored = newMonitored

	return nil

}

// GetCDRsSource ...
func (asterisk *Asterisk) GetCDRsSource() CDRsSource {
	return asterisk.CDRsSource
}
