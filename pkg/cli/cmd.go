package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ylallemant/go-picam-streamer/pkg/cli/binary/upgrade"
	"github.com/ylallemant/go-picam-streamer/pkg/cli/binary/version"
)

var rootCmd = &cobra.Command{
	Use:          "picam-streamer",
	Short:        "picam-streamer provides a toolset facilitating complex git-hook workflows",
	SilenceUsage: true,
	Long:         ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("please use a subcommand...")
		cmd.Usage()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upgrade.Command())
	rootCmd.AddCommand(version.Command())
}

func Command() *cobra.Command {
	return rootCmd
}
