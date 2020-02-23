package main

import (
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	constants "github.com/pojntfx/ipxebuilderd/cmd"
	iPXEBuilder "github.com/pojntfx/ipxebuilderd/pkg/proto/generated"
	"github.com/pojntfx/ipxebuilderd/pkg/svc"
	"github.com/pojntfx/ipxebuilderd/pkg/workers"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/bloom42/libs/rz-go"
	"gitlab.com/bloom42/libs/rz-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	keyPrefix         = "ipxebuilderd."
	configFileDefault = ""
	configFileKey     = keyPrefix + "configFile"
	listenHostPortKey = keyPrefix + "listenHostPort"
)

var rootCmd = &cobra.Command{
	Use:   "ipxebuilderd",
	Short: "ipxebuilderd is an iPXE build daemon",
	Long: `ipxebuilderd is an iPXE build daemon.

Find more information at:
https://pojntfx.github.io/ipxebuilderd/`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viper.SetEnvPrefix("ipxebuilderd")
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if !(viper.GetString(configFileKey) == configFileDefault) {
			viper.SetConfigFile(viper.GetString(configFileKey))

			if err := viper.ReadInConfig(); err != nil {
				return err
			}
		}
		builder := workers.Builder{
			BasePath: filepath.Join(os.TempDir(), "ipxebuilderd", uuid.NewV4().String()),
		}

		listener, err := net.Listen("tcp", viper.GetString(listenHostPortKey))
		if err != nil {
			return err
		}

		server := grpc.NewServer()
		reflection.Register(server)

		iPXEService := svc.IPXEBuilder{
			Builder: &builder,
		}

		iPXEBuilder.RegisterIPXEBuilderServer(server, &iPXEService)

		interrupt := make(chan os.Signal, 2)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-interrupt

			// Allow manually killing the process
			go func() {
				<-interrupt

				os.Exit(1)
			}()

			log.Info("Gracefully stopping server (this might take a few seconds)")

			server.GracefulStop()
		}()

		if err := iPXEService.Extract(); err != nil {
			return err
		}

		log.Info("Starting server")

		if err := server.Serve(listener); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	var (
		configFileFlag string
		hostPortFlag   string
	)

	rootCmd.PersistentFlags().StringVarP(&configFileFlag, configFileKey, "f", configFileDefault, constants.ConfigurationFileDocs)
	rootCmd.PersistentFlags().StringVarP(&hostPortFlag, listenHostPortKey, "l", constants.IPXEBuilderDHostPortDefault, "TCP listen host:port.")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatal(constants.CouldNotBindFlagsErrorMessage, rz.Err(err))
	}

	viper.AutomaticEnv()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(constants.CouldNotStartRootCommandErrorMessage, rz.Err(err))
	}
}
