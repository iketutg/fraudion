package softswitch

import (
	"fmt"

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
	Version    string
	CDRsSource CDRsSource
}

// GetCDRsSource ...
func (asterisk *Asterisk) GetCDRsSource() CDRsSource {
	return asterisk.CDRsSource
}

// CDRsSource ...
type CDRsSource interface{}

// CDRsSourceDatabase ...
type CDRsSourceDatabase struct {
	UserName     string
	UserPassword string
	DatabaseName string
	TableName    string
	MysqlOptions string
	connection   *sql.DB
}

// Connect ...
func (cdrSource *CDRsSourceDatabase) Connect() error {

	connections, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?allowOldPasswords=1", cdrSource.UserName, cdrSource.UserPassword, cdrSource.DatabaseName))
	if err != nil {
		return fmt.Errorf("Could not create Database connections")
	}

	cdrSource.connection = connections

	return nil

}
