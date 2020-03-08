package cmd

const (
	keyPrefix         = "ipxe."
	configFileDefault = ""
	serverHostPortKey = keyPrefix + "serverHostPort"
	configFileKey     = keyPrefix + "configFile"
	platformKey       = keyPrefix + "platform"
	driverKey         = keyPrefix + "driver"
	extensionKey      = keyPrefix + "extension"
	scriptKey         = keyPrefix + "script"
)

var (
	serverHostPortFlag string
	configFileFlag     string
)
