package start

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/go-picam-streamer/pkg/cli/start/options"
	"github.com/ylallemant/go-picam-streamer/pkg/globals"
	"github.com/ylallemant/go-picam-streamer/pkg/server"
)

var rootCmd = &cobra.Command{
	Use:   "start",
	Short: "outputs dependency installation path",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, err := server.New()
		if err != nil {
			return errors.Wrap(err, "failed to start server")
		}

		return srv.Start()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&options.Current.Child, "child", options.Current.Child, "forces the path to be relative to the current project")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.FallbackConfig, "fallback-config", globals.Current.FallbackConfig, "if no configuration was found, fallback to the default one")
	//rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.Debug, "debug", globals.Current.Debug, "outputs processing information")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
