package softswitches

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"database/sql"
	"os/exec"

	"github.com/andmar/fraudion/system"

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
	// NOTE: Currently supported dial string format for Asterisk
	asteriskDialString = "(?:SIP|DAHDI)/[^@&]+/([0-9]+)"
)

// Monitored ...
var Monitored Softswitch

// Softswitch ...
type Softswitch interface {
	GetCDRsSource() CDRsSource
	GetHits(func(string, ...uint32) (string, bool, error), time.Duration, bool) (map[string]*Hits, error)
	GetCurrentActiveCalls(uint32) (uint32, error)
}

// Asterisk ...
type Asterisk struct {
	Version    string
	CDRsSource CDRsSource
}

// GetHits ...
func (asterisk *Asterisk) GetHits(matches func(string, ...uint32) (string, bool, error), considerCDRsFromLast time.Duration, considerCallDuration bool) (map[string]*Hits, error) {

	log := marlog.MarLog

	switch asterisk.CDRsSource.(type) {
	case *CDRsSourceDatabase:

		log.LogS("DEBUG", "CDRs Source is Database")

		cdrsSource, ok := asterisk.CDRsSource.(*CDRsSourceDatabase)
		if ok == false {
			return nil, fmt.Errorf("could not convert CDRs Source to the appropriate type")
		}

		if err := cdrsSource.GetConnections().Ping(); err != nil {
			return nil, err
		}

		log.LogS("DEBUG", "Database connection is A-Ok!")

		// TODO: The Query format should depend on DBMS?
		// TODO: This query should consider CDRs only from StartUpTime onwards. This is currenctly the only though
		// of way of resetting the system state to Normal, if we don't do this, when we restart the system will
		// still be in alarm and will be until considerCDRsFromLast.Hours passes...
		var rows *sql.Rows
		var err error
		if system.DEBUG {
			rows, err = cdrsSource.GetConnections().Query(fmt.Sprintf("SELECT * FROM cdr WHERE calldate >= DATE_SUB(CURDATE(), INTERVAL %v HOUR) ORDER BY calldate DESC;", uint32(considerCDRsFromLast.Hours())))
		} else {
			sinceStartUp := time.Since(system.State.StartUpTime)
			rows, err = cdrsSource.GetConnections().Query(fmt.Sprintf("SELECT * FROM cdr WHERE calldate >= DATE_SUB(CURDATE(), INTERVAL %v HOUR) AND calldate >= DATE_SUB(CURDATE(), INTERVAL %v HOUR) ORDER BY calldate DESC;", uint32(considerCDRsFromLast.Hours()), sinceStartUp.Hours()))
		}

		if err != nil {
			log.LogS("ERROR", "could not query the database")
			return nil, err
		}

		result := make(map[string]*Hits)

		numberOfCDRsTotal := 0
		numberOfCDRsSuitable := 0
		numberOfCDRsMatched := 0

		for rows.Next() {

			var calldate, clid, src, dst, dcontext, channel, dstchannel, lastapp, lastdata, disposition, accountcode, uniqueid, userfield string
			var duration, billsec, amaflags uint32

			err := rows.Scan(&calldate, &clid, &src, &dst, &dcontext, &channel, &dstchannel, &lastapp, &lastdata, &duration, &billsec, &disposition, &amaflags, &accountcode, &uniqueid, &userfield)
			if err != nil {
				log.LogS("ERROR", "Could not bring query results to variables")
				return nil, err
			}

			numberOfCDRsTotal++

			//fmt.Println(calldate, clid, src, dst, dcontext, channel, dstchannel, lastapp, lastdata, duration, billsec, disposition, amaflags, accountcode, uniqueid, userfield)

			matchesDialString := regexp.MustCompile(asteriskDialString)
			matchedString := matchesDialString.FindString(lastdata)
			// NOTE: Ignore if "lastapp" is not Dial and "lastdata" does not contain an expected dial string
			if lastapp != "Dial" || matchedString == "" {
				continue
			}

			numberOfCDRsSuitable++

			dialedNumber := matchesDialString.FindStringSubmatch(lastdata)[1]

			var prefix string
			var matched bool
			if !considerCallDuration {
				prefix, matched, err = matches(dialedNumber)
			} else {
				prefix, matched, err = matches(dialedNumber, billsec)
			}
			if err != nil {
				log.LogS("ERROR", "Number not suitable")
				return nil, err
			}

			if matched == true {

				numberOfCDRsMatched++

				// NOTE: If the prefix doesn't have matches already, create a new hits object ELSE add to the count for that prefix
				if _, found := result[prefix]; found != true {
					result[prefix] = new(Hits)
					result[prefix].Prefix = prefix
				}
				result[prefix].NumberOfHits++
				result[prefix].Destinations = append(result[prefix].Destinations, dialedNumber)

			}

		}

		log.LogS("INFO", "Results: Suitable: "+strconv.Itoa(numberOfCDRsSuitable)+", Matched: "+strconv.Itoa(numberOfCDRsMatched)+", Total: "+strconv.Itoa(numberOfCDRsTotal))

		return result, nil

	default:
		return nil, fmt.Errorf("unknown CDRs Source object type)")
	}

}

// GetCurrentActiveCalls ...
func (asterisk *Asterisk) GetCurrentActiveCalls(minimumNumberLength uint32) (uint32, error) {

	log := marlog.MarLog

	// TODO: Make this depend on the Asterisk version because command format and result parsing may vary!

	//command := exec.Command("asterisk", "-rx", "core show channels concise") // NOTE: Fraudion has to have the permission to do this...
	command := exec.Command("cat", "asterisk.output") // NOTE: This is a just a test code to test Asterisk output without a local Asterisk via the non-commited text file: asterisk.output where one can put example output

	output, err := command.Output()
	if err != nil {
		log.LogS("ERROR: ", err.Error())
		return 0, err
	}

	numberOfCalls := 0
	numberOfLines := 0

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {

		numberOfLines++

		lineItems := strings.Split(scanner.Text(), "!")

		if len(lineItems) == 14 {

			matchesDialString := regexp.MustCompile(asteriskDialString)
			matchedString := matchesDialString.FindString(lineItems[6])

			// NOTE: Ignore if "lastapp" is not Dial and "lastdata" does not contain an expected dial string
			if lineItems[5] == "Dial" || matchedString != "" {

				dialedNumber := matchesDialString.FindStringSubmatch(lineItems[6])[1]

				if uint32(len(dialedNumber)) > minimumNumberLength {
					numberOfCalls++
				} else {
					log.LogS("DEBUG", "Number \""+dialedNumber+"\" is ignored due to length")
				}

			}

		} else {
			log.LogS("DEBUG", "Line has weird item count: "+strconv.Itoa(len(lineItems)))
		}

	}

	log.LogS("DEBUG", "Analized "+strconv.Itoa(numberOfLines)+" lines and found "+strconv.Itoa(numberOfCalls)+" suitable calls")

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
		return fmt.Errorf("could not create Database connections (" + err.Error() + ")")
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
	Prefix       string
	NumberOfHits uint32
	Destinations []string
}
