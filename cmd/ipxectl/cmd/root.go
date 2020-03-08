package cmd

import (
	"strings"

	constants "github.com/pojntfx/ipxebuilderd/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/bloom42/libs/rz-go"
	"gitlab.com/bloom42/libs/rz-go/log"
)

var rootCmd = &cobra.Command{
	Use:   "ipxectl",
	Short: "ipxectl manages ipxebuilderd, the iPXE build daemon.",
	Long: `ipxectl manages ipxebuilderd, the iPXE build daemon.

Find more information at:
https://pojntfx.github.io/ipxe/`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viper.SetEnvPrefix("ipxe")
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	},
}

// Execute starts the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(constants.CouldNotStartRootCommandErrorMessage, rz.Err(err))
	}
}
