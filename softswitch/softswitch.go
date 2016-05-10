package softswitch

import (
	"fmt"

	"database/sql"
)

const (
	// SystemAsterisk ...
	SystemAsterisk = "*asterisk"
	// CDRSourceDatabase ...
	CDRSourceDatabase = "*database"
)

// Monitored ...
var Monitored Softswitch

// Softswitch ...
type Softswitch interface {
	GetCDRsSource() CDRsSource
}

// Asterisk ...
type Asterisk struct {
	Version    string
	CDRsSource CDRsSource
}

// GetCDRsSource ...
func (asterisk *Asterisk) GetCDRsSource() CDRsSource {
	return asterisk.CDRsSource
}

// CDRsSource ...
type CDRsSource interface {
	GetCDRs()
}

// CDRsSourceDatabase ...
type CDRsSourceDatabase struct {
	Type         string
	DBMS         string
	UserName     string
	UserPassword string
	DatabaseName string
	TableName    string
	connections  *sql.DB
}

// GetCDRs ...
func (cdrSource *CDRsSourceDatabase) GetCDRs() {
}

// Connect ...
func (cdrSource *CDRsSourceDatabase) Connect() error {

	connections, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?allowOldPasswords=1", cdrSource.UserName, cdrSource.UserPassword, cdrSource.DatabaseName))
	if err != nil {
		return fmt.Errorf("Could not create Database connections")
	}

	cdrSource.connections = connections

	return nil

}

// GetConnections ...
func (cdrSource *CDRsSourceDatabase) GetConnections() *sql.DB {

	return cdrSource.connections

}
