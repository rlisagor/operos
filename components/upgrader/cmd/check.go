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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/paxautoma/operos/components/common"
	"github.com/paxautoma/operos/components/upgrader/pkg/check"
)

var gatekeeperAddress = flag.String("gatekeeper", "gatekeeper.paxautoma.com:57345", "address of the Gatkeeper server (host:port)")
var noGatekeeperTLS = flag.Bool("no-gatekeeper-tls", false, "do not use TLS with Gatekeeper")
var version = flag.String("version", "", "current version")
var clusterID = flag.String("cluster", "", "cluster ID")
var targetPath = flag.String("dir", ".", "directory where the downloaded files should go")

func main() {
	common.SetupLogging()
	defer common.LogPanic()

	flag.Parse()

	if *version == "" {
		fmt.Fprintln(os.Stderr, "The version flag is mandatory")
		os.Exit(1)
	}

	if *clusterID == "" {
		fmt.Fprintln(os.Stderr, "The cluster flag is mandatory")
		os.Exit(1)
	}

	c := check.UpgradeCheck{
		GatekeeperAddress: *gatekeeperAddress,
		NoGatekeeperTLS:   *noGatekeeperTLS,
		Version:           *version,
		ClusterID:         *clusterID,
		TargetPath:        *targetPath,
	}
	if err := c.DoCheck(); err != nil {
		log.Fatal(err.Error())
	}
}
