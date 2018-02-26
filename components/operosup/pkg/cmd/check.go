/*
Copyright 2018 Pax Automa Systems, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/paxautoma/operos/components/operosup/pkg/client"
	"github.com/paxautoma/operos/components/operosup/pkg/steps"
)

func NewCheckCommand(dialer *client.Dialer) *cobra.Command {
	var flags steps.UpgradeCheckFlags

	var cmd = &cobra.Command{
		Use:   "check",
		Short: "Check for upgrades",
		Long: `Check Pax Automa servers for the latest Operos upgrade. If an upgrade is available,
					flag this in the controller state.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return steps.DoCheck(dialer, &flags)
		},
	}

	cmd.Flags().StringVarP(
		&flags.Version, "version", "v", os.Getenv("OPEROS_VERSION"),
		"current Operos version")
	cmd.Flags().StringVarP(
		&flags.ClusterID, "cluster", "c", os.Getenv("OPEROS_INSTALL_ID"),
		"cluster install ID")

	return cmd
}
