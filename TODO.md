# TODO
 * Add some string validations in the config schema (phone numbers, e-mails, URLs)
 * Build system. Build on the Mac to run on CentOS 5: sudo GOOS=linux GOARCH=386 GOPATH=/Users/andre/Dev/go go install github.com/andmar/fraudion
 * Add a blank line followed by a start time message when the system starts logging to a file
 * See which information we want to have in each debug stamp, INFO/DEBUG in Monitores are not really correct as they are now
 * Concept of exitChan that makes it so that Monitors can trigger the end of main process

 * [DONE?] Consider if validation should also check dependencies in information that currently is done only on the config loading phase
 * [This is done using consts because compiler optimizes out stuff that depend on false consts] See how to add configuration specific to the code in a release to that release "package", godoc.org does not seem to support it
 * [This will delay log start] Make logging Stamps configurable via the General section of the config file OR via command line arguments which seems to be the best option because config is parsed after the logging is initiated
