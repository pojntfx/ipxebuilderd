package cmd

const (
	IPXEBuilderdHostPortPortDefault = "0.0.0.0:1010"                    // IPXEBuilderDPortDefault is the default Host:port of `ipxebuilderd`.
	HostPortDocs                    = "Host:port of the server to use." // HostPortDocs is the documentation for the host:port flag.
	ConfigurationFileDocs           = "Configuration file to use."      // ConfigurationFileDocs is the documentation for the configuration file flag.
)

const (
	CouldNotBindFlagsErrorMessage        = "Could not bind flags"         // CouldNotBindFlagsErrorMessage is the error message to throw if binding the flags has failed.
	CouldNotStartRootCommandErrorMessage = "Could not start root command" // CouldNotStartRootCommandErrorMessage is the error message to throw if starting the root command has failed.
)
