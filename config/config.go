package config

const (
	ConstOriginURL = iota
	ConstOriginFile
)

const (
	ConstDefaultOrigin         = ConstOriginFile
	ConstDefaultConfigDir      = "."
	ConstDefaultConfigFilename = "fraudion.json"
	ConstDefaultConfigURL      = "http://mirrors.voipit.pt/fraudion.json"
)
