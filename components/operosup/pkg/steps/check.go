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
package steps

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/paxautoma/operos/components/common/gatekeeper"
	"github.com/paxautoma/operos/components/operosup/pkg/client"
)

type UpgradeCheckFlags struct {
	Version   string
	ClusterID string
}

func DoCheck(dialer *client.Dialer, flags *UpgradeCheckFlags) error {
	log.Infof("current version is v%s", flags.Version)

	var res *gatekeeper.UpgradeCheckResp
	var err error
	if res, err = checkGatekeeper(dialer, flags.Version, flags.ClusterID); err != nil {
		return errors.Wrap(err, "failed to check gatekeeper")
	}

	newVersion := res.GetAvailableVersion()
	if newVersion == flags.Version {
		log.Info("up to date")
		return nil
	}

	log.Infof("update available to version %s", newVersion)

	if err := notifyTeamster(dialer); err != nil {
		return errors.Wrap(err, "failed to notify teamster")
	}

	return nil
}

func checkGatekeeper(dialer *client.Dialer, version, clusterID string) (*gatekeeper.UpgradeCheckResp, error) {
	conn, client, err := dialer.DialGatekeeper()
	defer conn.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial gatekeeper")
	}

	res, err := client.UpgradeCheck(context.Background(), &gatekeeper.UpgradeCheckReq{
		ClusterId:      clusterID,
		CurrentVersion: version,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed execute check")
	}

	return res, nil
}

func notifyTeamster(dialer *client.Dialer) error {
	return nil
}
