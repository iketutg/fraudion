package softswitches

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"database/sql"
	"os/exec"

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

const (
	asteriskDialString = "(?:SIP|DAHDI)/[^@&]+/([0-9]+)" // NOTE: Currently supported dial string format
)

// Monitored ...
var Monitored Softswitch

// Softswitch ...
type Softswitch interface {
	GetCDRsSource() CDRsSource
	GetHits(func(string) (string, bool, error), time.Duration) (map[string]*Hits, error)
	GetCurrentActiveCalls(uint32) (uint32, error)
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

			matchesDialString := regexp.MustCompile(asteriskDialString)
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

// GetCurrentActiveCalls ...
func (asterisk *Asterisk) GetCurrentActiveCalls(minimumNumberLength uint32) (uint32, error) {

	// TODO: Make this depend on the Asterisk version because command format and result parsing may vary!

	command := exec.Command("asterisk", "-rx 'core show channels concise'")
	//command := exec.Command("cat", "asterisk.output") // NOTE: This is a just a test code to test Asterisk output without a local Asterisk via the non-commited text file: asterisk.output where one can put example output

	output, err := command.Output()
	if err != nil {
		return 0, err
	}

	numberOfCalls := 0

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {

		lineItems := strings.Split(scanner.Text(), "!")

		matchesDialString := regexp.MustCompile(asteriskDialString)
		matchedString := matchesDialString.FindString(lineItems[6])

		if lineItems[5] == "Dial" || matchedString != "" { // NOTE: Ignore if "lastapp" is not Dial and "lastdata" does not contain an expected dial string

			dialedNumber := matchesDialString.FindStringSubmatch(lineItems[6])[1]

			if uint32(len(dialedNumber)) > minimumNumberLength {
				numberOfCalls++
			}
		}

	}

	return uint32(numberOfCalls), nil

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
