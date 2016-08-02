package config

const (
	ConstOriginURL  = "url"
	ConstOriginFile = "file"
)

const (
	ConstDefaultOrigin         = ConstOriginFile
	ConstDefaultConfigDir      = "."
	ConstDefaultConfigFilename = "fraudion.json"
	// TODO This "DefaultConfigURL" should not have a default right?
	ConstDefaultConfigURL = "http://mirrors.voipit.pt/fraudion.json"
)
