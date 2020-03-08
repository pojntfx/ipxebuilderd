package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/cheggaaa/pb/v3"
	constants "github.com/pojntfx/ipxebuilderd/cmd"
	globalConstants "github.com/pojntfx/ipxebuilderd/pkg/constants"
	iPXEBuilder "github.com/pojntfx/ipxebuilderd/pkg/proto/generated"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/bloom42/libs/rz-go"
	"gitlab.com/bloom42/libs/rz-go/log"
	"google.golang.org/grpc"
)

var applyCmd = &cobra.Command{
	Use:     "apply",
	Aliases: []string{"a"},
	Short:   "Apply a ipxe",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !(viper.GetString(configFileKey) == configFileDefault) {
			viper.SetConfigFile(viper.GetString(configFileKey))

			if err := viper.ReadInConfig(); err != nil {
				return err
			}
		}

		conn, err := grpc.Dial(viper.GetString(serverHostPortKey), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return err
		}
		defer conn.Close()

		client := iPXEBuilder.NewIPXEBuilderClient(conn)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		stream, err := client.Create(ctx, &iPXEBuilder.IPXE{
			Platform:  viper.GetString(platformKey),
			Driver:    viper.GetString(driverKey),
			Extension: viper.GetString(extensionKey),
			Script:    viper.GetString(scriptKey),
		})
		if err != nil {
			return err
		}

		bar := pb.Full.Start(globalConstants.TotalCompileSteps)

		waitChan := make(chan struct{})
		var id string
		go func() {
			for {
				update, err := stream.Recv()
				if err == io.EOF {
					close(waitChan)

					bar.Finish()

					return
				}
				if err != nil {
					bar.Finish()

					log.Fatal("Could not receive status update", rz.Err(err))
				}

				bar.Increment()

				id = update.GetId()
			}
		}()

		<-waitChan

		fmt.Printf("iPXE \"%s\" created\n", id)

		return nil
	},
}

func init() {
	var (
		platformFlag, driverFlag, extensionFlag string
		scriptFlag                              string
	)

	applyCmd.PersistentFlags().StringVarP(&serverHostPortFlag, serverHostPortKey, "s", constants.IPXEBuilderdHostPortPortDefault, constants.HostPortDocs)
	applyCmd.PersistentFlags().StringVarP(&configFileFlag, configFileKey, "f", configFileDefault, constants.ConfigurationFileDocs)
	applyCmd.PersistentFlags().StringVarP(&platformFlag, platformKey, "p", "bin-x86_64-efi", "Platform to build the iPXE for.")
	applyCmd.PersistentFlags().StringVarP(&driverFlag, driverKey, "d", "ipxe", "Driver to build the iPXE with.")
	applyCmd.PersistentFlags().StringVarP(&extensionFlag, extensionKey, "e", "efi", "Extension build the iPXE for.")
	applyCmd.PersistentFlags().StringVarP(&scriptFlag, scriptKey, "a", `#!ipxe
autoboot`, "Script to embed in the iPXE.")

	if err := viper.BindPFlags(applyCmd.PersistentFlags()); err != nil {
		log.Fatal(constants.CouldNotBindFlagsErrorMessage, rz.Err(err))
	}

	viper.AutomaticEnv()

	rootCmd.AddCommand(applyCmd)
}
