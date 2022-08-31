/*
Copyright 2021 The Dapr Authors
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
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dapr/cli/pkg/kubernetes"
	"github.com/dapr/cli/pkg/print"
	"github.com/dapr/cli/pkg/standalone"
)

var (
	uninstallNamespace        string
	uninstallKubernetes       bool
	uninstallAll              bool
	uninstallContainerRuntime string
	uninstallDir              string
)

// UninstallCmd is a command from removing a Dapr installation.
var UninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Dapr runtime. Supported platforms: Kubernetes and self-hosted",
	Example: `
# Uninstall from self-hosted mode
dapr uninstall

# Uninstall from self-hosted mode and remove .dapr directory, Redis, Placement and Zipkin containers
dapr uninstall --all

# Uninstall from Kubernetes
dapr uninstall -k

# Uninstall Dapr from non-default install directory (default is $HOME/.dapr)
dapr uninstall --install-dir <path-to-install-directory>
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("network", cmd.Flags().Lookup("network"))
		viper.BindPFlag("install-path", cmd.Flags().Lookup("install-path"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if uninstallKubernetes {
			if len(strings.TrimSpace(uninstallDir)) != 0 {
				print.FailureStatusEvent(os.Stderr, "--install-dir is only valid for self-hosted mode")
				os.Exit(1)
			}

			print.InfoStatusEvent(os.Stdout, "Removing Dapr from your cluster...")
			err = kubernetes.Uninstall(uninstallNamespace, uninstallAll, timeout)
		} else {
			print.InfoStatusEvent(os.Stdout, "Removing Dapr from your machine...")
			dockerNetwork := viper.GetString("network")
			err = standalone.Uninstall(uninstallAll, dockerNetwork, uninstallContainerRuntime, uninstallDir)
		}

		if err != nil {
			print.FailureStatusEvent(os.Stderr, fmt.Sprintf("Error removing Dapr: %s", err))
		} else {
			print.SuccessStatusEvent(os.Stdout, "Dapr has been removed successfully")
		}
	},
}

func init() {
	defaultInstallDir := standalone.DefaultDaprDirPath()

	UninstallCmd.Flags().BoolVarP(&uninstallKubernetes, "kubernetes", "k", false, "Uninstall Dapr from a Kubernetes cluster")
	UninstallCmd.Flags().UintVarP(&timeout, "timeout", "", 300, "The timeout for the Kubernetes uninstall")
	UninstallCmd.Flags().BoolVar(&uninstallAll, "all", false, "Remove .dapr directory, Redis, Placement and Zipkin containers on local machine, and CRDs on a Kubernetes cluster")
	UninstallCmd.Flags().String("network", "", "The Docker network from which to remove the Dapr runtime")
	UninstallCmd.Flags().StringVarP(&uninstallNamespace, "namespace", "n", "dapr-system", "The Kubernetes namespace to uninstall Dapr from")
	UninstallCmd.Flags().BoolP("help", "h", false, "Print this help message")
	UninstallCmd.Flags().StringVarP(&uninstallContainerRuntime, "container-runtime", "", "docker", "The container runtime to use (defaults to docker)")
	UninstallCmd.Flags().StringVarP(&uninstallDir, "install-dir", "", defaultInstallDir, "The directory to uninstall dapr from, for example: /usr/local/dapr")
	RootCmd.AddCommand(UninstallCmd)
}
