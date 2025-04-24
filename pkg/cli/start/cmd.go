package start

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/go-picam-streamer/pkg/api"
	"github.com/ylallemant/go-picam-streamer/pkg/cli/start/options"
	"github.com/ylallemant/go-picam-streamer/pkg/globals"
	"github.com/ylallemant/go-picam-streamer/pkg/server"
)

var rootCmd = &cobra.Command{
	Use:   "start",
	Short: "outputs dependency installation path",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		globals.ProcessGlobals()

		serverOptions := &api.ServerOptions{
			Port:    options.Current.Port,
			Address: options.Current.Address,
		}

		cameraOptions := &api.CameraOption{
			CaptureHeight: options.Current.CaptureHeight,
			CaptureWidth:  options.Current.CaptureWidth,
		}

		srv, err := server.New(serverOptions, cameraOptions)
		if err != nil {
			return errors.Wrap(err, "failed to start server")
		}

		return srv.Start()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&options.Current.Address, "address", "a", options.Current.Address, "server listener address")
	rootCmd.PersistentFlags().StringVarP(&options.Current.Port, "port", "p", options.Current.Port, "server listener port")
	rootCmd.PersistentFlags().IntVarP(&options.Current.CaptureHeight, "camera-capture-height", "y", options.Current.CaptureHeight, "camera capture height in pixels")
	rootCmd.PersistentFlags().IntVarP(&options.Current.CaptureWidth, "camera-capture-width", "w", options.Current.CaptureWidth, "camera capture width in pixels")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.FallbackConfig, "fallback-config", globals.Current.FallbackConfig, "if no configuration was found, fallback to the default one")
	//rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.Debug, "debug", globals.Current.Debug, "outputs processing information")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
