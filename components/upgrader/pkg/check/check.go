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
package check

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/paxautoma/operos/components/common/gatekeeper"
)

type UpgradeCheck struct {
	GatekeeperAddress string
	NoGatekeeperTLS   bool
	Version           string
	ClusterID         string
	TargetPath        string

	DialFunc func(string, bool) (io.Closer, gatekeeper.GatekeeperClient, error)
}

func DialGatekeeper(address string, noTLS bool) (io.Closer, gatekeeper.GatekeeperClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTimeout(10 * time.Second),
	}
	if noTLS {
		opts = append(opts, grpc.WithInsecure())
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to dial gatekeeper")
	}

	client := gatekeeper.NewGatekeeperClient(conn)
	return conn, client, nil
}

func (c *UpgradeCheck) DoCheck() error {
	if c.DialFunc == nil {
		c.DialFunc = DialGatekeeper
	}

	log.Infof("current version is v%s", c.Version)

	var res *gatekeeper.UpgradeCheckResp
	var err error
	if res, err = c.checkGatekeeper(); err != nil {
		return errors.Wrap(err, "failed to check gatekeeper")
	}

	newVersion := res.GetAvailableVersion()
	if newVersion == c.Version {
		log.Info("up to date")
		return nil
	}

	log.Infof("update available to version %s", newVersion)

	for _, url := range res.GetDownloadUrls() {
		var fileName string
		if fileName, err = c.downloadFile(c.TargetPath, url); err != nil {
			return errors.Wrapf(err, "failed to download file %s", url)
		}

		if err := c.untarFile(c.TargetPath, fileName); err != nil {
			return errors.Wrapf(err, "failed to untar file %s", fileName)
		}
	}

	return nil
}

func (c *UpgradeCheck) checkGatekeeper() (*gatekeeper.UpgradeCheckResp, error) {
	conn, client, err := c.DialFunc(c.GatekeeperAddress, c.NoGatekeeperTLS)
	defer conn.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial gatekeeper")
	}

	res, err := client.UpgradeCheck(context.Background(), &gatekeeper.UpgradeCheckReq{
		ClusterId:      c.ClusterID,
		CurrentVersion: c.Version,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed execute check")
	}

	return res, nil
}

func (c *UpgradeCheck) downloadFile(targetPath, url string) (string, error) {
	client := grab.NewClient()
	req, _ := grab.NewRequest(targetPath, url)

	log.Infof("downloading file: %s", url)
	resp := client.Do(req)

	var t *time.Ticker
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		t = time.NewTicker(500 * time.Millisecond)
	} else {
		t = time.NewTicker(5 * time.Second)
	}
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			log.Debugf("transferred %v / %v bytes (%.2f%%)",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress())
		case <-resp.Done:
			break Loop
		}
	}

	if err := resp.Err(); err != nil {
		return "", errors.Wrap(err, "download failed")
	}

	log.Infof("download saved to %v", resp.Filename)
	return resp.Filename, nil
}

func (c *UpgradeCheck) untarFile(targetPath, fileName string) error {
	r, err := os.Open(fileName)
	defer r.Close()
	if err != nil {
		return errors.Wrap(err, "could not open file for reading")
	}

	gzr, err := gzip.NewReader(r)
	log.Debugf("hello: %v %v", gzr, err)
	if err != nil {
		return errors.Wrap(err, "could not read from file")
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return errors.Wrap(err, "failed to read valid tar data")
		}

		target := path.Join(targetPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return errors.Wrapf(err, "failed to create directory %s", target)
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return errors.Wrapf(err, "failed to write file %s", target)
			}
			defer f.Close()

			if _, err := io.Copy(f, tr); err != nil {
				return errors.Wrapf(err, "failed to write to file %s", target)
			}
		}
	}
}
