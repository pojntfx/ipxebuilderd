package cmd

const (
	keyPrefix         = "ipxe."
	configFileDefault = ""
	serverHostPortKey = keyPrefix + "serverHostPort"
	configFileKey     = keyPrefix + "configFile"
)

var (
	serverHostPortFlag string
	configFileFlag     string
)
