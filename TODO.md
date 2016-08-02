
# TODO

 * Add some string validations in the config schema (phone numbers, e-mails, URLs)
 * See which information we want to have in each debug stamp, INFO/DEBUG/VERBOSE in Monitores are not really correct as they are now
 * Concept of exitChan that makes it so that Monitors can trigger the end of main process

 * [This is done using consts because compiler optimizes out stuff that depend on false consts] See how to add configuration specific to the code in a release to that release "package", godoc.org does not seem to support it
 * [This will delay log start] Make logging Stamps configurable via the General section of the config file OR via command line arguments which seems to be the best option because config is parsed after the logging is initiated
 * [DONE?] Consider if validation should also check dependencies in information that currently is done only on the config loading phase
 * [DONE?] Add capability to download config JSON from URL (centralized configuration? provisioning server by hostname?)
