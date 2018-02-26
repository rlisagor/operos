package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/paxautoma/operos/components/operosup/pkg/client"
	"github.com/paxautoma/operos/components/operosup/pkg/steps"
)

func NewUpgradeCommand(dialer *client.Dialer) *cobra.Command {
	var flags steps.UpgradeFlags

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Perform upgrade",
		Long:  `If an upgrade is available, download and apply it`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if flags.UpgradeOnly && flags.UpgradeFile == "" {
				return errors.New("--upgrade-only requires --file to be set")
			}

			if flags.UpgradeOnly && flags.DownloadOnly {
				return errors.New("cannot set both --upgrade-only and --download-only")
			}

			return steps.DoUpgrade(dialer, &flags)
		},
	}

	cmd.Flags().StringVar(
		&flags.RootPath, "root", "/", "use an alternate installation path")
	cmd.Flags().BoolVar(
		&flags.DownloadOnly, "download-only", false,
		"download the package, but do not perform upgrade")
	cmd.Flags().BoolVar(
		&flags.UpgradeOnly, "upgrade-only", false,
		"perform the upgrade using a pre-downloaded file (requires --file)")
	cmd.Flags().StringVar(
		&flags.UpgradeFile, "file", "",
		"the path to the upgrade file (if doing a download or upgrade only)")

	return cmd
}
