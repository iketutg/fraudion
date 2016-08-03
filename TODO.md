# TODO
 * Add some missing string validations in the config schema (phone numbers, e-mails, URLs)
 * Review logging, which stamps to use where and what to output to each one of them
 * Make active logging Stamps configurable via command line arguments

 * (?) Fail main thread when one of the monitors fails? OR what to do when this happens?
 * (?) See how to add configurations specific to the code in a release, to that release's "package", godoc.org does not seem to support it [This is done using consts because compiler optimizes out stuff that depend on consts with value false]
