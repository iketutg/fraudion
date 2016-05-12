package softswitches

import (
	"fmt"
	"regexp"
	"time"

	"database/sql"

	"github.com/andmar/marlog"
)

const (
	// TypeAsterisk ...
	TypeAsterisk = "*asterisk"
	// TypeFreeSwitch ...
	TypeFreeSwitch = "*freeswitch"
	// CDRSourceDatabase ...
	CDRSourceDatabase = "*database"
)

// Monitored ...
var Monitored Softswitch

// Softswitch ...
type Softswitch interface {
	GetCDRsSource() CDRsSource
	GetHits(func(string) (string, bool, error), time.Duration) (map[string]*Hits, error)
}

// Asterisk ...
type Asterisk struct {
	Version    string
	CDRsSource CDRsSource
}

// GetHits ...
func (asterisk *Asterisk) GetHits(matches func(string) (string, bool, error), considerCDRsFromLast time.Duration) (map[string]*Hits, error) { // TODO: Change this to accept the monitor that is monitoring the Softswitch?

	log := marlog.MarLog

	switch asterisk.CDRsSource.(type) {
	case *CDRsSourceDatabase:

		cdrsSource, ok := asterisk.CDRsSource.(*CDRsSourceDatabase)
		if ok == false {
			return nil, fmt.Errorf("could not convert CDRs Source to the appropriate type")
		}

		if err := cdrsSource.GetConnections().Ping(); err != nil {
			return nil, err
		}

		log.LogS("DEBUG", "Database connection is A-Ok!")

		rows, err := cdrsSource.GetConnections().Query(fmt.Sprintf("SELECT * FROM cdr WHERE calldate >= DATE_SUB(CURDATE(), INTERVAL %v HOUR) ORDER BY calldate DESC;", uint32(considerCDRsFromLast.Hours()))) // TODO: The Query format should depend on DBMS?
		if err != nil {
			return nil, err
		}

		result := make(map[string]*Hits)

		for rows.Next() {

			var calldate, clid, src, dst, dcontext, channel, dstchannel, lastapp, lastdata, disposition, accountcode, uniqueid, userfield string
			var duration, billsec, amaflags uint32

			err := rows.Scan(&calldate, &clid, &src, &dst, &dcontext, &channel, &dstchannel, &lastapp, &lastdata, &duration, &billsec, &disposition, &amaflags, &accountcode, &uniqueid, &userfield)
			if err != nil {
				return nil, err
			}

			//fmt.Println(calldate, clid, src, dst, dcontext, channel, dstchannel, lastapp, lastdata, duration, billsec, disposition, amaflags, accountcode, uniqueid, userfield)

			matchesDialString := regexp.MustCompile("(?:SIP|DAHDI)/[^@&]+/([0-9]+)") // NOTE: Currently supported dial string format
			matchedString := matchesDialString.FindString(lastdata)
			if lastapp != "Dial" || matchedString == "" { // NOTE: Ignore if "lastapp" is not Dial and "lastdata" does not contain an expected dial string
				continue
			}

			dialedNumber := matchesDialString.FindStringSubmatch(lastdata)[1]

			prefix, matched, err := matches(dialedNumber)
			if err != nil {
				return nil, err
			}
			if matched == true {

				if _, found := result[prefix]; found != true { // NOTE: If the prefix doesn't have matches already, create a new hits object ELSE add to the count for that prefix
					result[prefix] = new(Hits)
				}
				result[prefix].NumberOfHits++
				result[prefix].Destinations = append(result[prefix].Destinations, dialedNumber)

			}

		}

		return result, nil

	default:
		return nil, fmt.Errorf("unknown CDRs Source object type)")
	}

}

// GetCDRsSource ...
func (asterisk *Asterisk) GetCDRsSource() CDRsSource {
	return asterisk.CDRsSource
}

// CDRsSource ...
type CDRsSource interface {
	//GetCDRs()
}

// CDRsSourceDatabase ...
type CDRsSourceDatabase struct {
	DBMS         string
	UserName     string
	UserPassword string
	DatabaseName string
	TableName    string
	connections  *sql.DB
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
func (cdrSource CDRsSourceDatabase) GetConnections() *sql.DB {

	return cdrSource.connections

}

// Hits ...
type Hits struct {
	NumberOfHits uint32
	Destinations []string
}
