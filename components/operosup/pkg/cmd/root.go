package cmd

import (
	"github.com/spf13/cobra"

	"github.com/paxautoma/operos/components/operosup/pkg/client"
)

func NewRootCommand() *cobra.Command {
	dialer := client.Dialer{}

	cmd := &cobra.Command{
		Use:          "operosup",
		Short:        "Operos upgrade tool",
		SilenceUsage: true,
	}

	cmd.PersistentFlags().StringVar(
		&dialer.GatekeeperAddress, "gatekeeper", "localhost:57345",
		"address:port of the Gatekeeper instance")
	cmd.PersistentFlags().BoolVar(
		&dialer.NoGatekeeperTLS, "no-gatekeeper-tls", false,
		"do not use TLS when communicating with Gatekeeper")
	cmd.PersistentFlags().StringVar(
		&dialer.TeamsterAddress, "teamster", "localhost:4780",
		"address:port of the Teamster instance")

	cmd.AddCommand(NewCheckCommand(&dialer))
	cmd.AddCommand(NewUpgradeCommand(&dialer))

	return cmd
}
